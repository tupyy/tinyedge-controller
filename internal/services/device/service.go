package device

import (
	"context"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
)

type Service struct {
	pgDeviceRepo DeviceReaderWriter
}

func New(pgDeviceRepo DeviceReaderWriter) *Service {
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
		return errService.NewResourceAlreadyExistsError("namespace", namespace.Name)
	} else if _, ok := err.(errService.ResourseNotFoundError); !ok {
		return err
	}

	return w.pgDeviceRepo.CreateNamespace(ctx, namespace)
}

func (w *Service) CreateSet(ctx context.Context, set entity.Set) error {
	_, err := w.GetSet(ctx, set.Name)
	if err == nil {
		return errService.NewResourceAlreadyExistsError("set", set.Name)
	} else if _, ok := err.(errService.ResourseNotFoundError); !ok {
		return err
	}

	// check for namespace
	if set.NamespaceID == "" {
		return fmt.Errorf("unable to create set. namespace is missing")
	}

	_, err = w.GetNamespace(ctx, set.NamespaceID)
	if err != nil {
		return err
	}

	return w.pgDeviceRepo.CreateSet(ctx, set)
}
