package entity

import (
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
	HeartbeatPeriod time.Duration `json:"heartbeat_period"`
}
