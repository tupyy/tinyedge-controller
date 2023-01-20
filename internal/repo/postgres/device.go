package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pgclient "github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres/mappers"
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
	return mappers.DeviceToEntity(m), nil
}

func (d *DeviceRepo) GetDevices(ctx context.Context) ([]entity.Device, error) {
	if !d.circuitBreaker.IsAvailable() {
		return []entity.Device{}, common.ErrPostgresNotAvailable
	}

	m := []models.Device{}

	if err := d.getDb(ctx).Find(&m).Error; err != nil {
		if d.checkNetworkError(err) {
			return []entity.Device{}, common.ErrPostgresNotAvailable
		}
		return []entity.Device{}, err
	}

	if len(m) == 0 {
		return []entity.Device{}, nil
	}

	return mappers.DevicesToEntity(m), nil
}

func (d *DeviceRepo) GetConfiguration(ctx context.Context, id string) (entity.Configuration, error) {
	if !d.circuitBreaker.IsAvailable() {
		return entity.Configuration{}, common.ErrPostgresNotAvailable
	}
	m := models.Configuration{}

	if err := d.getDb(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if d.checkNetworkError(err) {
			return entity.Configuration{}, common.ErrPostgresNotAvailable
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Configuration{}, common.ErrResourceNotFound
		}
		return entity.Configuration{}, err
	}
	return mappers.ConfigurationToEntity(m), nil

}

func (d *DeviceRepo) GetSet(ctx context.Context, id string) (entity.Set, error) {
	if !d.circuitBreaker.IsAvailable() {
		return entity.Set{}, common.ErrPostgresNotAvailable
	}

	s := []models.SetJoin{}

	tx := d.getDb(ctx).Table("device_set").
		Select(`device_set.*, device.id as device_id,
		configuration.heartbeat_period_seconds as configuration_heartbeat_period_seconds, configuration.log_level as configuration_log_level,
		sets_references.manifest_reference_id as manifest_id`).
		Joins("LEFT JOIN device ON device.device_set_id = device_set.id").
		Joins("LEFT JOIN sets_references ON sets_references.device_set_id = device_set.id").
		Joins("LEFT JOIN configuration ON device_set.configuration_id = configuration.id").
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

	return mappers.SetToEntity(s), nil
}

func (d *DeviceRepo) GetSets(ctx context.Context) ([]entity.Set, error) {
	if !d.circuitBreaker.IsAvailable() {
		return []entity.Set{}, common.ErrPostgresNotAvailable
	}

	s := []models.SetJoin{}

	tx := d.getDb(ctx).Table("device_set").
		Select(`device_set.*, device.id as device_id,
		configuration.heartbeat_period_seconds as configuration_heartbeat_period_seconds, configuration.log_level as configuration_log_level,
		sets_references.manifest_reference_id as manifest_id`).
		Joins("LEFT JOIN device ON device.device_set_id = device_set.id").
		Joins("LEFT JOIN sets_references ON sets_references.device_set_id = device_set.id").
		Joins("LEFT JOIN configuration ON device_set.configuration_id = configuration.id")

	if err := tx.Find(&s).Error; err != nil {
		if d.checkNetworkError(err) {
			return []entity.Set{}, common.ErrPostgresNotAvailable
		}
		return []entity.Set{}, err
	}

	if len(s) == 0 {
		return []entity.Set{}, nil
	}

	return mappers.SetsToEntity(s), nil
}

