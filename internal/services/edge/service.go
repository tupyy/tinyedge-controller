package edge

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/certificate"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
	"go.uber.org/zap"
)

const (
	DefaultCertificateTTL = 3600 * 24 * 365 * time.Second
	BaseDomain            = "home.net"
	defaultNamespace      = "default"
)

type Service struct {
	deviceReaderWriter DeviceReaderWriter
	confReader         ConfigurationReader
	certService        *certificate.Service
}

func New(dr DeviceReaderWriter, confReader ConfigurationReader, certWriter *certificate.Service) *Service {
	return &Service{dr, confReader, certWriter}
}

// Enrol tries to enrol a device. If enable-auto-enrolment is true then the device is automatically
// enrolled. If false, the device is created but not enroled yet.
func (s *Service) Enrol(ctx context.Context, deviceID string) (status entity.EnrolStatus, err error) {
	d, err := s.deviceReaderWriter.GetDevice(ctx, deviceID)
	if err != nil {
		if !errService.IsResourceNotFound(err) {
			return entity.NotEnroledStatus, err
		}
		// device not found. create the device
		device := entity.Device{
			ID:          deviceID,
			NamespaceID: defaultNamespace,
			EnrolStatus: entity.EnroledStatus,
			EnroledAt:   time.Now().UTC(),
		}
		device.EnrolStatus = entity.EnroledStatus
		err = s.deviceReaderWriter.Create(ctx, device)
		if err != nil {
			return entity.NotEnroledStatus, err
		}
		zap.S().Infow("device enroled", "device_id", deviceID, "enrol_status", d.EnrolStatus)
		return device.EnrolStatus, nil
	}

	zap.S().Infow("device enroled", "device_id", deviceID, "enrol_status", d.EnrolStatus)
	return d.EnrolStatus, nil
}

func (s *Service) IsEnroled(ctx context.Context, deviceID string) (bool, error) {
	device, err := s.deviceReaderWriter.GetDevice(ctx, deviceID)
	if err != nil {
		if !errService.IsResourceNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return device.EnrolStatus == entity.EnroledStatus, nil
}

func (s *Service) Register(ctx context.Context, deviceID string, csr string) (entity.CertificateGroup, error) {
	device, err := s.deviceReaderWriter.GetDevice(ctx, deviceID)
	if err != nil {
		return entity.CertificateGroup{}, err
	}

	if device.EnrolStatus != entity.EnroledStatus {
		return entity.CertificateGroup{}, fmt.Errorf("unable to register the device. The device %s is not enroled yet", deviceID)
	}

	cn := fmt.Sprintf("%s.%s", deviceID, BaseDomain)
	certificate, err := s.certService.SignCSR(ctx, bytes.NewBufferString(csr).Bytes(), cn, DefaultCertificateTTL)
	if err != nil {
		return entity.CertificateGroup{}, fmt.Errorf("unable to sign the csr: %w", err)
	}

	device.Registred = true
	device.RegisteredAt = time.Now().UTC()
	device.CertificateSerialNumber = certificate.GetSerialNumber()

	if err := s.deviceReaderWriter.Update(ctx, device); err != nil {
		return entity.CertificateGroup{}, fmt.Errorf("unable to update device %q: %w", deviceID, err)
	}

	zap.S().Infow("device registered", "device_id", deviceID, "certificate_sn", device.CertificateSerialNumber)

	return certificate, nil
}

func (s *Service) IsRegistered(ctx context.Context, deviceID string) (bool, error) {
	device, err := s.deviceReaderWriter.GetDevice(ctx, deviceID)
	if err != nil {
		return false, err
	}
	return device.Registred, nil
}

func (s *Service) GetConfiguration(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error) {
	return s.confReader.GetDeviceConfiguration(ctx, deviceID)
}

// Heartbeat writes metrics from heartbeat.
func (s *Service) Heartbeat(ctx context.Context, heartbeat entity.Heartbeat) error {
	return nil
}
