package manifest

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type ReferenceReader interface {
	GetRepositoryReferences(ctx context.Context, repo entity.Repository) ([]entity.ManifestReference, error)
	GetReference(ctx context.Context, id string) (entity.ManifestReference, error)
}

type GitReader interface {
	GetManifest(ctx context.Context, ref entity.ManifestReference) (entity.ManifestWork, error)
}

type SecretReader interface {
	GetSecret(ctx context.Context, path, key string) (entity.Secret, error)
}
