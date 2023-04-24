package workers

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/services"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
	"go.uber.org/zap"
)

type GitOpsWorker struct {
	manifestService   *services.Manifest
	repositoryService *services.Repository
	confService       *services.Configuration
}

func NewGitOpsWorker(r *services.Repository, m *services.Manifest, c *services.Configuration) *GitOpsWorker {
	return &GitOpsWorker{
		manifestService:   m,
		repositoryService: r,
		confService:       c,
	}
}

func (g *GitOpsWorker) Do(ctx context.Context) error {
	repos, err := g.repositoryService.GetRepositories(ctx)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		err := g.repositoryService.Open(ctx, repo)
		if err != nil {
			if errService.IsResourceNotFound(err) {
				// clone it
				clone, err := g.repositoryService.Clone(ctx, repo)
				if err != nil {
					zap.S().Errorw("unable to clone repository", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
					continue
				}
				// save the clone and exit the loop
				if err := g.repositoryService.Update(ctx, clone); err != nil {
					zap.S().Errorw("unable to update repository", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
				}
				continue
			}
		}

		r, err := g.repositoryService.PullRepository(ctx, repo)
		if err != nil {
			zap.S().Errorw("unable to pull repository", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
			continue
		}

		if r.TargetHeadSha == r.CurrentHeadSha {
			zap.S().Debugw("repo is up to date. skipping...", "repo.url", repo.Url, "head_sha", repo.TargetHeadSha)
			continue
		}

		zap.S().Infow("changes detected in repo", "repo_url", repo.Url, "head sha", r.TargetHeadSha, "repo_current_sha", r.CurrentHeadSha)

		if err := g.manifestService.UpdateManifests(ctx, r); err != nil {
			zap.S().Errorw("unable to update repository's manifests", "error", err, "repo_id", r.Id, "repo_url", r.Url)
			continue
		}

		// all done. set current sha to target sha
		r.CurrentHeadSha = r.TargetHeadSha
		if err := g.repositoryService.Update(ctx, r); err != nil {
			zap.S().Errorw("unable to update current sha of the repository", "error", err, "repo_id", r.Id)
		}

		zap.S().Infow("repository and references updated", "repo_id", r.Id, "repo_url", r.Url, "repo_current_sha", r.CurrentHeadSha)
	}
	return nil
}

func (g *GitOpsWorker) Name() string {
	return "gitOpsWorker"
}
