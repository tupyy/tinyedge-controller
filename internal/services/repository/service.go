package repository

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
)

type Service struct {
	gitRepo common.GitReaderWriter
	pgRepo  common.RepositoryReaderWriter
}

func NewRepositoryService(p common.RepositoryReaderWriter, g common.GitReaderWriter) *Service {
	return &Service{gitRepo: g, pgRepo: p}
}

func (r *Service) GetRepositories(ctx context.Context) ([]entity.Repository, error) {
	repos, err := r.pgRepo.GetRepositories(ctx)
	if err != nil {
		return []entity.Repository{}, err
	}
	return repos, nil

}

func (r *Service) Open(ctx context.Context, repo entity.Repository) error {
	if repo.LocalPath == "" {
		return common.ErrResourceNotFound
	}
	_, err := r.gitRepo.Open(ctx, repo)
	if err != nil {
		return err
	}
	return nil
}

func (r *Service) Clone(ctx context.Context, url, name string) (entity.Repository, error) {
	return r.gitRepo.Clone(ctx, url, name)
}

func (w *Service) PullRepository(ctx context.Context, repo entity.Repository) (entity.Repository, error) {
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

func (w *Service) Update(ctx context.Context, r entity.Repository) error {
	if err := w.pgRepo.UpdateRepository(ctx, r); err != nil {
		return err
	}
	return nil
}

func (w *Service) Add(ctx context.Context, r entity.Repository) error {
	return w.pgRepo.InsertRepository(ctx, r)
}
