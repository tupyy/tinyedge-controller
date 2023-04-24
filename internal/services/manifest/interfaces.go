package manifest

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

//go:generate moq -out device_reader_moq.go . DeviceReader
type DeviceReader interface {
	GetDevice(ctx context.Context, id string) (entity.Device, error)
	GetNamespace(ctx context.Context, id string) (entity.Namespace, error)
	GetSet(ctx context.Context, id string) (entity.Set, error)
}

type ManifestReader interface {
	GetManifest(ctx context.Context, id string) (entity.Manifest, error)
	GetManifests(ctx context.Context, repo entity.Repository) ([]entity.Manifest, error)
}

type ManifestWriter interface {
	InsertManifest(ctx context.Context, manifest entity.Manifest) error
	UpdateManifest(ctx context.Context, manifest entity.Manifest) error
	DeleteManifest(ctx context.Context, id string) error

	CreateRelation(ctx context.Context, relation entity.Relation) error
	DeleteRelation(ctx context.Context, relation entity.Relation) error
}

//go:generate moq -out manifest_rw_moq.go . ManifestReaderWriter
type ManifestReaderWriter interface {
	ManifestReader
	ManifestWriter
}

//go:generate moq -out git_reader_moq.go . GitReader
type GitReader interface {
	GetManifests(ctx context.Context, repo entity.Repository) ([]entity.Manifest, error)
}
