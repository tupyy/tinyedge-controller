package device

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type SecretReaderWriter interface {
	Get(ctx context.Context, id string) (string, error)
	Put(ctx context.Context, secret entity.Secret) error
}

type WorkloadReader interface {
	Get(ctx context.Context, id string) (entity.Workload, error)
}

type DeviceReader interface {
	Get(ctx context.Context, id string) (entity.Device, error)
	GetBySet(ctx context.Context, setID string) ([]entity.Device, error)
}

type DeviceWriter interface {
	Create(ctx context.Context, device entity.Device) error
	Update(ctx context.Context, device entity.Device) error
	Delete(ctx context.Context, device entity.Device) error
}
