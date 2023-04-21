package postgres

import (
	"context"
	"errors"
	"fmt"

	pgclient "github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/repo/models/mappers"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ManifestRepository struct {
	db             *gorm.DB
	client         pgclient.Client
	circuitBreaker pgclient.CircuitBreaker
}

func NewManifestRepository(client pgclient.Client) (*ManifestRepository, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &ManifestRepository{}, err
	}

	return &ManifestRepository{gormDB, client, client.GetCircuitBreaker()}, nil
}

func (m *ManifestRepository) InsertManifest(ctx context.Context, manifest entity.Manifest) error {
	if !m.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("manifest repository")
	}

	exists, err := m.isExists(ctx, manifest.GetID())
	if err != nil {
		return err
	}

	if exists {
		return errService.NewResourceAlreadyExistsError("manifest", manifest.GetID())
	}

	model := mappers.ManifestEntityToModel(manifest)

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return errService.NewPostgresNotAvailableError("manifest repository")
		}
		return err
	}
	return nil
}

func (m *ManifestRepository) UpdateManifest(ctx context.Context, manifest entity.Manifest) error {
	if !m.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("manifest repository")
	}

	model := mappers.ManifestEntityToModel(manifest)

	if err := m.getDb(ctx).Where("id = ?", model.ID).Save(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return errService.NewPostgresNotAvailableError("manifest repository")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errService.NewResourceNotFoundError("manifest", manifest.GetID())
		}
		return err
	}

	return nil
}

func (m *ManifestRepository) DeleteManifest(ctx context.Context, id string) error {
	if !m.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("manifest repository")
	}

	exists, err := m.isExists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	if err := m.getDb(ctx).Where("id = ?", id).Delete(&models.Manifest{}).Error; err != nil {
		if m.checkNetworkError(err) {
			return errService.NewPostgresNotAvailableError("manifest repository")
		}
		return err
	}

	return nil
}

func (m *ManifestRepository) CreateRelation(ctx context.Context, relation entity.Relation) error {
	switch relation.Type {
	case entity.NamespaceRelationType:
		return m.createNamespaceRelation(ctx, relation.ResourceID, relation.ManifestID)
	case entity.SetRelationType:
		return m.createSetRelation(ctx, relation.ResourceID, relation.ManifestID)
	case entity.DeviceRelationType:
		return m.createDeviceRelation(ctx, relation.ResourceID, relation.ManifestID)
	default:
		return errors.New("unknown relation type")
	}
}

func (m *ManifestRepository) DeleteRelation(ctx context.Context, relation entity.Relation) error {
	switch relation.Type {
	case entity.NamespaceRelationType:
		return m.deleteNamespaceRelation(ctx, relation.ResourceID, relation.ManifestID)
	case entity.SetRelationType:
		return m.deleteSetRelation(ctx, relation.ResourceID, relation.ManifestID)
	case entity.DeviceRelationType:
		return m.deleteDeviceRelation(ctx, relation.ResourceID, relation.ManifestID)
	default:
		return errors.New("unknown relation type")
	}

}

func (m *ManifestRepository) createNamespaceRelation(ctx context.Context, namespaceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("manifest repository")
	}

	model := models.NamespacesManifests{
		NamespaceID: namespaceID,
		ManifestID:  manifestID,
	}

	exists, err := m.isRelationExists(ctx, func(db *gorm.DB) *gorm.DB {
		var m models.NamespacesManifests
		return db.Where("namespace_id = ? AND manifest_id = ?", namespaceID, manifestID).First(&m)
	})
	if err != nil {
		return err
	}

	if exists {
		return errService.NewResourceAlreadyExistsError(fmt.Sprintf("relation between namespace %q and manifest", namespaceID), manifestID)
	}

	// check if the relation already exists
	var dummy models.NamespacesManifests
	if err := m.getDb(ctx).Where("namespace_id = ? AND manifest_id = ?", namespaceID, manifestID).First(&dummy).Error; err == nil {
		return nil
	}

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return errService.NewPostgresNotAvailableError("manifest repository")
		}
		return err
	}

	return nil
}

func (m *ManifestRepository) deleteNamespaceRelation(ctx context.Context, namespaceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("manifest repository")
	}

	exists, err := m.isRelationExists(ctx, func(db *gorm.DB) *gorm.DB {
		var m models.NamespacesManifests
		return db.Where("namespace_id = ? AND manifest_id = ?", namespaceID, manifestID).First(&m)
	})
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	model := models.NamespacesManifests{}
	if err := m.getDb(ctx).Where("namespace_id = ? AND manifest_id = ?", namespaceID, manifestID).Delete(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return errService.NewPostgresNotAvailableError("manifest repository")
		}
		return err
	}

	return nil
}

