package entity

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	goyaml "github.com/go-yaml/yaml"
)

const (
	// name of the manifest work
	ManifestWorkFilename = "manifest_work.yaml"
)

// Repository holds the information about the git repository where the ManifestWork are to be found.
type Repository struct {
	Id             string
	Url            string
	Branch         string
	LocalPath      string
	CurrentHeadSha string
	TargetHeadSha  string
	PullPeriod     time.Duration
}

type ManifestWorkV1 struct {
	Id           string
	Version      string           `yaml:"version"`
	Name         string           `yaml:"name"`
	Description  string           `yaml:"description"`
	Selector     Selector         `yaml:"selectors"`
	Secrets      []ManifestSecret `yaml:"secrets"`
	Resources    []Resource       `yaml:"resources"`
	Repo         Repository       `yaml:"-"`
	Path         string           `yaml:"-"`
	DeviceIDs    []string         `yaml:"-"`
	SetIDs       []string         `yaml:"-"`
	NamespaceIDs []string         `yaml:"-"`
}

func (m ManifestWorkV1) Hash() string {
	hash := sha256.New()
	content, _ := goyaml.Marshal(m)
	hash.Write(content)
	return string(hash.Sum(nil))
}

func (m ManifestWorkV1) Encode() string {
	content, _ := goyaml.Marshal(m)
	return base64.StdEncoding.EncodeToString(content)
}

func (m ManifestWorkV1) Decode(str string) (ManifestWorkV1, error) {
	content, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ManifestWorkV1{}, err
	}
	var mm ManifestWorkV1
	err = goyaml.Unmarshal(content, &mm)
	if err != nil {
		return ManifestWorkV1{}, err
	}
	return mm, nil
}

type Relation[T any] struct {
	ObjectType T
	ManifestID string
	ObjectID   string
}

type Selector struct {
	Namespaces []string `yaml:"namespaces"`
	Sets       []string `yaml:"sets"`
	Devices    []string `yaml:"devices"`
}

type ManifestSecret struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Key  string `yaml:"key"`
}

type Resource struct {
	Ref          string `yaml:"$ref"`
	Kind         string
	KubeResource interface{}
}
