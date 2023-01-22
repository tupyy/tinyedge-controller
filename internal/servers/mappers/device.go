package mappers

import (
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/common"
)

func NamespaceToProto(n entity.Namespace) *admin.Namespace {
	return &admin.Namespace{
		Id:            n.Name,
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
		Manifests: make([]string, 0, len(s.ManifestIDS)),
		Devices:   make([]string, 0, len(s.ManifestIDS)),
	}
	if s.Configuration != nil {
		set.Configuration = s.Configuration.ID
	}
	for _, m := range s.ManifestIDS {
		set.Manifests = append(set.Manifests, m)
	}
	for _, id := range s.DeviceIDs {
		set.Devices = append(set.Devices, id)
	}
	return set
}

func DeviceToProto(d entity.Device) *common.Device {
	dp := &common.Device{
		Id:            d.ID,
		Namespace:     d.NamespaceID,
		CertificateSn: d.CertificateSerialNumber,
		Manifests:     d.ManifestIDS,
		EnrolStatus:   d.EnrolStatus.String(),
		Registered:    d.Registred,
	}

	if d.EnrolStatus == entity.EnroledStatus {
		dp.EnroledAt = d.EnroledAt.Format(time.RFC3339)
	}

	if d.Registred {
		dp.RegisteredAt = d.RegisteredAt.Format(time.RFC3339)
	}

	if d.SetID != nil {
		dp.Set = *d.SetID
	}

	if d.Configuration != nil {
		dp.Configuration = d.Configuration.ID
	}

	return dp
}
