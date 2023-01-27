package auth

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

//go:generate moq -out device_reader_moq.go . DeviceReader
type DeviceReader interface {
	GetDevice(ctx context.Context, id string) (entity.Device, error)
}

//go:generate moq -out cert_reader_moq.go . CertificateReader
type CertificateReader interface {
	GetCertificate(ctx context.Context, sn string) (entity.CertificateGroup, error)
}