func (d *DeviceRepo) GetNamespace(ctx context.Context, id string) (entity.Namespace, error) {
	if !d.circuitBreaker.IsAvailable() {
		return entity.Namespace{}, common.ErrPostgresNotAvailable
	}

	n := []models.NamespaceJoin{}
	tx := d.getDb(ctx).Table("namespace").
		Select(`device_set.*,namespace.*,
			configuration.heartbeat_period_seconds as configuration_heartbeat_period_seconds, configuration.log_level as configuration_log_level,
			device.id as device_id, device_set.id as device_set_id, namespaces_references.manifest_reference_id as manifest_id`).
		Joins("LEFT JOIN device ON device.namespace_id = namespace.id").
		Joins("LEFT JOIN device_set ON device_set.namespace_id = namespace.id").
		Joins("LEFT JOIN namespaces_references ON namespaces_references.namespace_id = namespace.id").
		Joins("LEFT JOIN configuration ON namespace.configuration_id = configuration.id").
		Where("namespace.id = ?", id)

	if err := tx.Find(&n).Error; err != nil {
		if d.checkNetworkError(err) {
			return entity.Namespace{}, common.ErrPostgresNotAvailable
		}
		return entity.Namespace{}, err
	}

	if len(n) == 0 {
		return entity.Namespace{}, common.ErrResourceNotFound
	}

	return mappers.NamespaceModelToEntity(n), nil
}

func (d *DeviceRepo) GetNamespaces(ctx context.Context) ([]entity.Namespace, error) {
	if !d.circuitBreaker.IsAvailable() {
		return []entity.Namespace{}, common.ErrPostgresNotAvailable
	}

	n := []models.NamespaceJoin{}
	tx := d.getDb(ctx).Table("namespace").
		Select(`device_set.*,namespace.*,
			configuration.heartbeat_period_seconds as configuration_heartbeat_period_seconds, configuration.log_level as configuration_log_level,
			device.id as device_id, device_set.id as device_set_id, namespaces_references.manifest_reference_id as manifest_id`).
		Joins("LEFT JOIN device ON device.namespace_id = namespace.id").
		Joins("LEFT JOIN device_set ON device_set.namespace_id = namespace.id").
		Joins("LEFT JOIN namespaces_references ON namespaces_references.namespace_id = namespace.id").
		Joins("LEFT JOIN configuration ON namespace.configuration_id = configuration.id")

	if err := tx.Find(&n).Error; err != nil {
		if d.checkNetworkError(err) {
			return []entity.Namespace{}, common.ErrPostgresNotAvailable
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Namespace{}, common.ErrResourceNotFound
		}
		return []entity.Namespace{}, err
	}

	if len(n) == 0 {
		return []entity.Namespace{}, nil
	}

	return mappers.NamespacesModelToEntity(n), nil
}

func (d *DeviceRepo) CreateSet(ctx context.Context, set entity.Set) error {
	_, err := d.GetSet(ctx, set.Name)
	if err == nil {
		return fmt.Errorf("%w set %q already exists", common.ErrResourceAlreadyExists, set.Name)
	} else if !errors.Is(err, common.ErrResourceNotFound) {
		return err
	}

	model := mappers.SetToModel(set)
	if err := d.getDb(ctx).Create(&model).Error; err != nil {
		return err
	}

	return nil
}

func (d *DeviceRepo) CreateNamespace(ctx context.Context, namespace entity.Namespace) error {
	if !d.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	// try to find if we have already a default namespace. If there is none, enforce the is_default on the current namespace.
	tx := d.getDb(ctx).Begin()

	var oldDefaultNamespace models.Namespace
	if err := tx.Where("is_default = ?", true).First(&oldDefaultNamespace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// default namespace was not found. Enforce this one
			zap.S().Debugf("no default namespace was found. enforce default flag on namespace %q", namespace.Name)
			namespace.IsDefault = true
		} else {
			tx.Commit()
			return fmt.Errorf("unable to unset is_default column %w", err)
		}
	} else if namespace.IsDefault {
		oldDefaultNamespace.IsDefault = sql.NullBool{Valid: true, Bool: false}
		if err := tx.Save(&oldDefaultNamespace).Error; err != nil {
			tx.Commit()
			return fmt.Errorf("unable to unset is_default to false for namespace %q: %w", oldDefaultNamespace.ID, err)
		}
	}

	model := mappers.NamespaceToModel(namespace)
	if err := d.getDb(ctx).Create(&model).Error; err != nil {
		return err
	}

	return tx.Commit().Error
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

	if err := d.getDb(ctx).Save(&model).Error; err != nil {
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
