package mappers

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	goyaml "github.com/go-yaml/yaml"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	manifestv1 "github.com/tupyy/tinyedge-controller/internal/repo/models/manifest/v1"
	models "github.com/tupyy/tinyedge-controller/internal/repo/models/pg"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
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

func ManifestModelToEntity(mm []models.ManifestJoin) (entity.ManifestWork, error) {
	m := mm[0]
	e, err := basicManifestModelToEntity(m)
	if err != nil {
		return entity.ManifestWork{}, err
	}

	e.DeviceIDs = make([]string, 0, len(mm))
	e.SetIDs = make([]string, 0, len(mm))
	e.NamespaceIDs = make([]string, 0, len(mm))

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
	return e, err
}

func ManifestModelsToEntities(mm []models.ManifestJoin) ([]entity.ManifestWork, error) {
	entities := make(map[string]entity.ManifestWork)
	// makes sure that we add only once the id of the devices, sets or namespaces
	idMap := make(uniqueIds)
	for _, m := range mm {
		manifest, ok := entities[m.ID]
		if !ok {
			e, err := basicManifestModelToEntity(m)
			if err != nil {
				zap.S().Errorw("unable to map basic manifest to entity", "error", err, "manifest join", m)
			}
			manifest = e
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

	ee := make([]entity.ManifestWork, 0, len(entities))
	for _, v := range entities {
		ee = append(ee, v)
	}
	return ee, nil
}

func basicManifestModelToEntity(m models.ManifestJoin) (entity.ManifestWork, error) {
	content, err := os.ReadFile(m.PathManifestWork)
	if err != nil {
		return entity.ManifestWork{}, fmt.Errorf("unable to read manifest file %q from repo: %w", m.PathManifestWork, err)
	}
	gitModel := manifestv1.Manifest{}
	if err := yaml.Unmarshal(content, &gitModel); err != nil {
		return entity.ManifestWork{}, fmt.Errorf("unable to unmarshal content of file %q: %w", m.PathManifestWork, err)
	}

	e := entity.ManifestWork{
		ConfigMaps: make([]v1.ConfigMap, 0),
		Pods:       make([]v1.Pod, 0),
	}
	e.Id = m.ID
	e.Repo = entity.Repository{
		Id:  m.RepoID,
		Url: m.Repo_URL,
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
	// create resources
	for _, resource := range gitModel.Resources {
		configmaps, pods, err := createResources(resource, m.Repo_LocalPath.String)
		if err != nil {
			return entity.ManifestWork{}, err
		}
		e.ConfigMaps = append(e.ConfigMaps, configmaps...)
		e.Pods = append(e.Pods, pods...)
	}

	e.Path = m.PathManifestWork
	e.DeviceIDs = make([]string, 0)
	e.NamespaceIDs = make([]string, 0)
	e.SetIDs = make([]string, 0)
	e.Hash = m.Hash
	e.Path = m.PathManifestWork
	return e, nil
}

func createResources(resource manifestv1.Resource, basePath string) ([]v1.ConfigMap, []v1.Pod, error) {
	content, err := os.ReadFile(path.Join(basePath, resource.Ref))
	if err != nil {
		return []v1.ConfigMap{}, []v1.Pod{}, fmt.Errorf("unable to read file %q: %w", resource.Ref, err)
	}

	parts, err := splitYAML(content)
	if err != nil {
		return []v1.ConfigMap{}, []v1.Pod{}, fmt.Errorf("unable to decode resource file %q: %w", resource.Ref, err)
	}

	configMaps := make([]v1.ConfigMap, 0)
	pods := make([]v1.Pod, 0)
	allowedKinds := "ConfigMap|Pods"
	for _, part := range parts {
		kind, err := getKind(part)
		if err != nil {
			zap.S().Errorf("unable to get \"kind\" from yaml with error %q", err)
			continue
		}
		if kind == "" || !strings.Contains(allowedKinds, kind) {
			zap.S().Errorf("kind %q not allowed in manifest work", kind)
			continue
		}
		switch kind {
		case "ConfigMap":
			var c v1.ConfigMap
			err := yaml.Unmarshal(part, &c)
			if err != nil {
				return []v1.ConfigMap{}, []v1.Pod{}, fmt.Errorf("unable to unmarshal part %q: %v", string(part), err)
			}
			configMaps = append(configMaps, c)
		case "Pod":
			var p v1.Pod
			err := yaml.Unmarshal(part, &p)
			if err != nil {
				return []v1.ConfigMap{}, []v1.Pod{}, fmt.Errorf("unable to unmarshal part %q: %v", string(part), err)
			}
			pods = append(pods, p)
		}
	}

	return configMaps, pods, nil
}

func getKind(content []byte) (string, error) {
	type anonymousStruct struct {
		Kind string `yaml:"kind"`
	}
	var a anonymousStruct
	if err := goyaml.Unmarshal(content, &a); err != nil {
		return "", fmt.Errorf("unknown struct: %s", err)
	}
	return a.Kind, nil
}

func splitYAML(resources []byte) ([][]byte, error) {
	dec := goyaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := goyaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}

func ManifestEntityToModel(e entity.ManifestWork) models.ManifestWork {
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
