package errors

import (
	"fmt"
)

type DeviceNotRegisteredError struct {
	DeviceID string
}

func (d DeviceNotRegisteredError) Error() string {
	return fmt.Sprintf("device %q is not registered", d.DeviceID)
}

func NewDeviceNotRegisteredError(deviceID string) DeviceNotRegisteredError {
	return DeviceNotRegisteredError{deviceID}
}

type DeviceNotEnroledError struct {
	DeviceID string
}

func (d DeviceNotEnroledError) Error() string {
	return fmt.Sprintf("device %q is not enroled", d.DeviceID)
}

func NewDeviceNotEnroledError(deviceID string) DeviceNotEnroledError {
	return DeviceNotEnroledError{deviceID}
}

func IsResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(ResourseNotFoundError)
	return ok
}

func IsResourceAlreadyExists(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(ResourceAlreadyExists)
	return ok
}

type ResourseNotFoundError struct {
	ResourceType string
	Reason       string
	ResourceID   string
	Err          error
}

func (r ResourseNotFoundError) Error() string {
	if r.Reason != "" {
		return r.Reason
	}
	return fmt.Sprintf("%s %q not found", r.ResourceType, r.ResourceID)
}

type ResourceAlreadyExists struct {
	ResourceType string
	ResourceID   string
	Err          error
}

func NewResourceNotFoundErrorWithErr(resourceType, resourceID string, err error) ResourseNotFoundError {
	return ResourseNotFoundError{ResourceType: resourceType, ResourceID: resourceID, Err: err}
}

func NewResourceNotFoundError(resourceType, resourceID string) ResourseNotFoundError {
	return ResourseNotFoundError{ResourceType: resourceType, ResourceID: resourceID, Err: fmt.Errorf("resource not found")}
}

func NewResourceNotFoundErrorWithReason(reason string) ResourseNotFoundError {
	return ResourseNotFoundError{Reason: reason, Err: fmt.Errorf("resource not found")}
}
func (r ResourceAlreadyExists) Error() string {
	return fmt.Sprintf("%s %q already exists", r.ResourceType, r.ResourceID)
}

func NewResourceAlreadyExistsError(resourceType, resourceID string) ResourceAlreadyExists {
	return ResourceAlreadyExists{resourceType, resourceID, fmt.Errorf("resource already exists")}
}

func NewResourceAlreadyExistsErrorWithErr(resourceType, resourceID string, err error) ResourceAlreadyExists {
	return ResourceAlreadyExists{resourceType, resourceID, err}
}

type PosgresNotAvailableError struct {
	Err error
}

func (p PosgresNotAvailableError) Error() string {
	return "postgres not available"
}

func NewPostgresNotAvailableError(service string) PosgresNotAvailableError {
	return PosgresNotAvailableError{fmt.Errorf("postgres not available in %q", service)}
}

type DeleteResourceError struct {
	ResourceType string
	ResourceID   string
	Reason       string
}

func (d DeleteResourceError) Error() string {
	return fmt.Sprintf("%s %q cannot be deleted: %s", d.ResourceType, d.ResourceID, d.Reason)
}

func NewDeleteResourceError(resourceType, resourceID, reason string) DeleteResourceError {
	return DeleteResourceError{resourceType, resourceID, reason}
}
