package postgres

import (
	"context"
	"errors"

	pgclient "github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres/mappers"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres/models"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
	"gorm.io/gorm"
)

type DeviceRepo struct {
	db             *gorm.DB
	client         pgclient.Client
	circuitBreaker pgclient.CircuitBreaker
}

func New(client pgclient.Client) (*DeviceRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &DeviceRepo{}, err
	}

	return &DeviceRepo{gormDB, client, client.GetCircuitBreaker()}, nil
}

func (d *DeviceRepo) Get(ctx context.Context, id string) (entity.Device, error) {
	m := models.Device{}

	if err := d.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Device{}, common.ErrDeviceNotFound
		}
		return entity.Device{}, err
	}
	return mappers.MapModelToEntity(m), nil
}

func (d *DeviceRepo) Create(ctx context.Context, device entity.Device) error {
	deviceModel := mappers.MapEntityToModel(device)

	if err := d.db.WithContext(ctx).Create(&deviceModel).Error; err != nil {
		return err
	}

	return nil
}

func (d *DeviceRepo) Update(ctx context.Context, device entity.Device) error {
	model := mappers.MapEntityToModel(device)

	if err := d.db.WithContext(ctx).Where("id = ?", model.ID).Save(&model).Error; err != nil {
		return err
	}

	return nil
}
