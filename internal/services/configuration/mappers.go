package configuration

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

func createConfigurationResponse(c entity.Configuration, manifests []entity.WorkloadManifest) entity.ConfigurationResponse {
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
		confResponse.Workloads = append(confResponse.Workloads, workloads(m)...)
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

func workloads(e entity.WorkloadManifest) []entity.Workload {
	workloads := make([]entity.Workload, 0, len(e.Pods))

	for i, p := range e.Pods {
		w := entity.Workload{
			Name:          fmt.Sprintf("%s-%d", p.Name, i),
			Specification: p.Spec,
			Rootless:      true,
			Labels:        e.Labels,
		}
		w.Hash = computeHash(w)
		w.ID = fmt.Sprintf("%s-%s", w.Name, w.Hash[:12])
		workloads = append(workloads, w)
	}

	for i := 0; i < len(workloads); i++ {
		w := &workloads[i]
		w.Configmaps = e.ConfigMaps
	}

	return workloads
}

func computeHash(w entity.Workload) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s", w.Name)

	for k, v := range w.Labels {
		fmt.Fprintf(&sb, "%s%s", k, v)
	}

	fmt.Fprintf(&sb, "%s", w)
	fmt.Fprintf(&sb, "%+v", w.WorkloadProfiles)
	fmt.Fprintf(&sb, "%v", w.Rootless)

	sum := sha256.Sum256(bytes.NewBufferString(sb.String()).Bytes())
	return fmt.Sprintf("%x", sum)
}
