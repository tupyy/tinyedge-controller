package mappers

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
)

type uniqueIds map[string]struct{}

func (u uniqueIds) exists(id string, prefix string) bool {
	_id := fmt.Sprintf("%s%s", prefix, id)
	_, ok := u[_id]
	return ok
}

func (u uniqueIds) add(id string, prefix string) {
	_id := fmt.Sprintf("%s%s", prefix, id)
	u[_id] = struct{}{}
}

func ManifestModelToEntity(mm []models.ManifestJoin) entity.ManifestReference {
	m := mm[0]
	e := entity.ManifestReference{
		Id:           m.ID,
		Valid:        m.Valid,
		Hash:         m.Hash,
		Path:         m.PathManifestWork,
		DeviceIDs:    make([]string, 0, len(mm)),
		SetIDs:       make([]string, 0, len(mm)),
		NamespaceIDs: make([]string, 0, len(mm)),
		Repo: entity.Repository{
			Id:  m.Repo_ID,
			Url: m.Repo_URL,
		},
	}
	if m.Repo_Branch.Valid {
		e.Repo.Branch = m.Repo_Branch.String
	}
	if m.Repo_LocalPath.Valid {
		e.Repo.LocalPath = m.Repo_LocalPath.String
	}
	if m.Repo_CurrentHeadSha.Valid {
		e.Repo.CurrentHeadSha = m.Repo_CurrentHeadSha.String
	}
	if m.Repo_TargetHeadSha.Valid {
		e.Repo.TargetHeadSha = m.Repo_TargetHeadSha.String
	}
	if m.Repo_PullPeriodSeconds.Valid {
		e.Repo.PullPeriod = time.Duration(m.Repo_PullPeriodSeconds.Int64) * time.Second
	}

	// makes sure that we add only once the id of the devices, sets or namespaces
	idMap := make(uniqueIds)

	for _, m := range mm {
		if m.DeviceId != "" && !idMap.exists(m.DeviceId, "device") {
			e.DeviceIDs = append(e.DeviceIDs, m.DeviceId)
			idMap.add(m.DeviceId, "device")
		}
		if m.NamespaceId != "" && !idMap.exists(m.NamespaceId, "namespace") {
			e.NamespaceIDs = append(e.NamespaceIDs, m.NamespaceId)
			idMap.add(m.NamespaceId, "namespace")
		}
		if m.SetId != "" && !idMap.exists(m.SetId, "set") {
			e.SetIDs = append(e.SetIDs, m.SetId)
			idMap.add(m.SetId, "set")
		}
	}
	return e
}

func ManifestModelsToEntities(mm []models.ManifestJoin) []entity.ManifestReference {
	entities := make(map[string]entity.ManifestReference)
	// makes sure that we add only once the id of the devices, sets or namespaces
	idMap := make(uniqueIds)
	for _, m := range mm {
		manifest, ok := entities[m.ID]
		if !ok {
			manifest = entity.ManifestReference{
				Id:    m.ID,
				Valid: m.Valid,
				Hash:  m.Hash,
				Path:  m.PathManifestWork,
				Repo: entity.Repository{
					Id:  m.Repo_ID,
					Url: m.Repo_URL,
				},
				DeviceIDs:    make([]string, 0, len(mm)),
				SetIDs:       make([]string, 0, len(mm)),
				NamespaceIDs: make([]string, 0, len(mm)),
			}
			if m.Repo_Branch.Valid {
				manifest.Repo.Branch = m.Repo_Branch.String
			}
			if m.Repo_LocalPath.Valid {
				manifest.Repo.LocalPath = m.Repo_LocalPath.String
			}
			if m.Repo_CurrentHeadSha.Valid {
				manifest.Repo.CurrentHeadSha = m.Repo_CurrentHeadSha.String
			}
			if m.Repo_TargetHeadSha.Valid {
				manifest.Repo.TargetHeadSha = m.Repo_TargetHeadSha.String
			}
			if m.Repo_PullPeriodSeconds.Valid {
				manifest.Repo.PullPeriod = time.Duration(m.Repo_PullPeriodSeconds.Int64) * time.Second
			}
		}
		if m.DeviceId != "" && !idMap.exists(m.DeviceId, "device") {
			manifest.DeviceIDs = append(manifest.DeviceIDs, m.DeviceId)
			idMap.add(m.DeviceId, "device")
		}
		if m.NamespaceId != "" && !idMap.exists(m.NamespaceId, "namespace") {
			manifest.NamespaceIDs = append(manifest.NamespaceIDs, m.NamespaceId)
			idMap.add(m.NamespaceId, "namespace")
		}
		if m.SetId != "" && !idMap.exists(m.SetId, "set") {
			manifest.SetIDs = append(manifest.SetIDs, m.SetId)
			idMap.add(m.SetId, "set")
		}
		entities[m.ID] = manifest
	}

	ee := make([]entity.ManifestReference, 0, len(entities))
	for _, v := range entities {
		ee = append(ee, v)
	}
	return ee
}

func ManifestEntityToModel(e entity.ManifestReference) models.ManifestWork {
	m := models.ManifestWork{
		ID:               e.Id,
		PathManifestWork: e.Path,
		RepoID:           e.Repo.Id,
		Valid:            e.Valid,
		Hash:             e.Hash,
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
