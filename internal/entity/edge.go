package entity

// DeviceConfiguration is the entity which maps the response to the device following the GetConfiguration call.
type DeviceConfiguration struct {
	Hash          string
	Configuration Configuration
	Workload      []byte
}

type Heartbeat struct{}
