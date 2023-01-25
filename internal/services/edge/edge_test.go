package edge_test

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/edge"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
)

const (
	certificate = `
-----BEGIN CERTIFICATE-----
MIIDujCCAqKgAwIBAgIIE31FZVaPXTUwDQYJKoZIhvcNAQEFBQAwSTELMAkGA1UE
BhMCVVMxEzARBgNVBAoTCkdvb2dsZSBJbmMxJTAjBgNVBAMTHEdvb2dsZSBJbnRl
cm5ldCBBdXRob3JpdHkgRzIwHhcNMTQwMTI5MTMyNzQzWhcNMTQwNTI5MDAwMDAw
WjBpMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwN
TW91bnRhaW4gVmlldzETMBEGA1UECgwKR29vZ2xlIEluYzEYMBYGA1UEAwwPbWFp
bC5nb29nbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEfRrObuSW5T7q
5CnSEqefEmtH4CCv6+5EckuriNr1CjfVvqzwfAhopXkLrq45EQm8vkmf7W96XJhC
7ZM0dYi1/qOCAU8wggFLMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAa
BgNVHREEEzARgg9tYWlsLmdvb2dsZS5jb20wCwYDVR0PBAQDAgeAMGgGCCsGAQUF
BwEBBFwwWjArBggrBgEFBQcwAoYfaHR0cDovL3BraS5nb29nbGUuY29tL0dJQUcy
LmNydDArBggrBgEFBQcwAYYfaHR0cDovL2NsaWVudHMxLmdvb2dsZS5jb20vb2Nz
cDAdBgNVHQ4EFgQUiJxtimAuTfwb+aUtBn5UYKreKvMwDAYDVR0TAQH/BAIwADAf
BgNVHSMEGDAWgBRK3QYWG7z2aLV29YG2u2IaulqBLzAXBgNVHSAEEDAOMAwGCisG
AQQB1nkCBQEwMAYDVR0fBCkwJzAloCOgIYYfaHR0cDovL3BraS5nb29nbGUuY29t
L0dJQUcyLmNybDANBgkqhkiG9w0BAQUFAAOCAQEAH6RYHxHdcGpMpFE3oxDoFnP+
gtuBCHan2yE2GRbJ2Cw8Lw0MmuKqHlf9RSeYfd3BXeKkj1qO6TVKwCh+0HdZk283
TZZyzmEOyclm3UGFYe82P/iDFt+CeQ3NpmBg+GoaVCuWAARJN/KfglbLyyYygcQq
0SgeDh8dRKUiaW3HQSoYvTvdTuqzwK4CXsr3b5/dAOY8uMuG/IAR3FgwTbZ1dtoW
RvOTa8hYiU6A475WuZKyEHcwnGYe57u2I2KbMgcKjPniocj4QzgYsVAVKW3IwaOh
yE+vPxsiUkvQHdO2fojCkY8jg70jxM+gu59tPDNbw3Uh/2Ij310FgTHsnGQMyA==
-----END CERTIFICATE-----`
)

