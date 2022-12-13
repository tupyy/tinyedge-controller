package entity

type HardwareInfo struct {
	Hostname      string
	OsInformation OsInformation
	Interfaces    []Interface
	SystemVendor  SystemVendor
}
type OsInformation struct {
	CommitID string
}

type SystemVendor struct {
	// manufacturer
	Manufacturer string
	// product name
	ProductName string
	// serial number
	SerialNumber string
	// Whether the machine appears to be a virtual machine or not
	Virtual bool
}

type Interface struct {
	// name
	Name string
	// HasCarrier
	HasCarrier bool
	// ipv4 address
	IPV4Address []string
	// mac address
	MacAddress string
}
