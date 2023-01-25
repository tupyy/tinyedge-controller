// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package edge

import (
	"context"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"sync"
	"time"
)

// Ensure, that CertificateWriterMock does implement CertificateWriter.
// If this is not the case, regenerate this file with moq.
var _ CertificateWriter = &CertificateWriterMock{}

// CertificateWriterMock is a mock implementation of CertificateWriter.
//
//	func TestSomethingThatUsesCertificateWriter(t *testing.T) {
//
//		// make and configure a mocked CertificateWriter
//		mockedCertificateWriter := &CertificateWriterMock{
//			SignCSRFunc: func(ctx context.Context, csr []byte, cn string, ttl time.Duration) (entity.CertificateGroup, error) {
//				panic("mock out the SignCSR method")
//			},
//		}
//
//		// use mockedCertificateWriter in code that requires CertificateWriter
//		// and then make assertions.
//
//	}
type CertificateWriterMock struct {
	// SignCSRFunc mocks the SignCSR method.
	SignCSRFunc func(ctx context.Context, csr []byte, cn string, ttl time.Duration) (entity.CertificateGroup, error)

	// calls tracks calls to the methods.
	calls struct {
		// SignCSR holds details about calls to the SignCSR method.
		SignCSR []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Csr is the csr argument value.
			Csr []byte
			// Cn is the cn argument value.
			Cn string
			// TTL is the ttl argument value.
			TTL time.Duration
		}
	}
	lockSignCSR sync.RWMutex
}

// SignCSR calls SignCSRFunc.
func (mock *CertificateWriterMock) SignCSR(ctx context.Context, csr []byte, cn string, ttl time.Duration) (entity.CertificateGroup, error) {
	if mock.SignCSRFunc == nil {
		panic("CertificateWriterMock.SignCSRFunc: method is nil but CertificateWriter.SignCSR was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Csr []byte
		Cn  string
		TTL time.Duration
	}{
		Ctx: ctx,
		Csr: csr,
		Cn:  cn,
		TTL: ttl,
	}
	mock.lockSignCSR.Lock()
	mock.calls.SignCSR = append(mock.calls.SignCSR, callInfo)
	mock.lockSignCSR.Unlock()
	return mock.SignCSRFunc(ctx, csr, cn, ttl)
}

// SignCSRCalls gets all the calls that were made to SignCSR.
// Check the length with:
//
//	len(mockedCertificateWriter.SignCSRCalls())
func (mock *CertificateWriterMock) SignCSRCalls() []struct {
	Ctx context.Context
	Csr []byte
	Cn  string
	TTL time.Duration
} {
	var calls []struct {
		Ctx context.Context
		Csr []byte
		Cn  string
		TTL time.Duration
	}
	mock.lockSignCSR.RLock()
	calls = mock.calls.SignCSR
	mock.lockSignCSR.RUnlock()
	return calls
}
