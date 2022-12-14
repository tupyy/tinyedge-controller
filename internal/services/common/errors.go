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
