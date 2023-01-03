package configuration

import (
	"context"
	"errors"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
)

type ConfigurationService struct {
	manifestReader  common.ManifestReader
	deviceReader    common.DeviceReader
	cacheReadWriter common.ConfigurationCacheReaderWriter
}

func New(deviceReader common.DeviceReader, manifestReader common.ManifestReader, cacheReaderWriter common.ConfigurationCacheReaderWriter) *ConfigurationService {
	return &ConfigurationService{
		manifestReader:  manifestReader,
		deviceReader:    deviceReader,
		cacheReadWriter: cacheReaderWriter,
	}
}

func (c *ConfigurationService) GetConfiguration(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error) {
	conf, err := c.cacheReadWriter.Get(ctx, deviceID)
	if err != nil {
		if !errors.Is(err, common.ErrResourceNotFound) {
			return entity.ConfigurationResponse{}, err
		}
		// create configuration from pg and save it to cache
	}
	return conf, nil
}

func (c *ConfigurationService) createConfiguration(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error) {
	device, err := c.deviceReader.GetDevice(ctx, deviceID)
	if err != nil {
		return entity.ConfigurationResponse{}, err
	}

	set, err := c.deviceReader.GetSet(ctx, id)
	if err != nil {
		return err
	}

	conf := set.Configuration
	if set.Configuration == nil {
		// get namespace's configuration
		namespace, err := c.deviceReader.GetNamespace(ctx, set.NamespaceID)
		if err != nil {
			return err
		}
		conf = &namespace.Configuration
	}

	manifests, err := c.manifestReader.GetSetManifests(ctx, id)
	if err != nil {
		return err
	}

	return c.cacheReadWriter.Put(ctx, id, createConfigurationResponse(*conf, manifests))
}
