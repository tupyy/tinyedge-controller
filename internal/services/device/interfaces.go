package device

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

// DeviceReader is an interface that groups all the methods allowing to query/get devices.
type DeviceReader interface {
	GetDevice(ctx context.Context, id string) (entity.Device, error)
	GetDevices(ctx context.Context) ([]entity.Device, error)
	GetNamespace(ctx context.Context, id string) (entity.Namespace, error)
	GetNamespaces(ctx context.Context) ([]entity.Namespace, error)
	GetSet(ctx context.Context, id string) (entity.Set, error)
	GetSets(ctx context.Context) ([]entity.Set, error)
	GetConfiguration(ctx context.Context, id string) (entity.Configuration, error)
}

// DeviceWriter allows creating a device.
type DeviceWriter interface {
	Create(ctx context.Context, device entity.Device) error
	Update(ctx context.Context, device entity.Device) error
	CreateSet(ctx context.Context, set entity.Set) error
	CreateNamespace(ctx context.Context, namespace entity.Namespace) error
}

type DeviceReaderWriter interface {
	DeviceReader
	DeviceWriter
}
