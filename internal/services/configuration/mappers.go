package configuration

import (
	"github.com/tupyy/tinyedge-controller/internal/entity"
)

func createConfigurationResponse(c entity.Configuration, manifests []entity.Workload) entity.ConfigurationResponse {
	confResponse := entity.ConfigurationResponse{
		Configuration: c,
	}

	return confResponse
}
