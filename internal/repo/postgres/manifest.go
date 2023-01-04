package postgres

import (
	"context"

	pgclient "github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres/mappers"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres/models"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ManifestRepo struct {
	db             *gorm.DB
	client         pgclient.Client
	circuitBreaker pgclient.CircuitBreaker
}

func NewManifestRepo(client pgclient.Client) (*ManifestRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &ManifestRepo{}, err
	}

	return &ManifestRepo{gormDB, client, client.GetCircuitBreaker()}, nil
}

func (m *ManifestRepo) GetManifests(ctx context.Context) ([]entity.ManifestWorkV1, error) {
	if !m.circuitBreaker.IsAvailable() {
		return []entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
	}

	manifests := []models.ManifestJoin{}

	tx := newManifestQuery(ctx, m.db).Build()
	if err := tx.Find(&manifests).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
		}
		return []entity.ManifestWorkV1{}, err
	}

	if len(manifests) == 0 {
		return []entity.ManifestWorkV1{}, nil
	}

	e, err := mappers.ManifestModelsToEntities(manifests)
	if err != nil {
		return []entity.ManifestWorkV1{}, err
	}

	return e, nil

}

func (m *ManifestRepo) GetManifest(ctx context.Context, id string) (entity.ManifestWorkV1, error) {
	if !m.circuitBreaker.IsAvailable() {
		return entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
	}

	manifests := []models.ManifestJoin{}

	tx := newManifestQuery(ctx, m.db).WithManifestID(id).Build()
	if err := tx.Find(&manifests).Error; err != nil {
		if m.checkNetworkError(err) {
			return entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
		}
		return entity.ManifestWorkV1{}, err
	}

	if len(manifests) == 0 {
		return entity.ManifestWorkV1{}, common.ErrResourceNotFound
	}

	e, err := mappers.ManifestModelToEntity(manifests)
	if err != nil {
		return entity.ManifestWorkV1{}, err
	}

	return e, nil
}

func (m *ManifestRepo) GetRepositories(ctx context.Context) ([]entity.Repository, error) {
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

func (m *ManifestRepo) GetRepoManifests(ctx context.Context, r entity.Repository) ([]entity.ManifestWorkV1, error) {
	if !m.circuitBreaker.IsAvailable() {
		return []entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
	}

	manifests := []models.ManifestJoin{}

	tx := newManifestQuery(ctx, m.db).WithRepoId(r.Id).Build()
	if err := tx.Find(&manifests).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
		}
		return []entity.ManifestWorkV1{}, err
	}

	if len(manifests) == 0 {
		return []entity.ManifestWorkV1{}, nil
	}

	e, err := mappers.ManifestModelsToEntities(manifests)
	if err != nil {
		return []entity.ManifestWorkV1{}, err
	}

	return e, nil
}

func (m *ManifestRepo) InsertRepo(ctx context.Context, r entity.Repository) error {
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

func (m *ManifestRepo) UpdateRepo(ctx context.Context, r entity.Repository) error {
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

func (m *ManifestRepo) InsertManifest(ctx context.Context, manifest entity.ManifestWorkV1) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := mappers.ManifestEntityToModel(manifest)

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}
	return nil
}

func (m *ManifestRepo) UpdateManifest(ctx context.Context, manifest entity.ManifestWorkV1) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := mappers.ManifestEntityToModel(manifest)

	if err := m.getDb(ctx).Where("id = ?", model.ID).Save(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ManifestRepo) DeleteManifest(ctx context.Context, manifest entity.ManifestWorkV1) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	if err := m.getDb(ctx).Where("id = ?", manifest.Id).Delete(&models.ManifestWork{}).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ManifestRepo) GetNamespaceManifests(ctx context.Context, namespaceID string) ([]entity.ManifestWorkV1, error) {
	if !m.circuitBreaker.IsAvailable() {
		return []entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
	}

	models := []models.ManifestJoin{}
	tx := newManifestQuery(ctx, m.db).WithNamespaceID(namespaceID).Build()
	if err := tx.Find(&models).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
		}
		return []entity.ManifestWorkV1{}, err
	}

	if len(models) == 0 {
		return []entity.ManifestWorkV1{}, nil
	}

	e, err := mappers.ManifestModelsToEntities(models)
	if err != nil {
		return []entity.ManifestWorkV1{}, err
	}

	return e, nil
}

func (m *ManifestRepo) GetSetManifests(ctx context.Context, setID string) ([]entity.ManifestWorkV1, error) {
	if !m.circuitBreaker.IsAvailable() {
		return []entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
	}

	models := []models.ManifestJoin{}

	tx := newManifestQuery(ctx, m.db).WithSetID(setID).Build()
	if err := tx.Find(&models).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
		}
		return []entity.ManifestWorkV1{}, err
	}

	if len(models) == 0 {
		return []entity.ManifestWorkV1{}, nil
	}

	e, err := mappers.ManifestModelsToEntities(models)
	if err != nil {
		return []entity.ManifestWorkV1{}, err
	}

	return e, nil
}

