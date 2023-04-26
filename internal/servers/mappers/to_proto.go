package mappers

import (
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
	edgepb "github.com/tupyy/tinyedge-controller/pkg/grpc/edge"
)

func MapConfigurationToProto(conf entity.ConfigurationResponse) *edgepb.ConfigurationResponse {
	response := &edgepb.ConfigurationResponse{
		Hash: conf.Hash,
	}
	response.Configuration = &common.Configuration{
		HeartbeatPeriod: uint32(conf.Configuration.HeartbeatPeriod.Seconds()),
	}
	return response
}
