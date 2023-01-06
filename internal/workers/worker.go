package workers

import (
	"context"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/repo/git"
	"github.com/tupyy/tinyedge-controller/internal/services/configuration"
	"github.com/tupyy/tinyedge-controller/internal/services/manifest"
	"go.uber.org/zap"
)

type GitOpsWorker struct {
	manifestService *manifest.Service
	confService     *configuration.ConfigurationService
	gitRepo         *git.GitRepo
}

func NewGitOpsWorker(w *manifest.Service, c *configuration.ConfigurationService, g *git.GitRepo) *GitOpsWorker {
	return &GitOpsWorker{w, c, g}
}

func (g *GitOpsWorker) Do(ctx context.Context) error {
	repos, err := g.manifestService.GetRepositories(ctx)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		_, err := g.gitRepo.Open(ctx, repo)
		if err != nil {
			zap.S().Errorw("unable to open repo", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
		}

		err = g.gitRepo.Pull(ctx, repo)
		if err != nil {
			zap.S().Errorw("unable to pull from origin", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
			continue
		}

		headSha, err := g.gitRepo.GetHeadSha(ctx, repo)
		if err != nil {
			zap.S().Errorw("unable to get head from repo", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
			continue
		}

		if headSha == repo.CurrentHeadSha {
			zap.S().Debugw("repo is up to date. skipping...", "repo.url", repo.Url, "head_sha", repo.TargetHeadSha)
			continue
		}

		zap.S().Infow("repo has been updated", "repo_url", repo.Url, "head sha", headSha, "repo_current_sha", repo.CurrentHeadSha)

		repo.TargetHeadSha = headSha
		if err := g.manifestService.UpdateRepository(ctx, repo); err != nil {
			zap.S().Errorw("unable to update target sha of the repository", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
			continue
		}

		// get all the manifest from the git repository
		gitManifests, err := g.gitRepo.GetManifests(ctx, repo)
		if err != nil {
			zap.S().Errorw("unable to fetch manifest from git repository", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
			continue
		}

		created, updated, _, err := g.manifestService.UpdateManifests(ctx, repo, gitManifests)
		if err != nil {
			zap.S().Errorw("unable to update repository's manifests", "error", err, "repo_id", repo.Id)
			continue
		}

		// create relations between namespaces, sets and devices for the new manifests
		for _, m := range created {
			zap.S().Infof("create relations from manifest %q", m.Name)
			if err := g.manifestService.CreateRelations(ctx, m); err != nil {
				return fmt.Errorf("unable to create relations for manifest %q: %w", repo.Id, err)
			}
		}

		// update the relations of the existing manifests
		for _, m := range updated {
			zap.S().Infof("update relations for manifest %+v", m)
			if err := g.manifestService.UpdateRelations(ctx, m); err != nil {
				return fmt.Errorf("unable to update relations for manifest %q: %w", repo.Id, err)
			}
		}

		// all done. set current sha to target sha
		repo.CurrentHeadSha = headSha
		if err := g.manifestService.UpdateRepository(ctx, repo); err != nil {
			zap.S().Errorw("unable to update current sha of the repository", "error", err, "repo_id", repo.Id)
		}

		zap.S().Infow("repository and manifests updated", "repo_id", repo.Id, "repo_url", repo.Url, "repo_current_sha", repo.CurrentHeadSha)
	}
	return nil
}

func (g *GitOpsWorker) Name() string {
	return "gitOpsWorker"
}
