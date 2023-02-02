package repository

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
)

type Service struct {
	gitReaderWriter  GitReaderWriter
	repoReaderWriter RepositoryReaderWriter
	secretReader     SecretReader
}

func NewRepositoryService(repoReaderWriter RepositoryReaderWriter, gitReaderWriter GitReaderWriter, secretReader SecretReader) *Service {
	return &Service{gitReaderWriter: gitReaderWriter, repoReaderWriter: repoReaderWriter, secretReader: secretReader}
}

func (r *Service) GetRepositories(ctx context.Context) ([]entity.Repository, error) {
	repos, err := r.repoReaderWriter.GetRepositories(ctx)
	if err != nil {
		return []entity.Repository{}, err
	}

	// get credential func
	for i := 0; i < len(repos); i++ {
		repo := &repos[i]
		if repo.AuthType != entity.NoRepositoryAuthType && repo.CredentialsSecretPath != "" {
			repo.Credentials = r.secretReader.GetCredentialsFunc(ctx, repo.AuthType, repo.CredentialsSecretPath)
		}
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

func (r *Service) Clone(ctx context.Context, remoteRepository entity.Repository) (entity.Repository, error) {
	return r.gitReaderWriter.Clone(ctx, remoteRepository)
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
