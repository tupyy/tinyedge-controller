package mappers

import (
	"encoding/json"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
	"go.uber.org/zap"
)

func ManifestToProto(m entity.Manifest) *admin.Manifest {
	manifest := &admin.Manifest{
		Id:          m.GetID(),
		Version:     m.GetVersion(),
		Hash:        m.GetHash(),
		Description: m.Description,
		Valid:       m.Valid,
		Path:        m.Path,
		Rootless:    m.Rootless,
		Devices:     m.Reference.DeviceIDs,
		Namespaces:  m.Reference.NamespaceIDs,
		Sets:        m.Reference.SetIDs,
	}

	manifest.Selectors = make([]*admin.Selector, 0, len(m.Selectors))
	for _, s := range m.Selectors {
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

	manifest.Secrets = make([]string, 0, len(m.Secrets))
	for _, secret := range m.Secrets {
		manifest.Secrets = append(manifest.Secrets, secret.Id)
	}

	manifest.Pods = make([]string, 0, len(m.Pods))
	for _, pod := range m.Pods {
		data, err := json.Marshal(pod)
		if err != nil {
			zap.S().Warnw("unable to marshal pod", "error", err, "pod", pod)
		}
		manifest.Pods = append(manifest.Pods, string(data))
	}

	manifest.Configmaps = make([]string, 0, len(m.ConfigMaps))
	for _, c := range m.ConfigMaps {
		data, err := json.Marshal(c)
		if err != nil {
			zap.S().Warnw("unable to marshal configmap", "error", err, "pod", c)
			continue
		}
		manifest.Configmaps = append(manifest.Configmaps, string(data))
	}
	return manifest
}
