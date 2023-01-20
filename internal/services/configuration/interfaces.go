package configuration

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type DeviceReader interface {
	GetDevice(ctx context.Context)
}
type ConfigurationReader interface {
	// GetConfiguration returns the configuration for a device.
	GetConfiguration(ctx context.Context, id string) (entity.Configuration, error)
}
