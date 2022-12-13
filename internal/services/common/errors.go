package common

import "errors"

var (
	ErrDeviceNotFound      = errors.New("device not found")
	ErrDeviceNotRegistered = errors.New("device is not registered")
)

type ServiceError struct {
	error
}
