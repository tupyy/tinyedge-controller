package common

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type ReferenceReader interface {
	GetReferences(ctx context.Context) ([]entity.ManifestReference, error)
	GetReference(ctx context.Context, id string) (entity.ManifestReference, error)
	GetRepositoryReferences(ctx context.Context, repo entity.Repository) ([]entity.ManifestReference, error)
	GetDeviceReferences(ctx context.Context, deviceID string) ([]entity.ManifestReference, error)
	GetSetReferences(ctx context.Context, setID string) ([]entity.ManifestReference, error)
	GetNamespaceReferences(ctx context.Context, setID string) ([]entity.ManifestReference, error)
}

type ReferenceWriter interface {
	InsertReference(ctx context.Context, ref entity.ManifestReference) error
	UpdateReference(ctx context.Context, ref entity.ManifestReference) error
	DeleteReference(ctx context.Context, ref entity.ManifestReference) error

	CreateRelation(ctx context.Context, relation entity.ReferenceRelation) error
	DeleteRelation(ctx context.Context, relation entity.ReferenceRelation) error
}

type ReferenceReaderWriter interface {
	ReferenceReader
	ReferenceWriter
}
