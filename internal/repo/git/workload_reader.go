package git

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"path"

	goyaml "github.com/go-yaml/yaml"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	manifestv1 "github.com/tupyy/tinyedge-controller/internal/repo/models/manifest/v1"
)

// parse parses the manifest file and verify that all the resources defined are valid k8s manifestv1.
// Returns false if one resource is not a valid ConfigMap or Pod.
func parseWorkloadManifest(content []byte, filename, basePath string) (entity.Manifest, error) {
	var manifest manifestv1.Manifest

	if err := goyaml.Unmarshal(content, &manifest); err != nil {
		return nil, err
	}

	e := entity.Workload{
		TypeMeta: entity.TypeMeta{
			Version: manifest.Version,
			Kind:    entity.WorkloadManifestKind,
		},
		ObjectMeta: entity.ObjectMeta{
			Name:   manifest.Name,
			Id:     hash(path.Join(basePath, filename)),
			Labels: make(map[string]string),
		},
		Description: manifest.Description,
		Selectors:   make([]entity.Selector, 0),
		Secrets:     make([]entity.Secret, 0, len(manifest.Secrets)),
		Resources:   make([]string, 0, len(manifest.Resources)),
	}

	for i := 0; true; i++ {
		keepGoing := false
		if i < len(manifest.Selector.Namespaces) {
			e.Selectors = append(e.Selectors, entity.Selector{
				Type:  entity.NamespaceSelector,
				Value: manifest.Selector.Namespaces[i],
			})
			keepGoing = true
		}
		if i < len(manifest.Selector.Sets) {
			e.Selectors = append(e.Selectors, entity.Selector{
				Type:  entity.SetSelector,
				Value: manifest.Selector.Sets[i],
			})
			keepGoing = true
		}
		if i < len(manifest.Selector.Devices) {
			e.Selectors = append(e.Selectors, entity.Selector{
				Type:  entity.DeviceSelector,
				Value: manifest.Selector.Devices[i],
			})
			keepGoing = true
		}
		if !keepGoing {
			break
		}
	}

	for _, s := range manifest.Secrets {
		e.Secrets = append(e.Secrets, entity.Secret{
			Path: s.Path,
			Id:   s.Name,
			Key:  s.Key,
		})
	}

	for _, resource := range manifest.Resources {
		e.Resources = append(e.Resources, path.Join(basePath, resource.Ref))
	}

	return e, nil
}

func hash(data string) string {
	b := bytes.NewBufferString(data).Bytes()
	hash := sha256.New()
	hash.Write(b)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
