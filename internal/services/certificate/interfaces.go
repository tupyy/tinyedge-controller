package certificate

import (
	"context"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type CertificateReader interface {
	GetCertificate(ctx context.Context, serialNumber string) ([]byte, bool, time.Time, error)
	GetCACertificate(ctx context.Context) ([]byte, error)
}

type CertificateWriter interface {
	GenerateCertificate(ctx context.Context, cn string, ttl time.Duration) ([]byte, []byte, []byte, error)
	SignCSR(ctx context.Context, csr []byte, cn string, ttl time.Duration) ([]byte, error)
}

type CertificateReaderWriter interface {
	CertificateReader
	CertificateWriter
}

type ConfigurationCacheReaderWriter interface {
	Get(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error)
	Put(ctx context.Context, deviceID string, conf entity.ConfigurationResponse) error
	Delete(ctx context.Context, deviceID string) error
}
