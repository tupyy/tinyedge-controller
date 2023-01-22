package configuration

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type DeviceReader interface {
	GetDevice(ctx context.Context, id string) (entity.Device, error)
	GetSet(ctx context.Context, id string) (entity.Set, error)
	GetNamespace(ctx context.Context, id string) (entity.Namespace, error)
}

type ManifestReader interface {
	GetManifest(ctx context.Context, id string) (entity.ManifestWork, error)
}
type ConfigurationReader interface {
	// GetConfiguration returns the configuration for a device.
	GetConfiguration(ctx context.Context, id string) (entity.Configuration, error)
}
