package postgres

import (
	"context"
	"errors"

	pgclient "github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres/mappers"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReferenceRepository struct {
	db             *gorm.DB
	client         pgclient.Client
	circuitBreaker pgclient.CircuitBreaker
}

func NewReferenceRepository(client pgclient.Client) (*ReferenceRepository, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &ReferenceRepository{}, err
	}

	return &ReferenceRepository{gormDB, client, client.GetCircuitBreaker()}, nil
}

func (m *ReferenceRepository) GetReferences(ctx context.Context) ([]entity.ManifestReference, error) {
	if !m.circuitBreaker.IsAvailable() {
		return []entity.ManifestReference{}, common.ErrPostgresNotAvailable
	}

	manifests := []models.ManifestJoin{}

	tx := newManifestQuery(ctx, m.db).Build()
	if err := tx.Find(&manifests).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.ManifestReference{}, common.ErrPostgresNotAvailable
		}
		return []entity.ManifestReference{}, err
	}

	if len(manifests) == 0 {
		return []entity.ManifestReference{}, nil
	}

	return mappers.ManifestModelsToEntities(manifests), nil

}

func (m *ReferenceRepository) GetReference(ctx context.Context, id string) (entity.ManifestReference, error) {
	if !m.circuitBreaker.IsAvailable() {
		return entity.ManifestReference{}, common.ErrPostgresNotAvailable
	}

	manifests := []models.ManifestJoin{}

	tx := newManifestQuery(ctx, m.db).WithReferenceID(id).Build()
	if err := tx.Find(&manifests).Error; err != nil {
		if m.checkNetworkError(err) {
			return entity.ManifestReference{}, common.ErrPostgresNotAvailable
		}
		return entity.ManifestReference{}, err
	}

	if len(manifests) == 0 {
		return entity.ManifestReference{}, common.ErrResourceNotFound
	}

	return mappers.ManifestModelToEntity(manifests), nil
}

func (m *ReferenceRepository) GetRepositories(ctx context.Context) ([]entity.Repository, error) {
	if !m.circuitBreaker.IsAvailable() {
		return []entity.Repository{}, common.ErrPostgresNotAvailable
	}

	repos := []models.Repo{}

	if err := m.getDb(ctx).Find(&repos).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.Repository{}, common.ErrPostgresNotAvailable
		}
		return []entity.Repository{}, err
	}

	if len(repos) == 0 {
		return []entity.Repository{}, nil
	}

	entities := make([]entity.Repository, 0, len(repos))
	for _, r := range repos {
		entities = append(entities, mappers.RepoModelToEntity(r))
	}

	return entities, nil
}

func (m *ReferenceRepository) GetRepositoryReferences(ctx context.Context, r entity.Repository) ([]entity.ManifestReference, error) {
	if !m.circuitBreaker.IsAvailable() {
		return []entity.ManifestReference{}, common.ErrPostgresNotAvailable
	}

	manifests := []models.ManifestJoin{}

	tx := newManifestQuery(ctx, m.db).WithRepoId(r.Id).Build()
	if err := tx.Find(&manifests).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.ManifestReference{}, common.ErrPostgresNotAvailable
		}
		return []entity.ManifestReference{}, err
	}

	if len(manifests) == 0 {
		return []entity.ManifestReference{}, nil
	}

	return mappers.ManifestModelsToEntities(manifests), nil
}

func (m *ReferenceRepository) InsertRepository(ctx context.Context, r entity.Repository) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := mappers.RepoEntityToModel(r)

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}
	return nil
}

func (m *ReferenceRepository) UpdateRepository(ctx context.Context, r entity.Repository) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := mappers.RepoEntityToModel(r)

	if err := m.getDb(ctx).Where("id = ?", model.ID).Save(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ReferenceRepository) InsertReference(ctx context.Context, ref entity.ManifestReference) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	exists, err := m.isExists(ctx, ref)
	if err != nil {
		return err
	}

	if exists {
		return common.ErrResourceAlreadyExists
	}

	model := mappers.ManifestEntityToModel(ref)

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}
	return nil
}

func (m *ReferenceRepository) UpdateReference(ctx context.Context, ref entity.ManifestReference) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := mappers.ManifestEntityToModel(ref)

	if err := m.getDb(ctx).Where("id = ?", model.ID).Save(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrResourceNotFound
		}
		return err
	}

	return nil
}

