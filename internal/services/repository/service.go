package repository

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
)

type Service struct {
	gitReaderWriter  GitReaderWriter
	repoReaderWriter RepositoryReaderWriter
}

func NewRepositoryService(repoReaderWriter RepositoryReaderWriter, gitReaderWriter GitReaderWriter) *Service {
	return &Service{gitReaderWriter: gitReaderWriter, repoReaderWriter: repoReaderWriter}
}

func (r *Service) GetRepositories(ctx context.Context) ([]entity.Repository, error) {
	repos, err := r.repoReaderWriter.GetRepositories(ctx)
	if err != nil {
		return []entity.Repository{}, err
	}
	return repos, nil

}

func (r *Service) Open(ctx context.Context, repo entity.Repository) error {
	if repo.LocalPath == "" {
		return errService.NewResourceNotFoundError("git repository", repo.Id)
	}
	_, err := r.gitReaderWriter.Open(ctx, repo)
	if err != nil {
		return err
	}
	return nil
}

func (r *Service) Clone(ctx context.Context, url, name string) (entity.Repository, error) {
	return r.gitReaderWriter.Clone(ctx, url, name)
}

func (w *Service) PullRepository(ctx context.Context, repo entity.Repository) (entity.Repository, error) {
	err := w.gitReaderWriter.Pull(ctx, repo)
	if err != nil {
		return entity.Repository{}, err
	}

	headSha, err := w.gitReaderWriter.GetHeadSha(ctx, repo)
	if err != nil {
		return entity.Repository{}, err
	}

	repo.TargetHeadSha = headSha

	return repo, nil
}

func (w *Service) Update(ctx context.Context, r entity.Repository) error {
	if err := w.repoReaderWriter.UpdateRepository(ctx, r); err != nil {
		return err
	}
	return nil
}

func (w *Service) Add(ctx context.Context, r entity.Repository) error {
	return w.repoReaderWriter.InsertRepository(ctx, r)
}
