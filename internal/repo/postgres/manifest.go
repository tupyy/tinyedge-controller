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

	manifests := []models.ManifestWork{}

	if err := m.getDb(ctx).Find(&manifests).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
		}
		return []entity.ManifestWorkV1{}, err
	}

	if len(manifests) == 0 {
		return []entity.ManifestWorkV1{}, nil
	}

	entities := make([]entity.ManifestWorkV1, 0, len(manifests))
	for _, m := range manifests {
		e, err := mappers.ManifestModelToEntity(m)
		if err != nil {
			zap.S().Errorw("unable to map manifest model to entity", "error", err, "content", m.Content)
			continue
		}
		entities = append(entities, e)
	}

	return entities, nil
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

	manifests := []models.ManifestWork{}

	if err := m.getDb(ctx).Where("repo_id = ?", r.Id).Find(&manifests).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.ManifestWorkV1{}, common.ErrPostgresNotAvailable
		}
		return []entity.ManifestWorkV1{}, err
	}

	if len(manifests) == 0 {
		return []entity.ManifestWorkV1{}, nil
	}

	entities := make([]entity.ManifestWorkV1, 0, len(manifests))
	for _, m := range manifests {
		e, err := mappers.ManifestModelToEntity(m)
		if err != nil {
			zap.S().Errorw("unable to map manifest model to entity", "error", err, "content", m.Content)
			continue
		}
		entities = append(entities, e)
	}

	return entities, nil

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

func (m *ManifestRepo) CreateNamespaceRelation(ctx context.Context, namespaceID, manifestID string) error {
	if !m.circuitBreaker.IsAvailable() {
		return common.ErrPostgresNotAvailable
	}

	model := models.NamespacesWorkloads{
		NamespaceID:    namespaceID,
		ManifestWorkID: manifestID,
	}

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
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

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
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

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
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
