package manifest

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

//go:generate moq -out reference_reader_moq.go . ReferenceReader
type ReferenceReader interface {
	GetReferences(ctx context.Context, repo entity.Repository) ([]entity.Reference, error)
}

//go:generate moq -out git_reader_moq.go . GitReader
type GitReader interface {
	GetWorkload(ctx context.Context, ref entity.Reference) (entity.Workload, error)
}

//go:generate moq -out secret_reader_moq.go . SecretReader
type SecretReader interface {
	GetSecret(ctx context.Context, path, key string) (entity.Secret, error)
}
