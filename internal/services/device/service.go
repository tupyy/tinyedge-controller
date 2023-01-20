package device

import (
	"context"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
)

type Service struct {
	pgDeviceRepo common.DeviceReaderWriter
}

func New(pgDeviceRepo common.DeviceReaderWriter) *Service {
	return &Service{pgDeviceRepo: pgDeviceRepo}
}

func (w *Service) GetNamespaces(ctx context.Context) ([]entity.Namespace, error) {
	return w.pgDeviceRepo.GetNamespaces(ctx)
}

func (w *Service) GetNamespace(ctx context.Context, id string) (entity.Namespace, error) {
	return w.pgDeviceRepo.GetNamespace(ctx, id)
}

func (w *Service) GetSets(ctx context.Context) ([]entity.Set, error) {
	return w.pgDeviceRepo.GetSets(ctx)
}

func (w *Service) GetSet(ctx context.Context, id string) (entity.Set, error) {
	return w.pgDeviceRepo.GetSet(ctx, id)
}

func (w *Service) GetDevice(ctx context.Context, id string) (entity.Device, error) {
	return w.pgDeviceRepo.GetDevice(ctx, id)
}

func (w *Service) GetDevices(ctx context.Context) ([]entity.Device, error) {
	return w.pgDeviceRepo.GetDevices(ctx)
}

func (w *Service) UpdateDevice(ctx context.Context, device entity.Device) error {
	return w.pgDeviceRepo.Update(ctx, device)
}

func (w *Service) CreateNamespace(ctx context.Context, namespace entity.Namespace) error {
	_, err := w.GetNamespace(ctx, namespace.Name)
	if err == nil {
		return fmt.Errorf("namespace %q already exists: %w", namespace.Name, common.ErrResourceAlreadyExists)
	} else if !common.IsResourceNotFound(err) {
		return err
	}

	return w.pgDeviceRepo.CreateNamespace(ctx, namespace)
}

func (w *Service) CreateSet(ctx context.Context, set entity.Set) error {
	_, err := w.GetSet(ctx, set.Name)
	if err == nil {
		return fmt.Errorf("set %q already exists: %w", set.Name, common.ErrResourceAlreadyExists)
	} else if !common.IsResourceNotFound(err) {
		return err
	}

	// check for namespace
	if set.NamespaceID == "" {
		return fmt.Errorf("unable to create set. namespace is missing")
	}

	_, err = w.GetNamespace(ctx, set.NamespaceID)
	if common.IsResourceNotFound(err) {
		return fmt.Errorf("%w: namespace %q", common.ErrResourceNotFound, set.NamespaceID)
	}

	return w.pgDeviceRepo.CreateSet(ctx, set)
}