func (m *ManifestRepo) GetDeviceManifests(ctx context.Context, deviceID string) ([]entity.ManifestWorkV1, error) {
	if !m.circuitBreaker.IsAvailable() {
		return []entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
	}

	models := []models.ManifestJoin{}

	tx := newManifestQuery(ctx, m.db).WithDeviceID(deviceID).Build()
	if err := tx.Find(&models).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
		}
		return []entity.ManifestWorkV1{}, err
	}

	if len(models) == 0 {
		return []entity.ManifestWorkV1{}, nil
	}

	e, err := mappers.ManifestModelsToEntities(models)
	if err != nil {
		return []entity.ManifestWorkV1{}, err
	}

	return e, nil
}

func (m *ManifestRepo) CreateNamespaceRelation(ctx context.Context, namespaceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := models.NamespacesWorkloads{
		NamespaceID:    namespaceID,
		ManifestWorkID: manifestID,
	}

	// check if the relation already exists
	var dummy models.NamespacesWorkloads
	if err := m.getDb(ctx).Where("namespace_id = ? AND manifest_work_id = ?", namespaceID, manifestID).First(&dummy).Error; err == nil {
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

func (m *ManifestRepo) DeleteNamespaceRelation(ctx context.Context, namespaceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := models.NamespacesWorkloads{}
	if err := m.getDb(ctx).Where("namespace_id = ? AND manifest_work_id = ?", namespaceID, manifestID).Delete(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ManifestRepo) CreateSetRelation(ctx context.Context, setID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := models.SetsWorkloads{
		DeviceSetID:    setID,
		ManifestWorkID: manifestID,
	}

	// check if the relation already exists
	var dummy models.SetsWorkloads
	if err := m.getDb(ctx).Where("device_set_id = ? AND manifest_work_id = ?", setID, manifestID).First(&dummy).Error; err == nil {
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

func (m *ManifestRepo) DeleteSetRelation(ctx context.Context, setID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := models.SetsWorkloads{}
	if err := m.getDb(ctx).Where("device_set_id = ? AND manifest_work_id = ?", setID, manifestID).Delete(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ManifestRepo) CreateDeviceRelation(ctx context.Context, deviceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := models.DevicesWorkloads{
		DeviceID:       deviceID,
		ManifestWorkID: manifestID,
	}

	// check if the relation already exists
	var dummy models.DevicesWorkloads
	if err := m.getDb(ctx).Where("device_id = ? AND manifest_work_id = ?", deviceID, manifestID).First(&dummy).Error; err == nil {
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

func (m *ManifestRepo) DeleteDeviceRelation(ctx context.Context, deviceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := models.DevicesWorkloads{}
	if err := m.getDb(ctx).Where("device_id = ? AND manifest_work_id = ?", deviceID, manifestID).Delete(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return common.ErrPostgresNotAvailable
		}
		return err
	}

	return nil
}

func (m *ManifestRepo) checkNetworkError(err error) (isOpen bool) {
	isOpen = m.circuitBreaker.BreakOnNetworkError(err)
	if isOpen {
		zap.S().Warn("circuit breaker is now open")
	}
	return
}

func (m *ManifestRepo) getDb(ctx context.Context) *gorm.DB {
	return m.db.Session(&gorm.Session{SkipHooks: true}).WithContext(ctx)
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
	tx := db.Session(&gorm.Session{SkipHooks: true}).WithContext(ctx).Table("manifest_work").
		Select(`manifest_work.*, devices_workloads.device_id as device_id, sets_workloads.device_set_id as set_id, namespaces_workloads.namespace_id as namespace_id,
		repo.id as repo_id, repo.url as repo_url, repo.branch as repo_branch, repo.local_path as repo_local_path,
		repo.current_head_sha as repo_current_head_sha, repo.target_head_sha as repo_target_head_sha,
		repo.pull_period_seconds as repo_pull_period_seconds`).
		Joins("LEFT JOIN namespaces_workloads ON namespaces_workloads.manifest_work_id = manifest_work.id").
		Joins("LEFT JOIN sets_workloads ON sets_workloads.manifest_work_id = manifest_work.id").
		Joins("LEFT JOIN devices_workloads ON devices_workloads.manifest_work_id = manifest_work.id").
		Joins("JOIN repo ON repo.id = manifest_work.repo_id")
	return &manifestQueryBuilder{tx}
}

func (mm *manifestQueryBuilder) WithRepoId(id string) *manifestQueryBuilder {
	mm.tx.Where("repo_id = ?", id)
	return mm
}

func (mm *manifestQueryBuilder) WithManifestID(id string) *manifestQueryBuilder {
	mm.tx.Where("manifest_work.id = ?", id)
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
