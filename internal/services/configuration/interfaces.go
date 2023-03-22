package configuration

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

//go:generate moq -out device_reader_moq.go . DeviceReader
type DeviceReader interface {
	GetDevice(ctx context.Context, id string) (entity.Device, error)
	GetSet(ctx context.Context, id string) (entity.Set, error)
	GetNamespace(ctx context.Context, id string) (entity.Namespace, error)
}

//go:generate moq -out manifest_reader_moq.go . ManifestReader
type ManifestReader interface {
	GetManifest(ctx context.Context, ref entity.Reference) (entity.WorkloadManifest, error)
}

//go:generate moq -out ref_reader_moq.go . ReferenceReader
type ReferenceReader interface {
	GetReference(ctx context.Context, id string) (entity.Reference, error)
}

//go:generate moq -out configuration_reader_moq.go . ConfigurationReader
type ConfigurationReader interface {
	// GetConfiguration returns the configuration for a device.
	GetConfiguration(ctx context.Context, id string) (entity.Configuration, error)
}
