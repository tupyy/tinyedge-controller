package workers

import (
	"context"
	"errors"

	"github.com/tupyy/tinyedge-controller/internal/services/common"
	"github.com/tupyy/tinyedge-controller/internal/services/configuration"
	reference "github.com/tupyy/tinyedge-controller/internal/services/manifest"
	"github.com/tupyy/tinyedge-controller/internal/services/repository"
	"go.uber.org/zap"
)

type GitOpsWorker struct {
	referenceService  *reference.Service
	repositoryService *repository.Service
	confService       *configuration.Service
}

func NewGitOpsWorker(ref *reference.Service, r *repository.Service, c *configuration.Service) *GitOpsWorker {
	return &GitOpsWorker{
		referenceService:  ref,
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
			if errors.Is(err, common.ErrResourceNotFound) {
				// clone it
				clone, err := g.repositoryService.Clone(ctx, repo.Url, repo.Id)
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

		if r.TargetHeadSha == repo.CurrentHeadSha {
			zap.S().Debugw("repo is up to date. skipping...", "repo.url", repo.Url, "head_sha", repo.TargetHeadSha)
			continue
		}

		zap.S().Infow("changes detected in repo", "repo_url", repo.Url, "head sha", repo.TargetHeadSha, "repo_current_sha", repo.CurrentHeadSha)

		if err := g.repositoryService.Update(ctx, repo); err != nil {
			zap.S().Errorw("unable to update target sha of the repository", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
			continue
		}

		if err := g.referenceService.UpdateReferences(ctx, repo); err != nil {
			zap.S().Errorw("unable to update repository's manifests", "error", err, "repo_id", repo.Id, "repo_url", repo.Url)
			continue
		}

		// all done. set current sha to target sha
		r.CurrentHeadSha = r.TargetHeadSha
		if err := g.repositoryService.Update(ctx, r); err != nil {
			zap.S().Errorw("unable to update current sha of the repository", "error", err, "repo_id", repo.Id)
		}

		zap.S().Infow("repository and references updated", "repo_id", repo.Id, "repo_url", repo.Url, "repo_current_sha", repo.CurrentHeadSha)
	}
	return nil
}

func (g *GitOpsWorker) Name() string {
	return "gitOpsWorker"
}
