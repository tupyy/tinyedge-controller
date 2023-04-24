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
		Configuration: n.Configuration.GetID(),
		Devices:       n.Devices,
		Sets:          n.Sets,
	}
}

func SetToProto(s entity.Set) *common.Set {
	set := &common.Set{
		Name:      s.Name,
		Namespace: s.NamespaceID,
		Manifests: make([]string, 0, len(s.Workloads)),
		Devices:   make([]string, 0, len(s.Workloads)),
	}
	if s.Configuration != nil {
		set.Configuration = s.Configuration.GetID()
	}
	for _, m := range s.Workloads {
		set.Manifests = append(set.Manifests, m.GetID())
	}
	for _, id := range s.Devices {
		set.Devices = append(set.Devices, id)
	}
	return set
}

func DeviceToProto(d entity.Device) *common.Device {
	dp := &common.Device{
		Id:            d.ID,
		Namespace:     d.NamespaceID,
		CertificateSn: d.CertificateSerialNumber,
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
		dp.Configuration = d.Configuration.GetID()
	}

	return dp
}