func (m *ReferenceRepository) DeleteReference(ctx context.Context, ref entity.ManifestReference) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	exists, err := m.isExists(ctx, ref)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	if err := m.getDb(ctx).Where("id = ?", ref.Id).Delete(&models.ManifestWork{}).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ReferenceRepository) GetNamespaceReferences(ctx context.Context, namespaceID string) ([]entity.ManifestReference, error) {
	if !m.circuitBreaker.IsAvailable() {
		return []entity.ManifestReference{}, common.ErrPostgresNotAvailable
	}

	models := []models.ManifestJoin{}
	tx := newManifestQuery(ctx, m.db).WithNamespaceID(namespaceID).Build()
	if err := tx.Find(&models).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.ManifestReference{}, common.ErrPostgresNotAvailable
		}
		return []entity.ManifestReference{}, err
	}

	if len(models) == 0 {
		return []entity.ManifestReference{}, nil
	}

	return mappers.ManifestModelsToEntities(models), nil
}

func (m *ReferenceRepository) GetSetReferences(ctx context.Context, setID string) ([]entity.ManifestReference, error) {
	if !m.circuitBreaker.IsAvailable() {
		return []entity.ManifestReference{}, common.ErrPostgresNotAvailable
	}

	models := []models.ManifestJoin{}

	tx := newManifestQuery(ctx, m.db).WithSetID(setID).Build()
	if err := tx.Find(&models).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.ManifestReference{}, common.ErrPostgresNotAvailable
		}
		return []entity.ManifestReference{}, err
	}

	if len(models) == 0 {
		return []entity.ManifestReference{}, nil
	}

	return mappers.ManifestModelsToEntities(models), nil
}

func (m *ReferenceRepository) GetDeviceReferences(ctx context.Context, deviceID string) ([]entity.ManifestReference, error) {
	if !m.circuitBreaker.IsAvailable() {
		return []entity.ManifestReference{}, common.ErrPostgresNotAvailable
	}

	models := []models.ManifestJoin{}

	tx := newManifestQuery(ctx, m.db).WithDeviceID(deviceID).Build()
	if err := tx.Find(&models).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.ManifestReference{}, common.ErrPostgresNotAvailable
		}
		return []entity.ManifestReference{}, err
	}

	if len(models) == 0 {
		return []entity.ManifestReference{}, nil
	}

	return mappers.ManifestModelsToEntities(models), nil
}

