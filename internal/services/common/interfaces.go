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
	GetManifests(ctx context.Context) ([]entity.ManifestReference, error)
	GetManifest(ctx context.Context, id string) (entity.ManifestReference, error)
	GetRepoManifests(ctx context.Context, r entity.Repository) ([]entity.ManifestReference, error)
	GetRepositories(ctx context.Context) ([]entity.Repository, error)
	GetDeviceManifests(ctx context.Context, deviceID string) ([]entity.ManifestReference, error)
	GetSetManifests(ctx context.Context, setID string) ([]entity.ManifestReference, error)
	GetNamespaceManifests(ctx context.Context, setID string) ([]entity.ManifestReference, error)
}

type ManifestWriter interface {
	InsertRepo(ctx context.Context, r entity.Repository) error
	UpdateRepo(ctx context.Context, r entity.Repository) error

	InsertManifest(ctx context.Context, manifest entity.ManifestReference) error
	UpdateManifest(ctx context.Context, manifest entity.ManifestReference) error
	DeleteManifest(ctx context.Context, manifest entity.ManifestReference) error

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
