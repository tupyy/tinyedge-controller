package repository

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
)

type RepositoryService struct {
	gitRepo common.GitReaderWriter
	pgRepo  common.RepositoryReaderWriter
}

func NewRepositoryService(p common.RepositoryReaderWriter, g common.GitReaderWriter) *RepositoryService {
	return &RepositoryService{gitRepo: g, pgRepo: p}
}

func (r *RepositoryService) GetRepositories(ctx context.Context) ([]entity.Repository, error) {
	repos, err := r.pgRepo.GetRepositories(ctx)
	if err != nil {
		return []entity.Repository{}, err
	}
	return repos, nil

}

func (r *RepositoryService) Open(ctx context.Context, repo entity.Repository) error {
	if repo.LocalPath == "" {
		return common.ErrResourceNotFound
	}
	_, err := r.gitRepo.Open(ctx, repo)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryService) Clone(ctx context.Context, url, name string) (entity.Repository, error) {
	return r.gitRepo.Clone(ctx, url, name)
}

func (w *RepositoryService) PullRepository(ctx context.Context, repo entity.Repository) (entity.Repository, error) {
	err := w.gitRepo.Pull(ctx, repo)
	if err != nil {
		return entity.Repository{}, err
	}

	headSha, err := w.gitRepo.GetHeadSha(ctx, repo)
	if err != nil {
		return entity.Repository{}, err
	}

	repo.TargetHeadSha = headSha

	return repo, nil
}

func (w *RepositoryService) Update(ctx context.Context, r entity.Repository) error {
	if err := w.pgRepo.UpdateRepository(ctx, r); err != nil {
		return err
	}
	return nil
}
