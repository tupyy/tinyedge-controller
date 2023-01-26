// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package manifest

import (
	"context"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"sync"
)

// Ensure, that ReferenceReaderMock does implement ReferenceReader.
// If this is not the case, regenerate this file with moq.
var _ ReferenceReader = &ReferenceReaderMock{}

// ReferenceReaderMock is a mock implementation of ReferenceReader.
//
//	func TestSomethingThatUsesReferenceReader(t *testing.T) {
//
//		// make and configure a mocked ReferenceReader
//		mockedReferenceReader := &ReferenceReaderMock{
//			GetRepositoryReferencesFunc: func(ctx context.Context, repo entity.Repository) ([]entity.ManifestReference, error) {
//				panic("mock out the GetRepositoryReferences method")
//			},
//		}
//
//		// use mockedReferenceReader in code that requires ReferenceReader
//		// and then make assertions.
//
//	}
type ReferenceReaderMock struct {
	// GetRepositoryReferencesFunc mocks the GetRepositoryReferences method.
	GetRepositoryReferencesFunc func(ctx context.Context, repo entity.Repository) ([]entity.ManifestReference, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetRepositoryReferences holds details about calls to the GetRepositoryReferences method.
		GetRepositoryReferences []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Repo is the repo argument value.
			Repo entity.Repository
		}
	}
	lockGetRepositoryReferences sync.RWMutex
}

// GetRepositoryReferences calls GetRepositoryReferencesFunc.
func (mock *ReferenceReaderMock) GetRepositoryReferences(ctx context.Context, repo entity.Repository) ([]entity.ManifestReference, error) {
	if mock.GetRepositoryReferencesFunc == nil {
		panic("ReferenceReaderMock.GetRepositoryReferencesFunc: method is nil but ReferenceReader.GetRepositoryReferences was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Repo entity.Repository
	}{
		Ctx:  ctx,
		Repo: repo,
	}
	mock.lockGetRepositoryReferences.Lock()
	mock.calls.GetRepositoryReferences = append(mock.calls.GetRepositoryReferences, callInfo)
	mock.lockGetRepositoryReferences.Unlock()
	return mock.GetRepositoryReferencesFunc(ctx, repo)
}

// GetRepositoryReferencesCalls gets all the calls that were made to GetRepositoryReferences.
// Check the length with:
//
//	len(mockedReferenceReader.GetRepositoryReferencesCalls())
func (mock *ReferenceReaderMock) GetRepositoryReferencesCalls() []struct {
	Ctx  context.Context
	Repo entity.Repository
} {
	var calls []struct {
		Ctx  context.Context
		Repo entity.Repository
	}
	mock.lockGetRepositoryReferences.RLock()
	calls = mock.calls.GetRepositoryReferences
	mock.lockGetRepositoryReferences.RUnlock()
	return calls
}
