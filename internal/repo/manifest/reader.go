package manifest

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	goyaml "github.com/go-yaml/yaml"
	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type ManifestReader func(r io.Reader, transformFn ...func(entity.Manifest) entity.Manifest) (entity.Manifest, error)

// ReadManifest parses the content of a reader and return a Manifest or error.
// TODO find a better way to set the id
func ReadManifest(r io.Reader, transformFn ...func(entity.Manifest) entity.Manifest) (entity.Manifest, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	kind, err := getKind(content)
	if err != nil {
		return nil, err
	}

	var manifest entity.Manifest
	switch kind {
	case entity.WorkloadManifestKind:
		manifest, err = parseWorkloadManifest(content)
	case entity.ConfigurationManifestKind:
		manifest, err = parseConfigurationManifest(content)
	}

	if err != nil {
		return nil, err
	}

	for _, fn := range transformFn {
		manifest = fn(manifest)
	}

	return manifest, nil
}

func getKind(content []byte) (entity.ManifestKind, error) {
	type anonymousStruct struct {
		Kind string `yaml:"kind"`
	}
	var a anonymousStruct
	if err := goyaml.Unmarshal(content, &a); err != nil {
		return 0, fmt.Errorf("unknown struct: %s", err)
	}
	switch strings.ToLower(a.Kind) {
	case "workload":
		return entity.WorkloadManifestKind, nil
	case "configuration":
		return entity.ConfigurationManifestKind, nil
	}
	return 0, fmt.Errorf("unknown kind: %q", a.Kind)
}
