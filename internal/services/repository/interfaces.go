package repository

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type RepositoryReader interface {
	GetRepositories(ctx context.Context) ([]entity.Repository, error)
}

type RepositoryWriter interface {
	InsertRepository(ctx context.Context, r entity.Repository) error
	UpdateRepository(ctx context.Context, r entity.Repository) error
}

type RepositoryReaderWriter interface {
	RepositoryReader
	RepositoryWriter
}

type GitReader interface {
	Open(ctx context.Context, r entity.Repository) (entity.Repository, error)
	Pull(ctx context.Context, r entity.Repository) error
	GetHeadSha(ctx context.Context, r entity.Repository) (string, error)
	GetManifests(ctx context.Context, repo entity.Repository) ([]entity.Manifest, error)
	GetManifest(ctx context.Context, repo entity.Repository, file string) (entity.Manifest, error)
}

type GitWriter interface {
	Clone(ctx context.Context, remoteRepo entity.Repository) (entity.Repository, error)
}

type GitReaderWriter interface {
	GitReader
	GitWriter
}

type SecretReader interface {
	GetCredentialsFunc(ctx context.Context, authType entity.RepositoryAuthType, secretPath string) entity.CredentialsFunc
}
