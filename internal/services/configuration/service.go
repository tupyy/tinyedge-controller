package configuration

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
	"github.com/tupyy/tinyedge-controller/internal/services/manifest"
	"go.uber.org/zap"
)

type ConfigurationService struct {
	manifestReader  *manifest.Service
	deviceReader    common.DeviceReader
	cacheReadWriter common.ConfigurationCacheReaderWriter
}

func New(deviceReader common.DeviceReader, manifestReader *manifest.Service, cacheReaderWriter common.ConfigurationCacheReaderWriter) *ConfigurationService {
	return &ConfigurationService{
		manifestReader:  manifestReader,
		deviceReader:    deviceReader,
		cacheReadWriter: cacheReaderWriter,
	}
}

func (c *ConfigurationService) GetConfiguration(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error) {
	// conf, err := c.cacheReadWriter.Get(ctx, deviceID)
	// if err != nil {
	// 	if !errors.Is(err, common.ErrResourceNotFound) {
	// 		return entity.ConfigurationResponse{}, err
	// 	}

	// create configuration from pg and save it to cache
	device, err := c.deviceReader.GetDevice(ctx, deviceID)
	if err != nil {
		return entity.ConfigurationResponse{}, err
	}
	configuration, err := c.getConfiguration(ctx, device)
	if err != nil {
		return entity.ConfigurationResponse{}, err
	}
	manifests, err := c.getManifests(ctx, device)
	if err != nil {
		return entity.ConfigurationResponse{}, err
	}
	confResponse := createConfigurationResponse(configuration, manifests)

	err = c.cacheReadWriter.Put(ctx, device.ID, confResponse)
	if err != nil {
		zap.S().Errorw("unable to save configuration to cache", "error", err)
	}

	zap.S().Debugw("configuration", "confi", confResponse)
	return confResponse, nil

	// }
	// return conf, nil
}

func (c *ConfigurationService) getConfiguration(ctx context.Context, device entity.Device) (entity.Configuration, error) {

	if device.Configuration != nil {
		return *device.Configuration, nil
	}

	// if device has no configuration look at the set of the device
	if device.SetID != nil {
		set, err := c.deviceReader.GetSet(ctx, *device.SetID)
		if err != nil {
			return entity.Configuration{}, err
		}
		if set.Configuration != nil {
			return *set.Configuration, nil
		}
	}

	// if the device has no set or the set has no configuration just grab the configuration from namespace
	namespace, err := c.deviceReader.GetNamespace(ctx, device.NamespaceID)
	if err != nil {
		return entity.Configuration{}, err
	}

	// namespace always has a configuration
	return namespace.Configuration, nil
}

func (c *ConfigurationService) getManifests(ctx context.Context, device entity.Device) ([]entity.ManifestWork, error) {
	getManifests := func(ctx context.Context, ids []string) ([]entity.ManifestWork, error) {
		manifests := make([]entity.ManifestWork, 0, len(device.ManifestIDS))
		for _, id := range ids {
			manifest, err := c.manifestReader.GetManifest(ctx, id)
			if err != nil {
				zap.S().Errorf("unable to get manifest", "error", err)
				continue
			}
			manifests = append(manifests, manifest)
		}
		return manifests, nil
	}

	if len(device.ManifestIDS) > 0 {
		return getManifests(ctx, device.ManifestIDS)
	}

	if device.SetID != nil {
		sets, err := c.deviceReader.GetSet(ctx, *device.SetID)
		if err != nil {
			return []entity.ManifestWork{}, err
		}
		if len(sets.ManifestIDS) > 0 {
			return getManifests(ctx, sets.ManifestIDS)
		}
	}

	namespace, err := c.deviceReader.GetNamespace(ctx, device.NamespaceID)
	if err != nil {
		return []entity.ManifestWork{}, err
	}
	return getManifests(ctx, namespace.ManifestIDS)
}
