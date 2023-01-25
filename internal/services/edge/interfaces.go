package edge

import (
	"context"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

// DeviceReader is an interface that groups all the methods allowing to query/get devices.
type DeviceReader interface {
	GetDevice(ctx context.Context, id string) (entity.Device, error)
}

// DeviceWriter allows creating a device.
type DeviceWriter interface {
	Create(ctx context.Context, device entity.Device) error
	Update(ctx context.Context, device entity.Device) error
}

//go:generate moq -out device_rw_moq.go . DeviceReaderWriter
type DeviceReaderWriter interface {
	DeviceReader
	DeviceWriter
}

//go:generate moq -out configuration_reader_moq.go . ConfigurationReader
type ConfigurationReader interface {
	GetDeviceConfiguration(ctx context.Context, id string) (entity.ConfigurationResponse, error)
}

//go:generate moq -out certficate_writer_moq.go . CertificateWriter
type CertificateWriter interface {
	SignCSR(ctx context.Context, csr []byte, cn string, ttl time.Duration) (entity.CertificateGroup, error)
}
