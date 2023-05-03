package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pgclient "github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/repo/manifest"
	"github.com/tupyy/tinyedge-controller/internal/repo/models/mappers"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DeviceRepo struct {
	db             *gorm.DB
	client         pgclient.Client
	circuitBreaker pgclient.CircuitBreaker
	manifestReader manifest.ManifestReader
}

func NewDeviceRepo(client pgclient.Client) (*DeviceRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
		Logger:                 logger.Default.LogMode(logger.Info),
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &DeviceRepo{}, err
	}

	return &DeviceRepo{gormDB, client, client.GetCircuitBreaker(), manifest.ReadManifest}, nil
}

func (d *DeviceRepo) GetDevice(ctx context.Context, id string) (entity.Device, error) {
	if !d.circuitBreaker.IsAvailable() {
		return entity.Device{}, errService.NewPostgresNotAvailableError("device repository")
	}
	m := []models.DeviceJoin{}

	tx := deviceQuery(d.getDb(ctx)).Where("device.id = ?", id)
	if err := tx.Find(&m).Error; err != nil {
		if d.checkNetworkError(err) {
			return entity.Device{}, errService.NewPostgresNotAvailableError("device repository")
		}
		return entity.Device{}, err
	}

	if len(m) == 0 {
		return entity.Device{}, errService.NewResourceNotFoundError("device", id)
	}

	return mappers.DeviceToEntity(m, d.manifestReader)
}

func (d *DeviceRepo) GetDevices(ctx context.Context) ([]entity.Device, error) {
	if !d.circuitBreaker.IsAvailable() {
		return []entity.Device{}, errService.NewPostgresNotAvailableError("device repository")
	}

	m := []models.DeviceJoin{}

	if err := deviceQuery(d.getDb(ctx)).Find(&m).Error; err != nil {
		if d.checkNetworkError(err) {
			return []entity.Device{}, errService.NewPostgresNotAvailableError("device repository")
		}
		return []entity.Device{}, err
	}

	if len(m) == 0 {
		return []entity.Device{}, nil
	}

	return mappers.DevicesToEntity(m, d.manifestReader)
}

func (d *DeviceRepo) CreateDevice(ctx context.Context, device entity.Device) error {
	if !d.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("device repository")
	}

	deviceModel := mappers.DeviceEntityToModel(device)
	if err := d.getDb(ctx).Create(&deviceModel).Error; err != nil {
		return err
	}

	return nil
}

func (d *DeviceRepo) DeleteDevice(ctx context.Context, id string) error {
	if !d.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("device repository")
	}
	return d.getDb(ctx).Where("id = ?", id).Delete(&models.Device{}).Error
}

func (d *DeviceRepo) UpdateDevice(ctx context.Context, device entity.Device) error {
	if !d.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("device repository")
	}

	_, err := d.GetDevice(ctx, device.ID)
	if err != nil {
		return err
	}

	model := mappers.DeviceEntityToModel(device)
	if err := d.getDb(ctx).Save(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errService.NewResourceNotFoundError("device", device.ID)
		}
		return err
	}

	return nil
}

func (d *DeviceRepo) GetSet(ctx context.Context, id string) (entity.Set, error) {
	if !d.circuitBreaker.IsAvailable() {
		return entity.Set{}, errService.NewPostgresNotAvailableError("device repository")
	}

	s := []models.SetJoin{}

	tx := setQuery(d.getDb(ctx)).Where("device_set.id = ?", id)
	if err := tx.Find(&s).Error; err != nil {
		if d.checkNetworkError(err) {
			return entity.Set{}, errService.NewPostgresNotAvailableError("device repository")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Set{}, errService.NewResourceNotFoundError("set", id)
		}
		return entity.Set{}, err
	}

	if len(s) == 0 {
		return entity.Set{}, errService.NewResourceNotFoundError("set", id)
	}

	return mappers.SetToEntity(s, d.manifestReader)
}

func (d *DeviceRepo) GetSets(ctx context.Context) ([]entity.Set, error) {
	if !d.circuitBreaker.IsAvailable() {
		return []entity.Set{}, errService.NewPostgresNotAvailableError("device repository")
	}

	s := []models.SetJoin{}

	if err := setQuery(d.getDb(ctx)).Find(&s).Error; err != nil {
		if d.checkNetworkError(err) {
			return []entity.Set{}, errService.NewPostgresNotAvailableError("device repository")
		}
		return []entity.Set{}, err
	}

	if len(s) == 0 {
		return []entity.Set{}, nil
	}

	return mappers.SetsToEntity(s, d.manifestReader)
}

func (d *DeviceRepo) CreateSet(ctx context.Context, set entity.Set) error {
	_, err := d.GetSet(ctx, set.Name)
	if err == nil {
		return errService.NewResourceAlreadyExistsError("set", set.Name)
	} else if _, ok := err.(errService.ResourseNotFoundError); !ok {
		return err
	}

	model := mappers.SetToModel(set)
	if err := d.getDb(ctx).Create(&model).Error; err != nil {
		return err
	}

	return nil
}

func (d *DeviceRepo) DeleteSet(ctx context.Context, id string) error {
	return d.getDb(ctx).Where("id = ?", id).Delete(&models.DeviceSet{}).Error
}

func (d *DeviceRepo) UpdateSet(ctx context.Context, set entity.Set) error {
	if !d.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("device repository")
	}

	oldSet := models.DeviceSet{}
	if err := d.getDb(ctx).Where("id = ?", set.Name).First(&oldSet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errService.NewResourceNotFoundError("set", set.Name)
		}
		return err
	}

	oldSet.NamespaceID = set.NamespaceID
	if set.Configuration != nil {
		oldSet.ConfigurationID = sql.NullString{Valid: true, String: set.Configuration.GetID()}
	} else {
		oldSet.ConfigurationID = sql.NullString{Valid: false}
	}

	return d.getDb(ctx).Save(oldSet).Error
}

