package entity

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
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

type ManifestReference struct {
	Id string
	// Valid is true if the content of the manifest is valid
	Valid bool
	// Hash of the manifest content
	Hash string
	// Repo - manifest's git repository
	Repo Repository
	// Path - filepath of the manifest in the local storage
	Path string
	// DeviceIDs - list of devices on which this manifest is applied.
	DeviceIDs []string
	// SetIDs - list of sets on which this manifest is applied.
	SetIDs []string
	// NamespaceIDs - list of namespaces on which this manifest is applied.
	NamespaceIDs []string
}

type ManifestWork struct {
	// Id - id of the manifest which is the hash of the filepath
	Id string
	// Version
	Version string
	// Name - name of the manifest
	Name string
	// Hash of the manifest
	Hash string
	// Description - description of the manifest
	Description string
	// Valid is true if the manifest content is valid
	Valid bool
	// path of the manifest file in the local repo
	Path string
	// Selectors list of selectors
	Selectors []Selector
	// Rootless - set the mode of podman execution: rootless or rootfull
	Rootless bool
	// Secrets - list of secrets without values. Values are retrieve from Vault.
	Secrets []Secret
	// Labels
	Labels map[string]string
	// Pods - list of pods
	Pods []v1.Pod
	// ConfigMaps -list of configmaps
	ConfigMaps []v1.ConfigMap
	// Reference
	Reference *ManifestReference
}

func (m ManifestWork) Workloads() []Workload {
	workloads := make([]Workload, 0, len(m.Pods))

	for i, p := range m.Pods {
		w := Workload{
			Name:          fmt.Sprintf("%s-%d", m.Name, i),
			Specification: p.Spec,
			Rootless:      true,
			Labels:        m.Labels,
		}
		w.Hash = m.computeHash(w)
		w.ID = fmt.Sprintf("%s-%s", w.Name, w.Hash[:12])
		workloads = append(workloads, w)
	}

	for i := 0; i < len(workloads); i++ {
		w := &workloads[i]
		w.Configmaps = m.ConfigMaps
	}

	return workloads
}

func (m ManifestWork) computeHash(w Workload) string {
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

type SelectorType int

const (
	NamespaceSelector SelectorType = iota
	SetSelector
	DeviceSelector
)

type Selector struct {
	Type  SelectorType
	Value string
}

type Secret struct {
	Id    string
	Path  string
	Key   string
	Hash  string
	Value string
}
