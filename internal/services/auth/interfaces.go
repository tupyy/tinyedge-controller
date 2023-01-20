package auth

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type DeviceReader interface {
	GetDevice(ctx context.Context, id string) (entity.Device, error)
}
