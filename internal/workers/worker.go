package workers

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/services/workload"
	"go.uber.org/zap"
)

type GitOpsWorker struct {
	service *workload.WorkloadService
}

func NewGitOpsWorker(s *workload.WorkloadService) *GitOpsWorker {
	return &GitOpsWorker{s}
}

func (g *GitOpsWorker) Do(ctx context.Context) error {
	repos, err := g.service.GetRepositories(ctx)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		if err := g.service.UpdateRepository(ctx, repo); err != nil {
			zap.S().Errorw("unable to update repository", "error", err, "repo_id", repo.Id)
			continue
		}
		if err := g.service.UpdateManifestsFromGit(ctx, repo); err != nil {
			zap.S().Errorw("unable to update repository's manifests", "error", err, "repo_id", repo.Id)
			continue
		}

		// get manifest and create relations
		manifests, err := g.service.GetManifests(ctx, repo)
		if err != nil {
			zap.S().Errorw("unable to get repository's manifests", "error", err, "repo_id", repo.Id)
			continue
		}

		for _, m := range manifests {
			zap.S().Infof("create relations from manifest %q", m.Name)
			if err := g.service.CreateRelations(ctx, m); err != nil {
				zap.S().Errorw("unable to create relations", "error", err, "repo_id", repo.Id)
			}
		}

		// zap.S().Infow("repository and manifests updated", "repo_id", repo.Id)
	}
	return nil
}

func (g *GitOpsWorker) Name() string {
	return "gitOpsWorker"
}
