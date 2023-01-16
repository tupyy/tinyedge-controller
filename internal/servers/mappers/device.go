package mappers

import (
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
)

func NamespaceToProto(n entity.Namespace) *admin.Namespace {
	return &admin.Namespace{
		Name:          n.Name,
		IsDefault:     n.IsDefault,
		Configuration: n.Configuration.ID,
		Devices:       n.DeviceIDs,
		Sets:          n.SetIDs,
		Manifests:     n.ManifestIDS,
	}
}

func SetToProto(s entity.Set) *common.Set {
	set := &common.Set{
		Name:      s.Name,
		Namespace: s.NamespaceID,
		Manifests: make([]string, 0),
	}
	if s.Configuration != nil {
		set.Configuration = s.Configuration.ID
	}
	for _, m := range s.ManifestIDS {
		set.Manifests = append(set.Manifests, m)
	}
	return set
}
