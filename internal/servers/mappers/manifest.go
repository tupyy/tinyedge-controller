package mappers

import (
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

func ManifestToProto(m entity.Manifest) *admin.Manifest {
	manifest := &admin.Manifest{
		Id:      m.GetID(),
		Version: m.GetVersion(),
		Hash:    m.GetHash(),
	}

	manifest.Selectors = make([]*admin.Selector, 0, len(m.GetSelectors()))
	for _, s := range m.GetSelectors() {
		var resourceType string
		switch s.Type {
		case entity.NamespaceSelector:
			resourceType = "namespace"
		case entity.SetSelector:
			resourceType = "set"
		case entity.DeviceSelector:
			resourceType = "device"
		}
		manifest.Selectors = append(manifest.Selectors, &admin.Selector{
			Value:        s.Value,
			ResourceType: resourceType,
		})
	}

	return manifest
}
