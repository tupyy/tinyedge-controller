package edge

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

// DeviceReader is an interface that groups all the methods allowing to query/get devices.
type DeviceReader interface {
	GetDevice(ctx context.Context, id string) (entity.Device, error)
}

// DeviceWriter allows creating a device.
type DeviceWriter interface {
	Create(ctx context.Context, device entity.Device) error
	Update(ctx context.Context, device entity.Device) error
}

type DeviceReaderWriter interface {
	DeviceReader
	DeviceWriter
}

type ConfigurationReader interface {
	GetDeviceConfiguration(ctx context.Context, id string) (entity.ConfigurationResponse, error)
}
