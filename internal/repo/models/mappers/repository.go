package mappers

import (
	"database/sql"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
)

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

	if r.AuthType != entity.NoRepositoryAuthType {
		switch r.AuthType {
		case entity.SSHRepositoryAuthType:
			m.AuthType = sql.NullString{Valid: true, String: "ssh"}
		case entity.TokenRepositoryAuthType:
			m.AuthType = sql.NullString{Valid: true, String: "token"}
		case entity.BasicRepositoryAuthType:
			m.AuthType = sql.NullString{Valid: true, String: "basic"}
		}
	}

	if r.CredentialsSecretPath != "" {
		m.AuthSecretPath = sql.NullString{Valid: true, String: r.CredentialsSecretPath}
	}

	return m
}

func RepoModelToEntity(m models.Repo) entity.Repository {
	e := entity.Repository{
		Id:         m.ID,
		Url:        m.URL,
		PullPeriod: 20 * time.Second,
		AuthType:   entity.NoRepositoryAuthType,
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

	if m.AuthSecretPath.Valid {
		e.CredentialsSecretPath = m.AuthSecretPath.String
	}

	if m.AuthType.Valid {
		switch m.AuthType.String {
		case "ssh":
			e.AuthType = entity.SSHRepositoryAuthType
		case "token":
			e.AuthType = entity.TokenRepositoryAuthType
		case "basic":
			e.AuthType = entity.BasicRepositoryAuthType
		default:
			e.AuthType = entity.NoRepositoryAuthType
		}
	}

	return e
}
