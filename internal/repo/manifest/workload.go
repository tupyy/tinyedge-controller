package manifest

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	goyaml "github.com/go-yaml/yaml"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	apiv1 "github.com/tupyy/tinyedge-controller/pkg/api/v1"
)

// parse parses the manifest file and verify that all the resources defined are valid k8s manifestv1.
// Returns false if one resource is not a valid ConfigMap or Pod.
func parseManifestV1(content []byte) (entity.Manifest, error) {
	var workload apiv1.Workload

	if err := goyaml.Unmarshal(content, &workload); err != nil {
		return nil, err
	}

	e := entity.ManifestV1{
		TypeMeta: entity.TypeMeta{
			Version: entity.ManifestVersionV1,
		},
		ObjectMeta: entity.ObjectMeta{
			Labels: make(map[string]string),
		},
		Description: workload.Description,
		Selectors:   make([]entity.Selector, 0),
		Secrets:     make([]entity.Secret, 0, len(workload.Secrets)),
		Resources:   make([]string, 0, len(workload.Resources)),
	}

	for i := 0; true; i++ {
		keepGoing := false
		if i < len(workload.Selector.Namespaces) {
			e.Selectors = append(e.Selectors, entity.Selector{
				Type:  entity.NamespaceSelector,
				Value: workload.Selector.Namespaces[i],
			})
			keepGoing = true
		}
		if i < len(workload.Selector.Sets) {
			e.Selectors = append(e.Selectors, entity.Selector{
				Type:  entity.SetSelector,
				Value: workload.Selector.Sets[i],
			})
			keepGoing = true
		}
		if i < len(workload.Selector.Devices) {
			e.Selectors = append(e.Selectors, entity.Selector{
				Type:  entity.DeviceSelector,
				Value: workload.Selector.Devices[i],
			})
			keepGoing = true
		}
		if !keepGoing {
			break
		}
	}

	for _, s := range workload.Secrets {
		e.Secrets = append(e.Secrets, entity.Secret{
			Path: s.Path,
			Id:   s.Name,
			Key:  s.Key,
		})
	}

	for _, resource := range workload.Resources {
		e.Resources = append(e.Resources, resource.Ref)
	}

	return e, nil
}

func hash(data string) string {
	b := bytes.NewBufferString(data).Bytes()
	hash := sha256.New()
	hash.Write(b)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
