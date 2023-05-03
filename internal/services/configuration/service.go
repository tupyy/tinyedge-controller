package configuration

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"go.uber.org/zap"
)

type Service struct {
	deviceReader DeviceReader
}

func New(deviceReader DeviceReader) *Service {
	return &Service{
		deviceReader: deviceReader,
	}
}

func (c *Service) GetDeviceConfiguration(ctx context.Context, deviceID string) (entity.DeviceConfiguration, error) {
	// conf, err := c.cacheReadWriter.Get(ctx, deviceID)
	// if err != nil {
	// 	if !errors.Is(err, common.ErrResourceNotFound) {
	// 		return entity.ConfigurationResponse{}, err
	// 	}

	// create configuration from pg and save it to cache
	device, err := c.deviceReader.GetDevice(ctx, deviceID)
	if err != nil {
		return entity.DeviceConfiguration{}, err
	}
	configuration, err := c.getConfiguration(ctx, device)
	if err != nil {
		return entity.DeviceConfiguration{}, err
	}
	manifests, err := c.getWorkloads(ctx, device)
	if err != nil {
		return entity.DeviceConfiguration{}, err
	}

	var confResponse entity.DeviceConfiguration
	if configuration != nil {
		confResponse = createConfigurationResponse(*configuration, manifests)
	}

	// err = c.cacheReadWriter.Put(ctx, device.ID, confResponse)
	// if err != nil {
	// 	zap.S().Errorw("unable to save configuration to cache", "error", err)
	// }

	zap.S().Debugw("configuration", "configuration", confResponse)
	return confResponse, nil

	// }
	// return conf, nil
}

func (c *Service) getConfiguration(ctx context.Context, device entity.Device) (*entity.Configuration, error) {
	if device.Configuration != nil {
		return device.Configuration, nil
	}

	// if device has no configuration look at the set of the device
	if device.SetID != nil {
		set, err := c.deviceReader.GetSet(ctx, *device.SetID)
		if err != nil {
			return nil, err
		}
		if set.Configuration != nil {
			return set.Configuration, nil
		}
	}

	// if the device has no set or the set has no configuration just grab the configuration from namespace
	namespace, err := c.deviceReader.GetNamespace(ctx, device.NamespaceID)
	if err != nil {
		return nil, err
	}

	// namespace always has a configuration
	return namespace.Configuration, nil
}

func (c *Service) getWorkloads(ctx context.Context, device entity.Device) ([]entity.Workload, error) {
	if len(device.Workloads) > 0 {
		return device.Workloads, nil
	}

	if device.SetID != nil {
		sets, err := c.deviceReader.GetSet(ctx, *device.SetID)
		if err != nil {
			return []entity.Workload{}, err
		}
		if len(sets.Workloads) > 0 {
			return sets.Workloads, nil
		}
	}

	namespace, err := c.deviceReader.GetNamespace(ctx, device.NamespaceID)
	if err != nil {
		return []entity.Workload{}, err
	}
	return namespace.Workloads, nil
}
