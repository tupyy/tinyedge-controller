package postgres

import (
	"context"
	"errors"

	pgclient "github.com/tupyy/tinyedge-controller/internal/clients/pg"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/repo/models/mappers"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	db             *gorm.DB
	client         pgclient.Client
	circuitBreaker pgclient.CircuitBreaker
}

func NewRepository(client pgclient.Client) (*Repository, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &Repository{}, err
	}

	return &Repository{gormDB, client, client.GetCircuitBreaker()}, nil
}

func (m *Repository) GetRepository(ctx context.Context, id string) (entity.Repository, error) {
	if !m.circuitBreaker.IsAvailable() {
		return entity.Repository{}, errService.NewPostgresNotAvailableError("repository")
	}

	repo := models.Repo{}

	if err := m.getDb(ctx).Where("id = ?", id).First(&repo).Error; err != nil {
		if m.checkNetworkError(err) {
			return entity.Repository{}, errService.NewPostgresNotAvailableError("repository")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Repository{}, errService.NewResourceNotFoundError("repository", id)
		}
		return entity.Repository{}, err
	}

	return mappers.RepoModelToEntity(repo), nil
}

func (m *Repository) GetRepositories(ctx context.Context) ([]entity.Repository, error) {
	if !m.circuitBreaker.IsAvailable() {
		return []entity.Repository{}, errService.NewPostgresNotAvailableError("repository")
	}

	repos := []models.Repo{}

	if err := m.getDb(ctx).Find(&repos).Error; err != nil {
		if m.checkNetworkError(err) {
			return []entity.Repository{}, errService.NewPostgresNotAvailableError("repository")
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

func (m *Repository) InsertRepository(ctx context.Context, r entity.Repository) error {
	if !m.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("repository")
	}

	model := mappers.RepoEntityToModel(r)

	if err := m.getDb(ctx).Create(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return errService.NewPostgresNotAvailableError("repository")
		}
		return err
	}
	return nil
}

func (m *Repository) UpdateRepository(ctx context.Context, r entity.Repository) error {
	if !m.circuitBreaker.IsAvailable() {
		return errService.NewPostgresNotAvailableError("repository")
	}

	model := mappers.RepoEntityToModel(r)

	if err := m.getDb(ctx).Where("id = ?", model.ID).Save(&model).Error; err != nil {
		if m.checkNetworkError(err) {
			return errService.NewPostgresNotAvailableError("repository")
		}
		return err
	}

	return nil
}

func (d *Repository) checkNetworkError(err error) (isOpen bool) {
	isOpen = d.circuitBreaker.BreakOnNetworkError(err)
	if isOpen {
		zap.S().Warn("circuit breaker is now open")
	}
	return
}

func (d *Repository) getDb(ctx context.Context) *gorm.DB {
	return d.db.Session(&gorm.Session{SkipHooks: true}).WithContext(ctx)
}
