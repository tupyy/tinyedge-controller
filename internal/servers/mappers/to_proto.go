package mappers

import (
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
	edgepb "github.com/tupyy/tinyedge-controller/pkg/grpc/edge"
	"go.uber.org/zap"
	"sigs.k8s.io/yaml"
)

func MapConfigurationToProto(conf entity.ConfigurationResponse) *edgepb.ConfigurationResponse {
	response := &edgepb.ConfigurationResponse{
		Hash:      conf.Hash(),
		Workloads: make([]*common.Workload, 0, len(conf.Workloads)),
		Secrets:   make([]*common.Secret, 0, len(conf.Secrets)),
	}
	for _, w := range conf.Workloads {
		response.Workloads = append(response.Workloads, MapWorkloadToProto(w))
	}
	for _, s := range conf.Secrets {
		response.Secrets = append(response.Secrets, MapSecretToProto(s))
	}
	response.Configuration = &common.Configuration{
		HeartbeatPeriod: uint32(conf.Configuration.HeartbeatPeriod.Seconds()),
	}
	return response
}

func MapWorkloadToProto(w entity.Workload) *common.Workload {
	configmaps := make([]string, 0, len(w.Configmaps))
	for _, c := range w.Configmaps {
		data, err := yaml.Marshal(c)
		if err != nil {
			zap.S().Errorf("unable to marshal configmap: %w", err)
			continue
		}
		configmaps = append(configmaps, string(data))
	}

	specs, err := yaml.Marshal(w.Specification)
	if err != nil {
		zap.S().Errorf("unable to marshal spec: %w", err)
		return &common.Workload{}
	}

	pb := common.Workload{
		Name:       w.Name,
		Id:         w.ID,
		Hash:       w.Hash,
		ConfigMaps: configmaps,
		Rootless:   w.Rootless,
		Spec:       string(specs),
		Labels:     w.Labels,
	}

	return &pb
}

func MapSecretToProto(s entity.Secret) *common.Secret {
	return &common.Secret{}
}
