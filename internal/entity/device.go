package entity

import (
	"time"
)

type EnrolStatus int

func (e EnrolStatus) String() string {
	switch e {
	case EnroledStatus:
		return "enroled"
	case PendingEnrolStatus:
		return "pending"
	case RefusedEnrolStatus:
		return "refused"
	default:
		return "not_enroled"
	}
}

func (e EnrolStatus) FromString(s string) EnrolStatus {
	switch s {
	case "enroled":
		return EnroledStatus
	case "pending":
		return PendingEnrolStatus
	case "refused":
		return RefusedEnrolStatus
	default:
		return NotEnroledStatus
	}
}

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
	// EnroledAt represents the time when the device was enroled
	EnroledAt time.Time
	// Registred set to true if the device is already registered.
	Registred bool
	// RegisteredAt represents the time when the device registered.
	RegisteredAt time.Time
	// Namespace in which the device is placed.
	Namespace string
	// CertificateSerialNumber holds the SN of the certificate used for authorization.
	// This is the certificate generate at registration time.
	CertificateSerialNumber string
	// ID of set in which the device is present
	Set string
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
