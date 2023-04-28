package mappers

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	reader "github.com/tupyy/tinyedge-controller/internal/repo/manifest"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
	"go.uber.org/zap"
)

func ManifestToEntity(mm []models.ManifestJoin) entity.Manifest {
	if kind(mm[0].RefType) == entity.WorkloadManifestKind {
		return createWorkload(mm)
	}
	return nil
}

func ManifestsToEntities(mm []models.ManifestJoin) []entity.Manifest {
	entities := make(map[string][]models.ManifestJoin)
	// makes sure that we add only once the id of the devices, sets or namespaces
	for _, m := range mm {
		list, ok := entities[m.ID]
		if !ok {
			entities[m.ID] = make([]models.ManifestJoin, 0, len(mm))
			continue
		}
		list = append(list, m)
		entities[m.ID] = list
	}

	ee := make([]entity.Manifest, 0, len(entities))
	for _, v := range entities {
		ee = append(ee, createWorkload(v))
	}
	return ee
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

func createWorkload(mm []models.ManifestJoin) entity.Workload {
	m := mm[0]

	e := entity.Workload{
		ObjectMeta: entity.ObjectMeta{
			Id: m.ID,
		},
		TypeMeta: entity.TypeMeta{
			Kind: kind(m.RefType),
		},
		Repository: entity.Repository{
			Id:  m.Repo_ID,
			Url: m.Repo_URL,
		},
		Devices:    make([]string, 0, len(mm)),
		Namespaces: make([]string, 0, len(mm)),
		Sets:       make([]string, 0, len(mm)),
	}
	if m.Repo_Branch.Valid {
		e.Repository.Branch = m.Repo_Branch.String
	}
	if m.Repo_LocalPath.Valid {
		e.Repository.LocalPath = m.Repo_LocalPath.String
	}
	if m.Repo_CurrentHeadSha.Valid {
		e.Repository.CurrentHeadSha = m.Repo_CurrentHeadSha.String
	}
	if m.Repo_TargetHeadSha.Valid {
		e.Repository.TargetHeadSha = m.Repo_TargetHeadSha.String
	}
	if m.Repo_PullPeriodSeconds.Valid {
		e.Repository.PullPeriod = time.Duration(m.Repo_PullPeriodSeconds.Int64) * time.Second
	}

	// makes sure that we add only once the id of the devices, sets or namespaces
	idMap := make(uniqueIds)

	for _, m := range mm {
		if m.DeviceId != "" && !idMap.exists(m.DeviceId, "device") {
			e.Devices = append(e.Devices, m.DeviceId)
			idMap.add(m.DeviceId, "device")
		}
		if m.NamespaceId != "" && !idMap.exists(m.NamespaceId, "namespace") {
			e.Namespaces = append(e.Namespaces, m.NamespaceId)
			idMap.add(m.NamespaceId, "namespace")
		}
		if m.SetId != "" && !idMap.exists(m.SetId, "set") {
			e.Sets = append(e.Sets, m.SetId)
			idMap.add(m.SetId, "set")
		}
	}

	content, err := os.ReadFile(path.Join(e.Repository.LocalPath, mm[0].Path))
	if err != nil {
		zap.S().Errorf("unable to read manifest %q from repository %q: %w", mm[0].Path, e.Repository.LocalPath, err)
		return e
	}

	// read manifest from repo to get the rest of stuff
	repoManifest, err := reader.ReadManifest(bytes.NewReader(content))
	if err != nil {
		zap.S().Errorf("unable to parse manifest %q content: %w", mm[0].Path, err)
		return e
	}

	w, ok := repoManifest.(entity.Workload)
	if !ok {
		zap.S().Errorf("manifest %q is not a workload", mm[0].Path)
	}

	e.Rootless = w.Rootless
	e.Resources = w.Resources
	e.Selectors = w.Selectors
	e.Secrets = w.Secrets
	e.Version = w.Version

	return e
}

func kind(refType string) entity.ManifestKind {
	switch refType {
	case "configuration":
		return entity.ConfigurationManifestKind
	case "workload":
		return entity.WorkloadManifestKind
	default:
		return -1
	}
}

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
