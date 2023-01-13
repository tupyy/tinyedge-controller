package common

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type GitReader interface {
	Open(ctx context.Context, r entity.Repository) (entity.Repository, error)
	Pull(ctx context.Context, r entity.Repository) error
	GetHeadSha(ctx context.Context, r entity.Repository) (string, error)
	GetManifests(ctx context.Context, repo entity.Repository) ([]entity.ManifestWork, error)
	GetManifest(ctx context.Context, ref entity.ManifestReference) (entity.ManifestWork, error)
}

type GitWriter interface {
	Clone(ctx context.Context, url, name string) (entity.Repository, error)
}

type GitReaderWriter interface {
	GitReader
	GitWriter
}

type ConfigurationReader interface {
	// GetConfiguration returns the configuration for a device.
	GetConfiguration(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error)
}
