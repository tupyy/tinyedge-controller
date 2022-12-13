package edge

import (
	"context"
	"errors"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
	"go.uber.org/zap"
)

// DeviceReader is an interface that groups all the methods allowing to query/get devices.
type DeviceReader interface {
	Get(ctx context.Context, id string) (entity.Device, error)
}

// DeviceWriter allows creating a device.
type DeviceWriter interface {
	Create(ctx context.Context, device entity.Device) error
}

type DeviceReaderWriter interface {
	DeviceReader
	DeviceWriter
}

type ConfigurationReader interface {
	// Get returns the configuration for a device.
	Get(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error)
}

type Service struct {
	deviceRepo DeviceReaderWriter
	confReader ConfigurationReader
}

func New(dr DeviceReaderWriter, confReader ConfigurationReader) *Service {
	return &Service{dr, confReader}
}

// Enrol tries to enrol a device. If enable-auto-enrolment is true then the device is automatically
// enrolled. If false, the device is created but not enroled yet.
func (s *Service) Enrol(ctx context.Context, device entity.Device) (status entity.EnrolStatus, err error) {
	d, err := s.deviceRepo.Get(ctx, device.ID)
	if err != nil {
		if !errors.Is(err, common.ErrDeviceNotFound) {
			return entity.NotEnroledStatus, err
		}
		// device not found. create the device
		device.EnrolStatus = entity.EnroledStatus
		err = s.deviceRepo.Create(ctx, device)
		if err != nil {
			return entity.NotEnroledStatus, err
		}
		zap.S().Infow("enrol device", "device_id", device.ID, "enrol status", d.EnrolStatus)
		return device.EnrolStatus, nil
	}

	zap.S().Infow("enrol device", "device_id", device.ID, "enrol status", d.EnrolStatus)
	return d.EnrolStatus, nil
}

func (s *Service) IsEnroled(ctx context.Context, deviceID string) (bool, error) {
	device, err := s.deviceRepo.Get(ctx, deviceID)
	if err != nil {
		if errors.Is(err, common.ErrDeviceNotFound) {
			return false, nil
		}
		return false, err
	}
	return device.EnrolStatus == entity.EnroledStatus, nil
}

func (s *Service) Register(ctx context.Context, device entity.Device, csr string) (deviceRegistered bool, err error) {
	return false, nil
}

func (s *Service) IsRegistered(ctx context.Context, deviceID string) (bool, error) {
	device, err := s.deviceRepo.Get(ctx, deviceID)
	if err != nil {
		return false, err
	}
	return device.Registred, nil
}

func (s *Service) GetConfiguration(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error) {
	return s.confReader.Get(ctx, deviceID)
}

// Heartbeat writes metrics from heartbeat.
func (s *Service) Heartbeat(ctx context.Context, heartbeat entity.Heartbeat) error {
	return nil
}
