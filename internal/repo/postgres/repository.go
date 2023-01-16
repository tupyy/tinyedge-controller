package postgres

import (
	"context"

	pgclient "github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres/mappers"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
	"gorm.io/gorm"
)

type Repository struct {
	db             *gorm.DB
	client         pgclient.Client
	circuitBreaker pgclient.CircuitBreaker
}

func NewRepository(client pgclient.Client) (*ReferenceRepository, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &ReferenceRepository{}, err
	}

	return &ReferenceRepository{gormDB, client, client.GetCircuitBreaker()}, nil
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
