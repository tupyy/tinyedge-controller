package common

import "errors"

var (
	ErrResourceNotFound      = errors.New("resource not found")
	ErrResourceAlreadyExists = errors.New("resources already exists")
	ErrDeviceNotRegistered   = errors.New("device is not registered")
	ErrDeviceNotEnroled      = errors.New("device is not enroled")
	ErrCertificateNotFound   = errors.New("certificate not found")

	ErrPostgresNotAvailable = errors.New("pg not available")
)

func IsResourceNotFound(err error) bool {
	return err != nil && errors.Is(err, ErrResourceNotFound)
}

func IsResourceAlreadyExists(err error) bool {
	return err != nil && errors.Is(err, ErrResourceAlreadyExists)
}
