package common

import (
	"context"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

// DeviceReader is an interface that groups all the methods allowing to query/get devices.
type DeviceReader interface {
	GetDevice(ctx context.Context, id string) (entity.Device, error)
	GetNamespace(ctx context.Context, id string) (entity.Namespace, error)
	GetSet(ctx context.Context, id string) (entity.Set, error)
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

type ReferenceReader interface {
	GetReferences(ctx context.Context) ([]entity.ManifestReference, error)
	GetReference(ctx context.Context, id string) (entity.ManifestReference, error)
	GetRepositoryReferences(ctx context.Context, r entity.Repository) ([]entity.ManifestReference, error)
	GetRepositories(ctx context.Context) ([]entity.Repository, error)
	GetDeviceReferences(ctx context.Context, deviceID string) ([]entity.ManifestReference, error)
	GetSetReferences(ctx context.Context, setID string) ([]entity.ManifestReference, error)
	GetNamespaceReferences(ctx context.Context, setID string) ([]entity.ManifestReference, error)
}

type ReferenceWriter interface {
	InsertRepository(ctx context.Context, r entity.Repository) error
	UpdateRepository(ctx context.Context, r entity.Repository) error

	InsertReference(ctx context.Context, ref entity.ManifestReference) error
	UpdateReference(ctx context.Context, ref entity.ManifestReference) error
	DeleteReference(ctx context.Context, ref entity.ManifestReference) error

	CreateRelation(ctx context.Context, relation entity.ReferenceRelation) error
	DeleteRelation(ctx context.Context, relation entity.ReferenceRelation) error
}

type ReferenceReaderWriter interface {
	ReferenceReader
	ReferenceWriter
}

type GitReader interface {
	Open(ctx context.Context, r entity.Repository) (entity.Repository, error)
	Pull(ctx context.Context, r entity.Repository) error
	GetHeadSha(ctx context.Context, r entity.Repository) (string, error)
	GetManifests(ctx context.Context, repo entity.Repository) ([]entity.ManifestWork, error)
	GetManifest(ctx context.Context, ref entity.ManifestReference) (entity.ManifestWork, error)
}

type ConfigurationReader interface {
	// GetConfiguration returns the configuration for a device.
	GetConfiguration(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error)
}

type ConfigurationCacheReaderWriter interface {
	Get(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error)
	Put(ctx context.Context, deviceID string, conf entity.ConfigurationResponse) error
	Delete(ctx context.Context, deviceID string) error
}

type SecretReader interface {
	GetSecret(ctx context.Context, path, key string) (entity.Secret, error)
}
