package workers

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/services/configuration"
	"github.com/tupyy/tinyedge-controller/internal/services/manifest"
	"go.uber.org/zap"
)

type GitOpsWorker struct {
	manifestService *manifest.Service
	confService     *configuration.ConfigurationService
}

func NewGitOpsWorker(w *manifest.Service, c *configuration.ConfigurationService) *GitOpsWorker {
	return &GitOpsWorker{w, c}
}

func (g *GitOpsWorker) Do(ctx context.Context) error {
	repos, err := g.manifestService.GetRepositories(ctx)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		r, err := g.manifestService.PullRepository(ctx, repo)
		if err != nil {
			zap.S().Errorf("unable to pull repository", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
			continue
		}

		if r.TargetHeadSha == repo.CurrentHeadSha {
			zap.S().Debugw("repo is up to date. skipping...", "repo.url", repo.Url, "head_sha", repo.TargetHeadSha)
			continue
		}

		zap.S().Infow("repo has been updated", "repo_url", repo.Url, "head sha", repo.TargetHeadSha, "repo_current_sha", repo.CurrentHeadSha)

		if err := g.manifestService.UpdateRepository(ctx, repo); err != nil {
			zap.S().Errorw("unable to update target sha of the repository", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
			continue
		}

		if err := g.manifestService.UpdateManifests(ctx, repo); err != nil {
			zap.S().Errorw("unable to update repository's manifests", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
			continue
		}

		// all done. set current sha to target sha
		r.CurrentHeadSha = r.TargetHeadSha
		if err := g.manifestService.UpdateRepository(ctx, r); err != nil {
			zap.S().Errorw("unable to update current sha of the repository", "error", err, "repo_id", repo.Id)
		}

		zap.S().Infow("repository and manifests updated", "repo_id", repo.Id, "repo_url", repo.Url, "repo_current_sha", repo.CurrentHeadSha)
	}
	return nil
}

func (g *GitOpsWorker) Name() string {
	return "gitOpsWorker"
}
