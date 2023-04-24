package device

import (
	"context"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
	"go.uber.org/zap"
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

func (w *Service) GetDefaultNamespace(ctx context.Context) (entity.Namespace, error) {
	return w.pgDeviceRepo.GetDefaultNamespace(ctx)
}

// DeleteNamespace removes a namespace. The last namespace cannot be removed
func (w *Service) DeleteNamespace(ctx context.Context, id string) (entity.Namespace, error) {
	namespace, err := w.GetNamespace(ctx, id)
	if err != nil {
		return entity.Namespace{}, err
	}

	// count the namespaces
	namespaces, err := w.GetNamespaces(ctx)
	if err != nil {
		return entity.Namespace{}, err
	}

	if len(namespaces) == 1 {
		return entity.Namespace{}, errService.NewDeleteResourceError("namespace", id, "cannot delete the last namespace")
	}

	// if namespace was the default one, set the next one as default
	if namespace.IsDefault {
		for _, n := range namespaces {
			if n.Name != namespace.Name {
				n.IsDefault = true
				if err := w.pgDeviceRepo.UpdateNamespace(ctx, n); err != nil {
					return entity.Namespace{}, err
				}
				break
			}
		}
	}

	defaultNamespace, err := w.GetDefaultNamespace(ctx)
	if err != nil {
		return entity.Namespace{}, nil
	}

	// set the namespace for all devices in the deleted namespace
	for _, id := range namespace.Devices {
		device, err := w.GetDevice(ctx, id)
		if err != nil {
			return entity.Namespace{}, err
		}
		device.NamespaceID = defaultNamespace.Name
		if err := w.UpdateDevice(ctx, device); err != nil {
			return entity.Namespace{}, err
		}
		zap.S().Infof("device %q was moved into the default namespace %q", id, device.NamespaceID)
	}

	if err := w.pgDeviceRepo.DeleteNamespace(ctx, id); err != nil {
		return entity.Namespace{}, err
	}

	zap.S().Infof("Namespace %q was deleted", id)
	return namespace, nil
}

func (w *Service) UpdateNamespace(ctx context.Context, namespace entity.Namespace) (entity.Namespace, error) {
	_, err := w.GetNamespace(ctx, namespace.Name)
	if err != nil {
		return entity.Namespace{}, err
	}

	if namespace.IsDefault {
		// get the default namespace and if different update it
		defaultNamespace, err := w.GetDefaultNamespace(ctx)
		if err != nil {
			return entity.Namespace{}, err
		}
		if defaultNamespace.Name != namespace.Name {
			defaultNamespace.IsDefault = false
			if err := w.pgDeviceRepo.UpdateNamespace(ctx, defaultNamespace); err != nil {
				return entity.Namespace{}, err
			}
		}
	}

	if err := w.pgDeviceRepo.UpdateNamespace(ctx, namespace); err != nil {
		return entity.Namespace{}, err
	}

	return namespace, nil
}

func (w *Service) GetSets(ctx context.Context) ([]entity.Set, error) {
	return w.pgDeviceRepo.GetSets(ctx)
}

func (w *Service) GetSet(ctx context.Context, id string) (entity.Set, error) {
	return w.pgDeviceRepo.GetSet(ctx, id)
}

func (w *Service) DeleteSet(ctx context.Context, id string) (entity.Set, error) {
	set, err := w.GetSet(ctx, id)
	if err != nil {
		return entity.Set{}, err
	}

	if err := w.pgDeviceRepo.DeleteSet(ctx, set.Name); err != nil {
		return set, err
	}

	return set, nil
}

func (w *Service) GetDevice(ctx context.Context, id string) (entity.Device, error) {
	return w.pgDeviceRepo.GetDevice(ctx, id)
}

func (w *Service) GetDevices(ctx context.Context) ([]entity.Device, error) {
	return w.pgDeviceRepo.GetDevices(ctx)
}

func (w *Service) UpdateDevice(ctx context.Context, device entity.Device) error {
	err := w.pgDeviceRepo.Update(ctx, device)
	if err != nil {
		return err
	}
	zap.S().Infof("Device %q updated.", device.ID)
	return nil
}

func (w *Service) CreateNamespace(ctx context.Context, namespace entity.Namespace) error {
	_, err := w.GetNamespace(ctx, namespace.Name)
	if err == nil {
		return errService.NewResourceAlreadyExistsError("namespace", namespace.Name)
	} else if _, ok := err.(errService.ResourseNotFoundError); !ok {
		return err
	}

	zap.S().Infof("Namespace %q was created", namespace.Name)

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

	zap.S().Infof("Set %q was created", set.Name)

	return w.pgDeviceRepo.CreateSet(ctx, set)
}
