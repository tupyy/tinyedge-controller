package mappers

import (
	"database/sql"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/repo/postgres/models"
)

func ManifestModelToEntity(m models.ManifestWork) (entity.ManifestWorkV1, error) {
	e, err := entity.ManifestWorkV1{}.Decode(m.Content)
	return e, err
}

func ManifestEntityToModel(e entity.ManifestWorkV1) models.ManifestWork {
	m := models.ManifestWork{
		ID:               e.Id,
		PathManifestWork: e.Path,
		RepoID:           e.Repo.Id,
		Content:          e.Encode(),
	}
	return m
}

func RepoEntityToModel(r entity.Repository) models.Repo {
	m := models.Repo{
		ID:                r.Id,
		URL:               r.Url,
		PullPeriodSeconds: sql.NullInt64{Valid: true, Int64: int64(r.PullPeriod.Seconds())},
	}

	if r.CurrentHeadSha != "" {
		m.CurrentHeadSha = sql.NullString{Valid: true, String: r.CurrentHeadSha}
	}

	if r.Branch != "" {
		m.Branch = sql.NullString{Valid: true, String: r.Branch}
	}

	if r.LocalPath != "" {
		m.LocalPath = sql.NullString{Valid: true, String: r.LocalPath}
	}

	if r.TargetHeadSha != "" {
		m.TargetHeadSha = sql.NullString{Valid: true, String: r.TargetHeadSha}
	}

	return m
}

func RepoModelToEntity(m models.Repo) entity.Repository {
	e := entity.Repository{
		Id:         m.ID,
		Url:        m.URL,
		PullPeriod: 20 * time.Second,
	}

	if m.CurrentHeadSha.Valid {
		e.CurrentHeadSha = m.CurrentHeadSha.String
	}

	if m.PullPeriodSeconds.Valid {
		e.PullPeriod = time.Duration(m.PullPeriodSeconds.Int64 * int64(time.Second))
	}

	if m.LocalPath.Valid {
		e.LocalPath = m.LocalPath.String
	}

	if m.Branch.Valid {
		e.Branch = m.Branch.String
	}

	if m.TargetHeadSha.Valid {
		e.TargetHeadSha = m.TargetHeadSha.String
	}

	return e
}
