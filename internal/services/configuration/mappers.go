package configuration

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

func createConfigurationResponse(c entity.Configuration, manifests []entity.ManifestWork) entity.ConfigurationResponse {
	hash := sha256.New()

	confResponse := entity.ConfigurationResponse{
		Configuration: c,
		Workloads:     make([]entity.Workload, 0, len(manifests)),
		Secrets:       make([]entity.Secret, 0, 2*len(manifests)),
	}

	data, err := json.Marshal(c)
	if err == nil {
		hash.Write(data)
	}

	for _, m := range manifests {
		confResponse.Workloads = append(confResponse.Workloads, m.Workloads()...)
		hash.Write([]byte(m.Hash))
	}

	for _, m := range manifests {
		for _, s := range m.Secrets {
			confResponse.Secrets = append(confResponse.Secrets, s)
		}
	}

	confResponse.Hash = fmt.Sprintf("%x", hash.Sum(nil))
	return confResponse
}
