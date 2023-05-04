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

	version, err := getVersion(content)
	if err != nil {
		return nil, err
	}

	var manifest entity.Manifest
	switch version {
	case entity.ManifestVersionV1:
		manifest, err = parseManifestV1(content)
	}

	if err != nil {
		return nil, err
	}

	for _, fn := range transformFn {
		manifest = fn(manifest)
	}

	return manifest, nil
}

func getVersion(content []byte) (entity.Version, error) {
	type anonymousStruct struct {
		Version string `yaml:"version"`
	}
	var a anonymousStruct
	if err := goyaml.Unmarshal(content, &a); err != nil {
		return entity.ManifestUnknownVersion, fmt.Errorf("unknown struct: %s", err)
	}
	switch strings.ToLower(a.Version) {
	case "v1":
		return entity.ManifestVersionV1, nil
	default:
		return entity.ManifestUnknownVersion, nil
	}
}
