package entity

type EnrolStatus int

const (
	EnroledStatus EnrolStatus = iota
	PendingEnrolStatus
	RefusedEnrolStatus
	NotEnroledStatus
)

type Device struct {
	// ID of the device
	ID string
	// EnrolStatus set to true if the device is enroled
	EnrolStatus EnrolStatus
	// Registred set to true if the device is already registered.
	Registred bool
	// Namespace in which the device is placed.
	Namespace string
	// List of sets in which the device is present
	Sets []string
	// Configuration of the device
	Configuration Configuration
	// List of workloads
	Workloads []Workload
	//HardwareInfo holds the information about the host's hardware
	HardwareInfo HardwareInfo
}

type Group struct {
	// Name of the group
	Name string
	// List of the devices in the group
	Devices []Device
}
