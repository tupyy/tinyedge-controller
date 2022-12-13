package entity

import (
	"encoding/json"
	"time"
)

type ConfigurationResponse struct {
	Configuration Configuration
	Workloads     []Workload
	Secrets       []Secret
}

func (c ConfigurationResponse) Hash() string {
	return "hash"
}

type Configuration struct {
	// list of profiles
	Profiles []Profile `json:"profiles"`
	// HeartbeatConfiguration hold the configuration of hearbeat
	HeartbeatConfiguration HeartbeatConfiguration `json:"heartbeat"`
}

type HeartbeatConfiguration struct {
	HardwareProfile HardwareProfileConfiguration
	// period in seconds
	Period time.Duration
}

func (h HeartbeatConfiguration) String() string {
	json, err := json.Marshal(h)
	if err != nil {
		return err.Error()
	}
	return string(json)
}

type Scope int

func (s Scope) String() string {
	switch s {
	case FullScope:
		return "full_scope"
	case DeltaScope:
		return "delta"
	default:
		return "unknown"
	}
}

const (
	FullScope Scope = iota
	DeltaScope
)

type HardwareProfileConfiguration struct {
	Include bool
	Scope   Scope
}
