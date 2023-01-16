package common

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