func (d *DeviceRepo) GetNamespace(ctx context.Context, id string) (entity.Namespace, error) {
	if !d.circuitBreaker.IsAvailable() {
		return entity.Namespace{}, errService.NewPostgresNotAvailableError("device repository")
	}

	n := []models.NamespaceJoin{}

	tx := namespaceQuery(d.getDb(ctx)).Where("namespace.id = ?", id)
	if err := tx.Find(&n).Error; err != nil {
		if d.checkNetworkError(err) {
			return entity.Namespace{}, errService.NewPostgresNotAvailableError("device repository")
		}
		return entity.Namespace{}, err
	}

	if len(n) == 0 {
		return entity.Namespace{}, errService.NewResourceNotFoundError("namespace", id)
	}

	return mappers.NamespaceModelToEntity(n, d.manifestReader)
}

func (d *DeviceRepo) GetDefaultNamespace(ctx context.Context) (entity.Namespace, error) {
	if !d.circuitBreaker.IsAvailable() {
		return entity.Namespace{}, errService.NewPostgresNotAvailableError("device repository")
	}

	n := []models.NamespaceJoin{}
	tx := namespaceQuery(d.getDb(ctx)).Where("namespace.is_default = ?", true)

	if err := tx.Find(&n).Error; err != nil {
		if d.checkNetworkError(err) {
			return entity.Namespace{}, errService.NewPostgresNotAvailableError("device repository")
		}
		return entity.Namespace{}, err
	}

	if len(n) == 0 {
		return entity.Namespace{}, errService.NewResourceNotFoundErrorWithReason("Default namespace not found")
	}

	return mappers.NamespaceModelToEntity(n, d.manifestReader)
}

func (d *DeviceRepo) GetNamespaces(ctx context.Context) ([]entity.Namespace, error) {
	if !d.circuitBreaker.IsAvailable() {
		return []entity.Namespace{}, errService.NewPostgresNotAvailableError("device repository")
	}

	n := []models.NamespaceJoin{}
	tx := namespaceQuery(d.getDb(ctx))

	if err := tx.Find(&n).Error; err != nil {
		if d.checkNetworkError(err) {
			return []entity.Namespace{}, errService.NewPostgresNotAvailableError("device repository")
		}
		return []entity.Namespace{}, err
	}

	if len(n) == 0 {
		return []entity.Namespace{}, nil
	}

	return mappers.NamespacesModelToEntity(n, d.manifestReader)
}

func (d *DeviceRepo) CreateNamespace(ctx context.Context, namespace entity.Namespace) error {
	if !d.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("device repository")
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

func (d *DeviceRepo) UpdateNamespace(ctx context.Context, namespace entity.Namespace) error {
	if !d.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("device repository")
	}

	model, err := d.GetNamespace(ctx, namespace.Name)
	if err != nil {
		return err
	}

	tx := d.getDb(ctx).Begin()

	if model.IsDefault && !namespace.IsDefault {
		// count the namespaces. We always have one default namespace
		count := []models.Namespace{}
		err := tx.Find(&count).Error
		if err != nil {
			tx.Commit()
			return fmt.Errorf("unble to count namespaces")
		}
		if len(count) == 1 {
			tx.Commit()
			return fmt.Errorf("cannot set is_default to false for the only namespace")
		}
	}

	if namespace.IsDefault && !model.IsDefault {
		var defaultNamespace models.Namespace
		if err := tx.Where("is_default = ?", true).First(&defaultNamespace).Error; err != nil {
			tx.Commit()
			return fmt.Errorf("unable to find the default namespace: %w", err)
		}
		defaultNamespace.IsDefault = sql.NullBool{Valid: true, Bool: false}
		if err := tx.Save(&defaultNamespace).Error; err != nil {
			tx.Commit()
			return fmt.Errorf("unable to unset is_default to false for namespace %q: %w", defaultNamespace.ID, err)
		}
	}

	m := mappers.NamespaceToModel(namespace)
	if err := tx.Save(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errService.NewResourceNotFoundError("namespace", namespace.Name)
		}
		return err
	}

	return tx.Commit().Error
}

func (d *DeviceRepo) DeleteNamespace(ctx context.Context, id string) error {
	if !d.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("device repository")
	}

	tx := d.getDb(ctx).Begin()

	var n models.Namespace
	if err := tx.Where("id = ?", id).First(&n).Error; err != nil {
		tx.Commit()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errService.NewResourceNotFoundError("namespace", id)
		}
		return errService.NewDeleteResourceError("namespace", id, err.Error())
	}

	if n.IsDefault.Bool {
		var nextNamespace models.Namespace
		if err := tx.Where("is_default = ?", false).First(&nextNamespace).Error; err != nil {
			tx.Commit()
			return errService.NewDeleteResourceError("namespace", id, err.Error())
		}
		nextNamespace.IsDefault = sql.NullBool{Valid: true, Bool: true}
		if err := tx.Save(&nextNamespace).Error; err != nil {
			tx.Commit()
			return errService.NewDeleteResourceError("namespace", id, err.Error())
		}
	}

	if err := tx.Where("id = ?", id).Delete(&models.Namespace{}).Error; err != nil {
		tx.Commit()
		return errService.NewDeleteResourceError("namespace", id, err.Error())
	}

	return tx.Commit().Error
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
