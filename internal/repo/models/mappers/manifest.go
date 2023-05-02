package mappers

import (
	"errors"
	"path"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/repo/manifest"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
)

func ManifestToEntity(mm []models.ManifestJoin, readFn manifest.ManifestReader) (entity.Manifest, error) {
	return parseManifest(mm, readFn)
}

func ManifestsToEntities(mm []models.ManifestJoin, readFn manifest.ManifestReader) ([]entity.Manifest, error) {
	entities := make(map[string][]models.ManifestJoin)
	// makes sure that we add only once the id of the devices, sets or namespaces
	for _, m := range mm {
		list, ok := entities[m.ID]
		if !ok {
			entities[m.ID] = make([]models.ManifestJoin, 0, len(mm))
		}
		list = append(list, m)
		entities[m.ID] = list
	}

	ee := make([]entity.Manifest, 0, len(entities))
	for _, v := range entities {
		m, err := parseManifest(v, readFn)
		if err != nil {
			return []entity.Manifest{}, err
		}
		ee = append(ee, m)
	}
	return ee, nil
}

func ManifestEntityToModel(e entity.Manifest) models.Manifest {
	m := models.Manifest{
		ID:      e.GetID(),
		RefType: e.GetKind().String(),
	}
	switch v := e.(type) {
	case entity.Configuration:
	case entity.Workload:
		m.RepoID = v.Repository.Id
		m.Path = v.Path
	}

	return m
}

func parseManifest(mm []models.ManifestJoin, readFn manifest.ManifestReader) (entity.Manifest, error) {
	m := mm[0]

	repo := entity.Repository{}
	repo.Id = m.RepoID
	repo.Url = m.RepoURL
	if m.RepoBranch.Valid {
		repo.Branch = m.RepoBranch.String
	}
	if m.RepoLocalPath.Valid {
		repo.LocalPath = m.RepoLocalPath.String
	}
	if m.RepoCurrentHeadSha.Valid {
		repo.CurrentHeadSha = m.RepoCurrentHeadSha.String
	}
	if m.RepoTargetHeadSha.Valid {
		repo.TargetHeadSha = m.RepoTargetHeadSha.String
	}
	if m.RepoPullPeriodSeconds.Valid {
		repo.PullPeriod = time.Duration(m.RepoPullPeriodSeconds.Int64) * time.Second
	}
	if m.RepoAuth.Valid {
		switch m.RepoAuth.String {
		case "ssh":
			repo.AuthType = entity.SSHRepositoryAuthType
		case "basic":
			repo.AuthType = entity.BasicRepositoryAuthType
		case "token":
			repo.AuthType = entity.TokenRepositoryAuthType
		default:
			repo.AuthType = entity.NoRepositoryAuthType
		}
	}
	if m.RepoAuth.Valid && m.RepoAuthSecret.Valid {
		repo.CredentialsSecretPath = m.RepoAuthSecret.String
	}

	// makes sure that we add only once the id of the devices, sets or namespaces
	idMap := make(uniqueIds)

	devices := make([]string, 0, len(mm))
	sets := make([]string, 0, len(mm))
	namespaces := make([]string, 0, len(mm))
	for _, m := range mm {
		if m.DeviceId != "" && !idMap.exists(m.DeviceId, "device") {
			devices = append(devices, m.DeviceId)
			idMap.add(m.DeviceId, "device")
		}
		if m.NamespaceId != "" && !idMap.exists(m.NamespaceId, "namespace") {
			namespaces = append(namespaces, m.NamespaceId)
			idMap.add(m.NamespaceId, "namespace")
		}
		if m.SetId != "" && !idMap.exists(m.SetId, "set") {
			sets = append(sets, m.SetId)
			idMap.add(m.SetId, "set")
		}
	}

	manifest, err := readManifest(path.Join(repo.LocalPath, mm[0].Path), mm[0].ID, readFn)
	if err != nil {
		return nil, err
	}

	w, ok := manifest.(entity.Workload)
	if ok {
		w.Repository = repo
		w.Devices = devices
		w.Namespaces = namespaces
		w.Sets = sets
		return w, nil
	}

	c, ok := manifest.(entity.Configuration)
	if ok {
		c.Repository = repo
		c.Devices = devices
		c.Namespaces = namespaces
		c.Sets = sets
		return c, nil
	}

	return nil, errors.New("unknown type")
}
