package entity

import (
	"encoding/json"

	v1 "k8s.io/api/core/v1"
)

type Workload struct {
	// ID of the workload
	ID string
	// Name of the workload
	Name string
	// hash of the workload
	Hash string
	// Rootless is true if workload is to be executed in podman rootless
	Rootless bool
	// Configmaps
	Configmaps []v1.ConfigMap
	// Workload labels
	Labels map[string]string
	// Workload profiles
	WorkloadProfiles []WorkloadProfile
	// specification
	Specification v1.PodSpec
}

func (w Workload) String() string {
	json, err := json.Marshal(w)
	if err != nil {
		return err.Error()
	}
	return string(json)
}

type WorkloadProfile struct {
	Name       string   `json:"name"`
	Conditions []string `json:"conditions"`
}