func (m *ReferenceRepository) CreateRelation(ctx context.Context, relation entity.ReferenceRelation) error {
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

func (m *ReferenceRepository) DeleteRelation(ctx context.Context, relation entity.ReferenceRelation) error {
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

func (m *ReferenceRepository) createNamespaceRelation(ctx context.Context, namespaceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := models.NamespacesWorkloads{
		NamespaceID:         namespaceID,
		ManifestReferenceID: manifestID,
	}

	exists, err := m.isRelationExists(ctx, func(db *gorm.DB) *gorm.DB {
		var m models.NamespacesWorkloads
		return db.Where("namespace_id = ? AND manifest_reference_id = ?", namespaceID, manifestID).First(&m)
	})
	if err != nil {
		return err
	}

	if exists {
		return common.ErrResourceAlreadyExists
	}

	// check if the relation already exists
	var dummy models.NamespacesWorkloads
	if err := m.getDb(ctx).Where("namespace_id = ? AND manifest_reference_id = ?", namespaceID, manifestID).First(&dummy).Error; err == nil {
		return nil
	}

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ReferenceRepository) deleteNamespaceRelation(ctx context.Context, namespaceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	exists, err := m.isRelationExists(ctx, func(db *gorm.DB) *gorm.DB {
		var m models.NamespacesWorkloads
		return db.Where("namespace_id = ? AND manifest_reference_id = ?", namespaceID, manifestID).First(&m)
	})
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	model := models.NamespacesWorkloads{}
	if err := m.getDb(ctx).Where("namespace_id = ? AND manifest_reference_id = ?", namespaceID, manifestID).Delete(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ReferenceRepository) createSetRelation(ctx context.Context, setID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := models.SetsWorkloads{
		DeviceSetID:         setID,
		ManifestReferenceID: manifestID,
	}

	exists, err := m.isRelationExists(ctx, func(db *gorm.DB) *gorm.DB {
		var m models.SetsWorkloads
		return db.Where("device_set_id = ? AND manifest_reference_id = ?", setID, manifestID).First(&m)
	})
	if err != nil {
		return err
	}

	if exists {
		return common.ErrResourceAlreadyExists
	}

	// check if the relation already exists
	var dummy models.SetsWorkloads
	if err := m.getDb(ctx).Where("device_set_id = ? AND manifest_reference_id = ?", setID, manifestID).First(&dummy).Error; err == nil {
		return nil
	}

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ReferenceRepository) deleteSetRelation(ctx context.Context, setID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	exists, err := m.isRelationExists(ctx, func(db *gorm.DB) *gorm.DB {
		var m models.SetsWorkloads
		return db.Where("device_set_id = ? AND manifest_reference_id = ?", setID, manifestID).First(&m)
	})
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	model := models.SetsWorkloads{}
	if err := m.getDb(ctx).Where("device_set_id = ? AND manifest_reference_id = ?", setID, manifestID).Delete(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ReferenceRepository) createDeviceRelation(ctx context.Context, deviceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	exists, err := m.isRelationExists(ctx, func(db *gorm.DB) *gorm.DB {
		var m models.DevicesWorkloads
		return db.Where("device_id = ? AND manifest_reference_id = ?", deviceID, manifestID).First(&m)
	})
	if err != nil {
		return err
	}

	if exists {
		return common.ErrResourceAlreadyExists
	}

	model := models.DevicesWorkloads{
		DeviceID:            deviceID,
		ManifestReferenceID: manifestID,
	}

	// check if the relation already exists
	var dummy models.DevicesWorkloads
	if err := m.getDb(ctx).Where("device_id = ? AND manifest_reference_id = ?", deviceID, manifestID).First(&dummy).Error; err == nil {
		return nil
	}

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ReferenceRepository) deleteDeviceRelation(ctx context.Context, deviceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	exists, err := m.isRelationExists(ctx, func(db *gorm.DB) *gorm.DB {
		var m models.DevicesWorkloads
		return db.Where("device_id = ? AND manifest_reference_id = ?", deviceID, manifestID).First(&m)
	})
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	model := models.DevicesWorkloads{}
	if err := m.getDb(ctx).Where("device_id = ? AND manifest_reference_id = ?", deviceID, manifestID).Delete(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ReferenceRepository) checkNetworkError(err error) (isOpen bool) {
	isOpen = m.circuitBreaker.BreakOnNetworkError(err)
	if isOpen {
		zap.S().Warn("circuit breaker is now open")
	}
	return
}

func (m *ReferenceRepository) getDb(ctx context.Context) *gorm.DB {
	return m.db.Session(&gorm.Session{SkipHooks: true}).WithContext(ctx)
}

func (m *ReferenceRepository) isExists(ctx context.Context, manifest entity.ManifestReference) (bool, error) {
	var model models.ManifestReference
	if err := m.getDb(ctx).Where("id = ?", manifest.Id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *ReferenceRepository) isRelationExists(ctx context.Context, relationQuery func(db *gorm.DB) *gorm.DB) (bool, error) {
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

type manifestQueryBuilder struct {
	tx *gorm.DB
}

func newManifestQuery(ctx context.Context, db *gorm.DB) *manifestQueryBuilder {
	tx := db.Session(&gorm.Session{SkipHooks: true}).WithContext(ctx).Table("manifest_reference").
		Select(`manifest_reference.*, devices_workloads.device_id as device_id, sets_workloads.device_set_id as set_id, namespaces_workloads.namespace_id as namespace_id,
		repo.id as repo_id, repo.url as repo_url, repo.branch as repo_branch, repo.local_path as repo_local_path,
		repo.current_head_sha as repo_current_head_sha, repo.target_head_sha as repo_target_head_sha,
		repo.pull_period_seconds as repo_pull_period_seconds`).
		Joins("LEFT JOIN namespaces_workloads ON namespaces_workloads.manifest_reference_id = manifest_reference.id").
		Joins("LEFT JOIN sets_workloads ON sets_workloads.manifest_reference_id = manifest_reference.id").
		Joins("LEFT JOIN devices_workloads ON devices_workloads.manifest_reference_id = manifest_reference.id").
		Joins("JOIN repo ON repo.id = manifest_reference.repo_id")
	return &manifestQueryBuilder{tx}
}

func (mm *manifestQueryBuilder) WithRepoId(id string) *manifestQueryBuilder {
	mm.tx.Where("repo_id = ?", id)
	return mm
}

func (mm *manifestQueryBuilder) WithReferenceID(id string) *manifestQueryBuilder {
	mm.tx.Where("manifest_reference.id = ?", id)
	return mm
}

func (mm *manifestQueryBuilder) WithNamespaceID(id string) *manifestQueryBuilder {
	mm.tx.Where("namespaces_workloads.namespace_id = ?", id)
	return mm
}

func (mm *manifestQueryBuilder) WithDeviceID(id string) *manifestQueryBuilder {
	mm.tx.Where("devices_workloads.device_id = ?", id)
	return mm
}

func (mm *manifestQueryBuilder) WithSetID(id string) *manifestQueryBuilder {
	mm.tx.Where("sets_workloads.device_set_id = ?", id)
	return mm
}

func (mm *manifestQueryBuilder) Build() *gorm.DB {
	return mm.tx
}
