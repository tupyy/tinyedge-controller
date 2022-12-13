package mappers

import (
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
	edgepb "github.com/tupyy/tinyedge-controller/pkg/grpc/edge"
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
	return response
}

func MapWorkloadToProto(w entity.Workload) *common.Workload {
	return &common.Workload{}
}

func MapSecretToProto(s entity.Secret) *common.Secret {
	return &common.Secret{}
}