func (m *ManifestRepository) createSetRelation(ctx context.Context, setID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("manifest repository")
	}

	model := models.SetsManifests{
		DeviceSetID: setID,
		ManifestID:  manifestID,
	}

	exists, err := m.isRelationExists(ctx, func(db *gorm.DB) *gorm.DB {
		var m models.SetsManifests
		return db.Where("device_set_id = ? AND manifest_id = ?", setID, manifestID).First(&m)
	})
	if err != nil {
		return err
	}

	if exists {
		return errService.NewResourceAlreadyExistsError(fmt.Sprintf("relation between set %q and reference", setID), manifestID)
	}

	// check if the relation already exists
	var dummy models.SetsManifests
	if err := m.getDb(ctx).Where("device_set_id = ? AND manifest_id = ?", setID, manifestID).First(&dummy).Error; err == nil {
		return nil
	}

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return errService.NewPostgresNotAvailableError("manifest repository")
		}
		return err
	}

	return nil
}

func (m *ManifestRepository) deleteSetRelation(ctx context.Context, setID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("manifest repository")
	}

	exists, err := m.isRelationExists(ctx, func(db *gorm.DB) *gorm.DB {
		var m models.SetsManifests
		return db.Where("device_set_id = ? AND manifest_id = ?", setID, manifestID).First(&m)
	})
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	model := models.SetsManifests{}
	if err := m.getDb(ctx).Where("device_set_id = ? AND manifest_id = ?", setID, manifestID).Delete(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return errService.NewPostgresNotAvailableError("manifest repository")
		}
		return err
	}

	return nil
}

func (m *ManifestRepository) createDeviceRelation(ctx context.Context, deviceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("manifest repository")
	}

	exists, err := m.isRelationExists(ctx, func(db *gorm.DB) *gorm.DB {
		var m models.DevicesManifests
		return db.Where("device_id = ? AND manifest_id = ?", deviceID, manifestID).First(&m)
	})
	if err != nil {
		return err
	}

	if exists {
		return errService.NewResourceAlreadyExistsError(fmt.Sprintf("relation between device %q and manifest", deviceID), manifestID)
	}

	model := models.DevicesManifests{
		DeviceID:   deviceID,
		ManifestID: manifestID,
	}

	// check if the relation already exists
	var dummy models.DevicesManifests
	if err := m.getDb(ctx).Where("device_id = ? AND manifest_id = ?", deviceID, manifestID).First(&dummy).Error; err == nil {
		return nil
	}

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return errService.NewPostgresNotAvailableError("manifest repository")
		}
		return err
	}

	return nil
}

func (m *ManifestRepository) deleteDeviceRelation(ctx context.Context, deviceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("manifest repository")
	}

	exists, err := m.isRelationExists(ctx, func(db *gorm.DB) *gorm.DB {
		var m models.DevicesManifests
		return db.Where("device_id = ? AND manifest_id = ?", deviceID, manifestID).First(&m)
	})
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	model := models.DevicesManifests{}
	if err := m.getDb(ctx).Where("device_id = ? AND manifest_id = ?", deviceID, manifestID).Delete(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return errService.NewPostgresNotAvailableError("manifest repository")
		}
		return err
	}

	return nil
}

func (m *ManifestRepository) checkNetworkError(err error) (isOpen bool) {
	isOpen = m.circuitBreaker.BreakOnNetworkError(err)
	if isOpen {
		zap.S().Warn("circuit breaker is now open")
	}
	return
}

func (m *ManifestRepository) getDb(ctx context.Context) *gorm.DB {
	return m.db.Session(&gorm.Session{SkipHooks: true}).WithContext(ctx)
}

func (m *ManifestRepository) isExists(ctx context.Context, id string) (bool, error) {
	var model models.Manifest
	if err := m.getDb(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *ManifestRepository) isRelationExists(ctx context.Context, relationQuery func(db *gorm.DB) *gorm.DB) (bool, error) {
	if err := relationQuery(m.getDb(ctx)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func getManifest[T any](db *gorm.DB, searchKey, id string) ([]T, error) {
	models := []T{}

	if err := db.Where(searchKey, id).Find(&models).Error; err != nil {
		return models, err
	}

	if len(models) == 0 {
		return models, nil
	}

	return models, nil
}
