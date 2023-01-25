// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package reference

import (
	"context"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"sync"
)

// Ensure, that DeviceReaderMock does implement DeviceReader.
// If this is not the case, regenerate this file with moq.
var _ DeviceReader = &DeviceReaderMock{}

// DeviceReaderMock is a mock implementation of DeviceReader.
//
//	func TestSomethingThatUsesDeviceReader(t *testing.T) {
//
//		// make and configure a mocked DeviceReader
//		mockedDeviceReader := &DeviceReaderMock{
//			GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
//				panic("mock out the GetDevice method")
//			},
//			GetNamespaceFunc: func(ctx context.Context, id string) (entity.Namespace, error) {
//				panic("mock out the GetNamespace method")
//			},
//			GetSetFunc: func(ctx context.Context, id string) (entity.Set, error) {
//				panic("mock out the GetSet method")
//			},
//		}
//
//		// use mockedDeviceReader in code that requires DeviceReader
//		// and then make assertions.
//
//	}
type DeviceReaderMock struct {
	// GetDeviceFunc mocks the GetDevice method.
	GetDeviceFunc func(ctx context.Context, id string) (entity.Device, error)

	// GetNamespaceFunc mocks the GetNamespace method.
	GetNamespaceFunc func(ctx context.Context, id string) (entity.Namespace, error)

	// GetSetFunc mocks the GetSet method.
	GetSetFunc func(ctx context.Context, id string) (entity.Set, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetDevice holds details about calls to the GetDevice method.
		GetDevice []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID string
		}
		// GetNamespace holds details about calls to the GetNamespace method.
		GetNamespace []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID string
		}
		// GetSet holds details about calls to the GetSet method.
		GetSet []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID string
		}
	}
	lockGetDevice    sync.RWMutex
	lockGetNamespace sync.RWMutex
	lockGetSet       sync.RWMutex
}

// GetDevice calls GetDeviceFunc.
func (mock *DeviceReaderMock) GetDevice(ctx context.Context, id string) (entity.Device, error) {
	if mock.GetDeviceFunc == nil {
		panic("DeviceReaderMock.GetDeviceFunc: method is nil but DeviceReader.GetDevice was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  string
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockGetDevice.Lock()
	mock.calls.GetDevice = append(mock.calls.GetDevice, callInfo)
	mock.lockGetDevice.Unlock()
	return mock.GetDeviceFunc(ctx, id)
}

// GetDeviceCalls gets all the calls that were made to GetDevice.
// Check the length with:
//
//	len(mockedDeviceReader.GetDeviceCalls())
func (mock *DeviceReaderMock) GetDeviceCalls() []struct {
	Ctx context.Context
	ID  string
} {
	var calls []struct {
		Ctx context.Context
		ID  string
	}
	mock.lockGetDevice.RLock()
	calls = mock.calls.GetDevice
	mock.lockGetDevice.RUnlock()
	return calls
}

// GetNamespace calls GetNamespaceFunc.
func (mock *DeviceReaderMock) GetNamespace(ctx context.Context, id string) (entity.Namespace, error) {
	if mock.GetNamespaceFunc == nil {
		panic("DeviceReaderMock.GetNamespaceFunc: method is nil but DeviceReader.GetNamespace was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  string
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockGetNamespace.Lock()
	mock.calls.GetNamespace = append(mock.calls.GetNamespace, callInfo)
	mock.lockGetNamespace.Unlock()
	return mock.GetNamespaceFunc(ctx, id)
}

// GetNamespaceCalls gets all the calls that were made to GetNamespace.
// Check the length with:
//
//	len(mockedDeviceReader.GetNamespaceCalls())
func (mock *DeviceReaderMock) GetNamespaceCalls() []struct {
	Ctx context.Context
	ID  string
} {
	var calls []struct {
		Ctx context.Context
		ID  string
	}
	mock.lockGetNamespace.RLock()
	calls = mock.calls.GetNamespace
	mock.lockGetNamespace.RUnlock()
	return calls
}

// GetSet calls GetSetFunc.
func (mock *DeviceReaderMock) GetSet(ctx context.Context, id string) (entity.Set, error) {
	if mock.GetSetFunc == nil {
		panic("DeviceReaderMock.GetSetFunc: method is nil but DeviceReader.GetSet was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  string
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockGetSet.Lock()
	mock.calls.GetSet = append(mock.calls.GetSet, callInfo)
	mock.lockGetSet.Unlock()
	return mock.GetSetFunc(ctx, id)
}

// GetSetCalls gets all the calls that were made to GetSet.
// Check the length with:
//
//	len(mockedDeviceReader.GetSetCalls())
func (mock *DeviceReaderMock) GetSetCalls() []struct {
	Ctx context.Context
	ID  string
} {
	var calls []struct {
		Ctx context.Context
		ID  string
	}
	mock.lockGetSet.RLock()
	calls = mock.calls.GetSet
	mock.lockGetSet.RUnlock()
	return calls
}
