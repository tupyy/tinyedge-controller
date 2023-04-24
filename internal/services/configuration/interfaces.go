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
