package entity

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
)

type Workload struct {
	Name string
	// Rootless is true if workload is to be executed in podman rootless
	Rootless bool
	// Configmaps
	Configmaps []string
	// secrets
	Secrets map[string]string
	// Workload labels
	Labels map[string]string
	// Workload profiles
	WorkloadProfiles []WorkloadProfile
	// specification
	Specification v1.PodSpec
}

func (p Workload) ID() string {
	return fmt.Sprintf("%s-%s", p.Name, p.Hash()[:12])
}

func (p Workload) String() string {
	json, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(json)
}

func (p Workload) Hash() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s", p.Name)
	for k, v := range p.Secrets {
		fmt.Fprintf(&sb, "%s%s", k, v)
	}

	for k, v := range p.Labels {
		fmt.Fprintf(&sb, "%s%s", k, v)
	}

	fmt.Fprintf(&sb, "%s", p.Specification.String())
	fmt.Fprintf(&sb, "%+v", p.WorkloadProfiles)
	fmt.Fprintf(&sb, "%v", p.Rootless)

	sum := sha256.Sum256(bytes.NewBufferString(sb.String()).Bytes())
	return fmt.Sprintf("%x", sum)
}

type WorkloadProfile struct {
	Name       string   `json:"name"`
	Conditions []string `json:"conditions"`
}
