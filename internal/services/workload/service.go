package workload

import (
	"context"
	"errors"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
	"go.uber.org/zap"
)

type WorkloadService struct {
	pgManifestRepo common.ManifestReaderWriter
	pgDeviceRepo   common.DeviceReader
	gitRepo        common.GitReaderWriter
}

func New(pgDeviceRepo common.DeviceReader, pgManifestRepo common.ManifestReaderWriter, gitRepo common.GitReaderWriter) *WorkloadService {
	return &WorkloadService{
		pgDeviceRepo:   pgDeviceRepo,
		pgManifestRepo: pgManifestRepo,
		gitRepo:        gitRepo,
	}
}

func (w *WorkloadService) GetRepositories(ctx context.Context) ([]entity.Repository, error) {
	repos, err := w.pgManifestRepo.GetRepositories(ctx)
	if err != nil {
		return []entity.Repository{}, err
	}
	return repos, nil

}

func (w *WorkloadService) GetManifests(ctx context.Context, repo entity.Repository) ([]entity.ManifestWorkV1, error) {
	return w.pgManifestRepo.GetRepoManifests(ctx, repo)
}

func (w *WorkloadService) UpdateRepository(ctx context.Context, r entity.Repository) error {
	repo, err := w.gitRepo.Open(ctx, r)
	if err != nil {
		return err
	}

	// if err := w.gitRepo.Pull(ctx, r); err != nil {
	// 	return err
	// }

	headSha, err := w.gitRepo.GetHeadSha(ctx, repo)
	if err != nil {
		return err
	}

	if repo.TargetHeadSha != headSha {
		repo.TargetHeadSha = headSha
		// update postgres
		if err := w.pgManifestRepo.UpdateRepo(ctx, repo); err != nil {
			return err
		}
	}

	return nil
}

func (w *WorkloadService) UpdateManifestsFromGit(ctx context.Context, repo entity.Repository) error {
	created, deleted, err := w.getManifests(ctx, repo)
	if err != nil {
		return err
	}

	w.insertManifests(ctx, created)
	w.deleteManifests(ctx, deleted)

	return nil
}

// GetUpdatedManifestWorks returns the manifest work which had been updated.
func (w *WorkloadService) GetUpdatedManifests(ctx context.Context, repos []entity.Repository) ([]entity.ManifestWorkV1, error) {
	updatedManifestWork := make([]entity.ManifestWorkV1, 0, 10)
	for _, r := range repos {
		if r.CurrentHeadSha != r.TargetHeadSha {
			manifestWorks, err := w.gitRepo.GetManifests(ctx, r)
			if err != nil {
				zap.S().Error("unable to get manifest works from repo", "error", err, "repo_id", r.Id, "url", r.Url)
				continue
			}
			updatedManifestWork = append(updatedManifestWork, manifestWorks...)
		}
	}

	return updatedManifestWork, nil
}

func (w *WorkloadService) CreateRelations(ctx context.Context, m entity.ManifestWorkV1) error {
	for _, name := range m.Selector.Namespaces {
		namespace, err := w.pgDeviceRepo.GetNamespace(ctx, name)
		if err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				zap.S().Warnw("unable to create relation. namespace does not exist", "namespace", name)
				continue
			}
			zap.S().Errorw("unable to get namespace", "namespace", name, "error", err)
			continue
		}
		if contains(namespace.ManifestIDS, m.Id) {
			continue
		}
		if err := w.pgManifestRepo.CreateNamespaceRelation(ctx, namespace.Name, m.Id); err != nil {
			zap.S().Errorw("unable to create namespace manifest relation", "error", err, "namespace", namespace.Name, "manifest_id", m.Id)
		}
	}

	for _, name := range m.Selector.Sets {
		set, err := w.pgDeviceRepo.GetSet(ctx, name)
		if err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				zap.S().Warnw("unable to create relation. set does not exist", "set", name)
				continue
			}
			zap.S().Errorw("unable to get set", "set", name, "error", err)
			continue
		}
		if contains(set.ManifestIDS, m.Id) {
			continue
		}
		if err := w.pgManifestRepo.CreateSetRelation(ctx, set.Name, m.Id); err != nil {
			zap.S().Errorw("unable to create set manifest relation", "error", err, "set", set.Name, "manifest_id", m.Id)
		}
	}

	for _, id := range m.Selector.Devices {
		device, err := w.pgDeviceRepo.GetDevice(ctx, id)
		if err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				zap.S().Warnw("unable to create relation. device does not exist", "device_id", id)
				continue
			}
			zap.S().Errorw("unable to get device", "device", id, "error", err)
			continue
		}
		if err := w.pgManifestRepo.CreateDeviceRelation(ctx, device.ID, m.Id); err != nil {
			zap.S().Errorw("unable to create device manifest relation", "error", err, "device", device.ID, "manifest_id", m.Id)
		}
	}
	return nil
}

func (w *WorkloadService) getManifests(ctx context.Context, repo entity.Repository) (created []entity.ManifestWorkV1, deleted []entity.ManifestWorkV1, err error) {
	created = make([]entity.ManifestWorkV1, 0)
	deleted = make([]entity.ManifestWorkV1, 0)

	pgManifest, err := w.pgManifestRepo.GetRepoManifests(ctx, repo)
	if err != nil {
		return created, deleted, err
	}

	gitManifests, err := w.gitRepo.GetManifests(ctx, repo)
	if err != nil {
		return created, deleted, err
	}

	created = substract(gitManifests, pgManifest)
	deleted = substract(pgManifest, gitManifests)

	return created, deleted, nil
}

func (w *WorkloadService) insertManifests(ctx context.Context, manifests []entity.ManifestWorkV1) {
	for _, m := range manifests {
		if err := w.pgManifestRepo.InsertManifest(ctx, m); err != nil {
			zap.S().Errorw("unable to insert manifest", "errro", err, "manifest_id", m.Id, "manifest_repo", m.Repo.LocalPath)
			continue
		}
	}
}

func (w *WorkloadService) deleteManifests(ctx context.Context, manifests []entity.ManifestWorkV1) {
	for _, m := range manifests {
		if err := w.pgManifestRepo.DeleteManifest(ctx, m); err != nil {
			zap.S().Error("unable to delete manifest", "errro", err, "manifest_id", m.Id, "manifest_repo", m.Repo.LocalPath)
			continue
		}
	}
}

// substract return all elements of a which are not found in b
func substract(a []entity.ManifestWorkV1, b []entity.ManifestWorkV1) []entity.ManifestWorkV1 {
	m1 := make(map[string]entity.ManifestWorkV1)
	m2 := make(map[string]entity.ManifestWorkV1)

	limit := len(a)
	if limit < len(b) {
		limit = len(b)
	}

	for i := 0; i < limit; i++ {
		if i < len(a) {
			m1[a[i].Id] = a[i]
		}

		if i < len(b) {
			m2[b[i].Id] = b[i]
		}
	}

	res := make([]entity.ManifestWorkV1, 0, len(a))
	for id, v := range m1 {
		if _, found := m2[id]; !found {
			res = append(res, v)
		}
	}

	return res
}

func contains(arr []string, id string) bool {
	for _, a := range arr {
		if a == id {
			return true
		}
	}
	return false
}
