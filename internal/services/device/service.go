package device

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
)

type Service struct {
	pgDeviceRepo common.DeviceReader
}

func New(pgDeviceRepo common.DeviceReader) *Service {
	return &Service{pgDeviceRepo: pgDeviceRepo}
}

func (w *Service) GetNamespaces(ctx context.Context) ([]entity.Namespace, error) {
	return w.pgDeviceRepo.GetNamespaces(ctx)
}

func (w *Service) GetSets(ctx context.Context) ([]entity.Set, error) {
	return w.pgDeviceRepo.GetSets(ctx)
}
