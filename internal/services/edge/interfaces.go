package edge

import (
	"context"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

// DeviceReader is an interface that groups all the methods allowing to query/get devices.
type DeviceReader interface {
	Get(ctx context.Context, id string) (entity.Device, error)
}

// DeviceWriter allows creating a device.
type DeviceWriter interface {
	Create(ctx context.Context, device entity.Device) error
	Update(ctx context.Context, device entity.Device) error
}

type DeviceReaderWriter interface {
	DeviceReader
	DeviceWriter
}

type ConfigurationReader interface {
	// Get returns the configuration for a device.
	Get(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error)
}

type CertificateWriter interface {
	SignCSR(ctx context.Context, csr []byte, cn string, ttl time.Duration) (entity.CertificateGroup, error)
	GetCertificate(ctx context.Context, serialNumber string) (entity.CertificateGroup, error)
}
