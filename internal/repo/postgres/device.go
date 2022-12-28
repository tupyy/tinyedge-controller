package postgres

import (
	"context"
	"errors"

	pgclient "github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres/mappers"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres/models"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DeviceRepo struct {
	db             *gorm.DB
	client         pgclient.Client
	circuitBreaker pgclient.CircuitBreaker
}

func NewDeviceRepo(client pgclient.Client) (*DeviceRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &DeviceRepo{}, err
	}

	return &DeviceRepo{gormDB, client, client.GetCircuitBreaker()}, nil
}

func (d *DeviceRepo) GetDevice(ctx context.Context, id string) (entity.Device, error) {
	if !d.circuitBreaker.IsAvailable() {
		return entity.Device{}, common.ErrPostgresNotAvailable
	}
	m := models.Device{}

	if err := d.getDb(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if d.checkNetworkError(err) {
			return entity.Device{}, common.ErrPostgresNotAvailable
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Device{}, common.ErrResourceNotFound
		}
		return entity.Device{}, err
	}
	return mappers.MapModelToEntity(m), nil
}

func (d *DeviceRepo) GetSet(ctx context.Context, id string) (entity.Set, error) {
	if !d.circuitBreaker.IsAvailable() {
		return entity.Set{}, common.ErrPostgresNotAvailable
	}

	s := []models.SetJoin{}

	tx := d.getDb(ctx).Table("device_set").
		Select(`device_set.*, device.id as device_id, sets_workloads.manifest_work_id as manifest_id`).
		Joins("LEFT JOIN device ON device.device_set_id = device_set.id").
		Joins("LEFT JOIN sets_workloads ON sets_workloads.device_set_id = device_set.id").
		Where("device_set.id = ?", id)

	if err := tx.Find(&s).Error; err != nil {
		if d.checkNetworkError(err) {
			return entity.Set{}, common.ErrPostgresNotAvailable
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Set{}, common.ErrResourceNotFound
		}
		return entity.Set{}, err
	}

	if len(s) == 0 {
		return entity.Set{}, common.ErrResourceNotFound
	}

	model := models.DeviceSet{
		ID:              s[0].ID,
		NamespaceID:     s[0].NamespaceID,
		ConfigurationID: s[0].ConfigurationID,
	}
	ids := make([]string, 0, len(s))
	manifests := make([]string, 0, len(s))
	for _, ss := range s {
		if ss.DeviceId != "" {
			ids = append(ids, ss.DeviceId)
		}
		if ss.ManifestId != "" {
			manifests = append(manifests, ss.ManifestId)
		}
	}

	return mappers.SetModelToEntity(model, ids, manifests), nil
}

func (d *DeviceRepo) GetNamespace(ctx context.Context, id string) (entity.Namespace, error) {
	if !d.circuitBreaker.IsAvailable() {
		return entity.Namespace{}, common.ErrPostgresNotAvailable
	}

	n := []models.NamespaceJoin{}
	tx := d.getDb(ctx).Table("namespace").
		Select(`device_set.*, device.id as device_id, device_set.id as device_set_id, namespaces_workloads.manifest_work_id as manifest_id`).
		Joins("LEFT JOIN device ON device.namespace_id = namespace.id").
		Joins("LEFT JOIN device_set ON device_set.namespace_id = namespace.id").
		Joins("LEFT JOIN namespaces_workloads ON namespaces_workloads.namespace_id = namespace.id").
		Where("namespace.id = ?", id)

	if err := tx.Find(&n).Error; err != nil {
		if d.checkNetworkError(err) {
			return entity.Namespace{}, common.ErrPostgresNotAvailable
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Namespace{}, common.ErrResourceNotFound
		}
		return entity.Namespace{}, err
	}

	if len(n) == 0 {
		return entity.Namespace{}, common.ErrResourceNotFound
	}

	namespace := models.Namespace{
		ID:              n[0].ID,
		IsDefault:       n[0].IsDefault,
		ConfigurationID: n[0].ConfigurationID,
	}
	sets := make([]string, 0, len(n))
	devices := make([]string, 0, len(n))
	manifests := make([]string, 0, len(n))
	for _, nn := range n {
		if nn.SetId != "" {
			sets = append(sets, nn.SetId)
		}
		if nn.DeviceId != "" {
			devices = append(devices, nn.DeviceId)
		}
		if nn.ManifestId != "" {
			manifests = append(manifests, nn.ManifestId)
		}
	}

	return mappers.NamespaceModelToEntity(namespace, sets, devices, manifests), nil
}

func (d *DeviceRepo) Create(ctx context.Context, device entity.Device) error {
	deviceModel := mappers.MapEntityToModel(device)

	if err := d.getDb(ctx).Create(&deviceModel).Error; err != nil {
		return err
	}

	return nil
}

func (d *DeviceRepo) Update(ctx context.Context, device entity.Device) error {
	model := mappers.MapEntityToModel(device)

	if err := d.getDb(ctx).Where("id = ?", model.ID).Save(&model).Error; err != nil {
		return err
	}

	return nil
}

func (d *DeviceRepo) checkNetworkError(err error) (isOpen bool) {
	isOpen = d.circuitBreaker.BreakOnNetworkError(err)
	if isOpen {
		zap.S().Warn("circuit breaker is now open")
	}
	return
}

func (d *DeviceRepo) getDb(ctx context.Context) *gorm.DB {
	return d.db.Session(&gorm.Session{SkipHooks: true}).WithContext(ctx)
}