var _ = Describe("Edge", func() {
	var (
		configureReader *edge.ConfigurationReaderMock
		certWriter      *edge.CertificateWriterMock
	)

	Describe("Enrol", func() {
		It("when device is enroled", func() {
			deviceReadWriter := &edge.DeviceReaderWriterMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{}, errService.NewResourceNotFoundError("device", id)
				},
				CreateFunc: func(ctx context.Context, device entity.Device) error {
					return nil
				},
			}

			service := edge.New(deviceReadWriter, configureReader, certWriter)
			status, err := service.Enrol(context.TODO(), "deviceID")
			Expect(err).To(BeNil())
			Expect(status).To(Equal(entity.EnroledStatus))
			calls := deviceReadWriter.CreateCalls()
			Expect(len(calls)).To(Equal(1))

			firstCall := calls[0]
			Expect(firstCall.Device.ID).To(Equal("deviceID"))
			Expect(firstCall.Device.EnrolStatus).To(Equal(entity.EnroledStatus))
		})

		It("device is already enroled", func() {
			deviceReadWriter := &edge.DeviceReaderWriterMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{
						ID:          "deviceID",
						EnrolStatus: entity.EnroledStatus,
					}, nil
				},
				CreateFunc: func(ctx context.Context, device entity.Device) error {
					return nil
				},
			}

			service := edge.New(deviceReadWriter, configureReader, certWriter)
			status, err := service.Enrol(context.TODO(), "deviceID")
			Expect(err).To(BeNil())
			Expect(status).To(Equal(entity.EnroledStatus))
			calls := deviceReadWriter.CreateCalls()
			Expect(len(calls)).To(Equal(0))
		})

		It("get device end with error", func() {
			deviceReadWriter := &edge.DeviceReaderWriterMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{}, errors.New("unknown error")
				},
			}

			service := edge.New(deviceReadWriter, configureReader, certWriter)
			status, err := service.Enrol(context.TODO(), "deviceID")
			Expect(err).NotTo(BeNil())
			Expect(status).To(Equal(entity.NotEnroledStatus))
			calls := deviceReadWriter.CreateCalls()
			Expect(len(calls)).To(Equal(0))
		})

		It("create device returns error", func() {
			deviceReadWriter := &edge.DeviceReaderWriterMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{}, errService.NewResourceNotFoundError("device", "deviceID")
				},
				CreateFunc: func(ctx context.Context, device entity.Device) error {
					return errors.New("cannot create device")
				},
			}

			service := edge.New(deviceReadWriter, configureReader, certWriter)
			status, err := service.Enrol(context.TODO(), "deviceID")
			Expect(err).NotTo(BeNil())
			Expect(status).To(Equal(entity.NotEnroledStatus))
			calls := deviceReadWriter.CreateCalls()
			Expect(len(calls)).To(Equal(1))
		})
	})

	Describe("Register", func() {
		It("device is registered", func() {
			deviceRW := &edge.DeviceReaderWriterMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{
						ID:          "deviceID",
						EnrolStatus: entity.EnroledStatus,
					}, nil
				},
				UpdateFunc: func(ctx context.Context, device entity.Device) error {
					return nil
				},
			}
			certWriter := &edge.CertificateWriterMock{
				SignCSRFunc: func(ctx context.Context, csr []byte, cn string, ttl time.Duration) (entity.CertificateGroup, error) {
					block, _ := pem.Decode([]byte(certificate))
					if block == nil {
						panic("failed to parse certificate PEM")
					}
					cert, err := x509.ParseCertificate(block.Bytes)
					if err != nil {
						panic("failed to parse certificate: " + err.Error())
					}
					certificate := entity.CertificateGroup{
						Certificate: cert,
					}
					return certificate, nil
				},
			}
			service := edge.New(deviceRW, configureReader, certWriter)
			csr := "csr"
			certificate, err := service.Register(context.TODO(), "deviceID", csr)
			Expect(err).To(BeNil())
			calls := deviceRW.UpdateCalls()
			Expect(len(calls)).To(Equal(1))
			Expect(calls[0].Device.Registred).To(BeTrue())
			Expect(calls[0].Device.CertificateSerialNumber).To(Equal("137D4565568F5D35"))
			Expect(certificate.GetSerialNumber()).To(Equal("137D4565568F5D35"))
		})
		It("update device return error", func() {
			deviceRW := &edge.DeviceReaderWriterMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{
						ID:          "deviceID",
						EnrolStatus: entity.EnroledStatus,
					}, nil
				},
				UpdateFunc: func(ctx context.Context, device entity.Device) error {
					return errors.New("unknown error")
				},
			}
			certWriter := &edge.CertificateWriterMock{
				SignCSRFunc: func(ctx context.Context, csr []byte, cn string, ttl time.Duration) (entity.CertificateGroup, error) {
					block, _ := pem.Decode([]byte(certificate))
					if block == nil {
						panic("failed to parse certificate PEM")
					}
					cert, err := x509.ParseCertificate(block.Bytes)
					if err != nil {
						panic("failed to parse certificate: " + err.Error())
					}
					certificate := entity.CertificateGroup{
						Certificate: cert,
					}
					return certificate, nil
				},
			}
			service := edge.New(deviceRW, configureReader, certWriter)
			csr := "csr"
			_, err := service.Register(context.TODO(), "deviceID", csr)
			Expect(err).NotTo(BeNil())
			calls := deviceRW.UpdateCalls()
			Expect(len(calls)).To(Equal(1))
		})
		It("unable to sign the csr", func() {
			deviceRW := &edge.DeviceReaderWriterMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{
						ID:          "deviceID",
						EnrolStatus: entity.EnroledStatus,
					}, nil
				},
				UpdateFunc: func(ctx context.Context, device entity.Device) error {
					return errors.New("unknown error")
				},
			}
			certWriter := &edge.CertificateWriterMock{
				SignCSRFunc: func(ctx context.Context, csr []byte, cn string, ttl time.Duration) (entity.CertificateGroup, error) {
					return entity.CertificateGroup{}, errors.New("unknown error")
				},
			}
			service := edge.New(deviceRW, configureReader, certWriter)
			csr := "csr"
			_, err := service.Register(context.TODO(), "deviceID", csr)
			Expect(err).NotTo(BeNil())
			calls := deviceRW.UpdateCalls()
			Expect(len(calls)).To(Equal(0))
		})
		It("device not found", func() {
			deviceRW := &edge.DeviceReaderWriterMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{}, errors.New("device not found")
				},
				UpdateFunc: func(ctx context.Context, device entity.Device) error {
					return errors.New("unknown error")
				},
			}
			service := edge.New(deviceRW, configureReader, certWriter)
			csr := "csr"
			_, err := service.Register(context.TODO(), "deviceID", csr)
			Expect(err).NotTo(BeNil())
			calls := deviceRW.UpdateCalls()
			Expect(len(calls)).To(Equal(0))
		})
		It("device is not enroled", func() {
			deviceRW := &edge.DeviceReaderWriterMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{
						ID:          "deviceID",
						EnrolStatus: entity.NotEnroledStatus,
					}, nil
				},
			}
			service := edge.New(deviceRW, configureReader, certWriter)
			csr := "csr"
			_, err := service.Register(context.TODO(), "deviceID", csr)
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("IsRegistered", func() {
		It("device is registered", func() {
			deviceRW := &edge.DeviceReaderWriterMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{
						ID:          "deviceID",
						EnrolStatus: entity.EnroledStatus,
						Registred:   true,
					}, nil
				},
			}
			service := edge.New(deviceRW, configureReader, certWriter)
			isRegisterd, err := service.IsRegistered(context.TODO(), "deviceID")
			Expect(err).To(BeNil())
			Expect(isRegisterd).To(BeTrue())
		})
		It("device is not registered", func() {
			deviceRW := &edge.DeviceReaderWriterMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{
						ID:          "deviceID",
						EnrolStatus: entity.EnroledStatus,
						Registred:   false,
					}, nil
				},
			}
			service := edge.New(deviceRW, configureReader, certWriter)
			isRegisterd, err := service.IsRegistered(context.TODO(), "deviceID")
			Expect(err).To(BeNil())
			Expect(isRegisterd).To(BeFalse())
		})
	})
})
