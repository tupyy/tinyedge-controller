package configuration

import (
	"github.com/tupyy/tinyedge-controller/internal/entity"
)

func createConfigurationResponse(c entity.Configuration, manifests []entity.ManifestV1) entity.DeviceConfiguration {
	confResponse := entity.DeviceConfiguration{
		Configuration: c,
	}

	return confResponse
}
