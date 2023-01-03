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
}

type CertificateWriter interface {
	GenerateCertificate(ctx context.Context, cn string, ttl time.Duration) ([]byte, []byte, error)
	SignCSR(ctx context.Context, csr []byte, cn string, ttl time.Duration) ([]byte, error)
}

type CertificateReaderWriter interface {
	CertificateReader
	CertificateWriter
}

type ManifestReader interface {
	GetManifests(ctx context.Context) ([]entity.ManifestWorkV1, error)
	GetManifest(ctx context.Context, id string) (entity.ManifestWorkV1, error)
	GetRepoManifests(ctx context.Context, r entity.Repository) ([]entity.ManifestWorkV1, error)
	GetRepositories(ctx context.Context) ([]entity.Repository, error)
	GetDeviceManifests(ctx context.Context, deviceID string) ([]entity.ManifestWorkV1, error)
	GetSetManifests(ctx context.Context, setID string) ([]entity.ManifestWorkV1, error)
	GetNamespaceManifests(ctx context.Context, setID string) ([]entity.ManifestWorkV1, error)
}

type ManifestWriter interface {
	InsertRepo(ctx context.Context, r entity.Repository) error
	UpdateRepo(ctx context.Context, r entity.Repository) error

	InsertManifest(ctx context.Context, manifest entity.ManifestWorkV1) error
	UpdateManifest(ctx context.Context, manifest entity.ManifestWorkV1) error
	DeleteManifest(ctx context.Context, manifest entity.ManifestWorkV1) error

	CreateNamespaceRelation(ctx context.Context, namespace, manifestID string) error
	CreateSetRelation(ctx context.Context, set, manifestID string) error
	CreateDeviceRelation(ctx context.Context, deviceID, manifestID string) error
	DeleteNamespaceRelation(ctx context.Context, namespace, manifestID string) error
	DeleteSetRelation(ctx context.Context, set, manifestID string) error
	DeleteDeviceRelation(ctx context.Context, deviceID, manifestID string) error
}

type ManifestReaderWriter interface {
	ManifestReader
	ManifestWriter
}

type GitReaderWriter interface {
	Open(ctx context.Context, r entity.Repository) (entity.Repository, error)
	Pull(ctx context.Context, r entity.Repository) error
	GetHeadSha(ctx context.Context, r entity.Repository) (string, error)
	GetManifests(ctx context.Context, repo entity.Repository) ([]entity.ManifestWorkV1, error)
}

type ConfigurationReader interface {
	// GetConfiguration returns the configuration for a device.
	GetConfiguration(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error)
}

type ConfigurationCacheReaderWriter interface {
	Get(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error)
	Put(ctx context.Context, deviceID string, conf entity.ConfigurationResponse) error
}
