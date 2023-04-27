// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package edge

import (
	"context"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"sync"
)

// Ensure, that ConfigurationReaderMock does implement ConfigurationReader.
// If this is not the case, regenerate this file with moq.
var _ ConfigurationReader = &ConfigurationReaderMock{}

// ConfigurationReaderMock is a mock implementation of ConfigurationReader.
//
// 	func TestSomethingThatUsesConfigurationReader(t *testing.T) {
//
// 		// make and configure a mocked ConfigurationReader
// 		mockedConfigurationReader := &ConfigurationReaderMock{
// 			GetDeviceConfigurationFunc: func(ctx context.Context, id string) (entity.ConfigurationResponse, error) {
// 				panic("mock out the GetDeviceConfiguration method")
// 			},
// 		}
//
// 		// use mockedConfigurationReader in code that requires ConfigurationReader
// 		// and then make assertions.
//
// 	}
type ConfigurationReaderMock struct {
	// GetDeviceConfigurationFunc mocks the GetDeviceConfiguration method.
	GetDeviceConfigurationFunc func(ctx context.Context, id string) (entity.DeviceConfiguration, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetDeviceConfiguration holds details about calls to the GetDeviceConfiguration method.
		GetDeviceConfiguration []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID string
		}
	}
	lockGetDeviceConfiguration sync.RWMutex
}

// GetDeviceConfiguration calls GetDeviceConfigurationFunc.
func (mock *ConfigurationReaderMock) GetDeviceConfiguration(ctx context.Context, id string) (entity.DeviceConfiguration, error) {
	if mock.GetDeviceConfigurationFunc == nil {
		panic("ConfigurationReaderMock.GetDeviceConfigurationFunc: method is nil but ConfigurationReader.GetDeviceConfiguration was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  string
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockGetDeviceConfiguration.Lock()
	mock.calls.GetDeviceConfiguration = append(mock.calls.GetDeviceConfiguration, callInfo)
	mock.lockGetDeviceConfiguration.Unlock()
	return mock.GetDeviceConfigurationFunc(ctx, id)
}

// GetDeviceConfigurationCalls gets all the calls that were made to GetDeviceConfiguration.
// Check the length with:
//     len(mockedConfigurationReader.GetDeviceConfigurationCalls())
func (mock *ConfigurationReaderMock) GetDeviceConfigurationCalls() []struct {
	Ctx context.Context
	ID  string
} {
	var calls []struct {
		Ctx context.Context
		ID  string
	}
	mock.lockGetDeviceConfiguration.RLock()
	calls = mock.calls.GetDeviceConfiguration
	mock.lockGetDeviceConfiguration.RUnlock()
	return calls
}
