package configuration

import "github.com/tupyy/tinyedge-controller/internal/entity"

func createConfigurationResponse(c entity.Configuration, manifests []entity.ManifestWork) entity.ConfigurationResponse {
	confResponse := entity.ConfigurationResponse{
		Configuration: c,
		Workloads:     make([]entity.Workload, 0, len(manifests)),
		Secrets:       make([]entity.Secret, 0, 2*len(manifests)),
	}

	for _, m := range manifests {
		confResponse.Workloads = append(confResponse.Workloads, m.Workloads()...)
	}

	for _, m := range manifests {
		for _, s := range m.Secrets {
			confResponse.Secrets = append(confResponse.Secrets, s)
		}
	}
	return confResponse
}
