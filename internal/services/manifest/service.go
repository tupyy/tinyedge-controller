package manifest

import (
	"context"
	"errors"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
	"go.uber.org/zap"
)

type Service struct {
	pgReferenceRepo common.ReferenceReaderWriter
	pgDeviceRepo    common.DeviceReader
	gitRepo         common.GitReader
}

func New(pgDeviceRepo common.DeviceReader, pgManifestRepo common.ReferenceReaderWriter, gitRepo common.GitReader) *Service {
	return &Service{
		pgDeviceRepo:    pgDeviceRepo,
		pgReferenceRepo: pgManifestRepo,
		gitRepo:         gitRepo,
	}
}

func (w *Service) GetRepositories(ctx context.Context) ([]entity.Repository, error) {
	repos, err := w.pgReferenceRepo.GetRepositories(ctx)
	if err != nil {
		return []entity.Repository{}, err
	}
	return repos, nil

}

func (w *Service) GetManifest(ctx context.Context, id string) (entity.ManifestWork, error) {
	ref, err := w.pgReferenceRepo.GetReference(ctx, id)
	if err != nil {
		return entity.ManifestWork{}, fmt.Errorf("unable to get manifest reference: %w", err)
	}

	manifest, err := w.gitRepo.GetManifest(ctx, ref)
	if err != nil {
		return entity.ManifestWork{}, fmt.Errorf("unable to get manifest: %w", err)
	}
	return manifest, nil
}

func (w *Service) GetManifestReferences(ctx context.Context, repo entity.Repository) ([]entity.ManifestReference, error) {
	return w.pgReferenceRepo.GetRepositoryReferences(ctx, repo)
}

func (w *Service) GetManifests(ctx context.Context, repo entity.Repository) ([]entity.ManifestWork, error) {
	refs, err := w.GetManifestReferences(ctx, repo)
	if err != nil {
		return []entity.ManifestWork{}, err
	}

	// for each ref get the real manifest and add devices, sets and namespaces
	manifests := make([]entity.ManifestWork, 0, len(refs))
	for _, ref := range refs {
		manifest, err := w.gitRepo.GetManifest(ctx, ref)
		if err != nil {
			return []entity.ManifestWork{}, fmt.Errorf("unable to get manifest %q from repo %q: %w", ref.Path, repo.Id, err)
		}
		manifest.Reference = &ref
		manifests = append(manifests, manifest)
	}

	return manifests, nil
}

func (w *Service) PullRepository(ctx context.Context, repo entity.Repository) (entity.Repository, error) {
	r, err := w.gitRepo.Open(ctx, repo)
	if err != nil {
		return entity.Repository{}, err
	}

	err = w.gitRepo.Pull(ctx, repo)
	if err != nil {
		return entity.Repository{}, err
	}

	headSha, err := w.gitRepo.GetHeadSha(ctx, r)
	if err != nil {
		return entity.Repository{}, err
	}

	r.TargetHeadSha = headSha

	return r, nil
}

func (w *Service) UpdateRepository(ctx context.Context, r entity.Repository) error {
	if err := w.pgReferenceRepo.UpdateRepository(ctx, r); err != nil {
		return err
	}

	return nil
}

func (w *Service) UpdateManifests(ctx context.Context, repo entity.Repository) error {
	references, err := w.pgReferenceRepo.GetRepositoryReferences(ctx, repo)
	if err != nil {
		return fmt.Errorf("unable to read references of repo %q: %w", repo.Id, err)
	}

	manifests, err := w.gitRepo.GetManifests(ctx, repo)
	if err != nil {
		return fmt.Errorf("unable to read manifest from repo %q: %w", repo.Id, err)
	}

	created := substract(manifests, references, func(item entity.ManifestWork) string { return item.Id }, func(item entity.ManifestReference) string { return item.Id })
	deleted := substract(references, manifests, func(item entity.ManifestReference) string { return item.Id }, func(item entity.ManifestWork) string { return item.Id })
	updated := intersect(manifests, references,
		func(item entity.ManifestWork) string { return item.Id },
		func(item entity.ManifestReference) string { return item.Id },
		func(manifest entity.ManifestWork, ref entity.ManifestReference) bool {
			return ref.Hash != manifest.Hash
		},
	)

	for _, c := range created {
		ref := w.createReference(c, repo)
		if err := w.pgReferenceRepo.InsertReference(ctx, ref); err != nil && !errors.Is(err, common.ErrResourceAlreadyExists) {
			return fmt.Errorf("unable to insert repo %q: %w", c.Id, err)
		}
		if err := w.UpdateRelations(ctx, ref); err != nil {
			return err
		}
	}

	for _, d := range deleted {
		if err := w.pgReferenceRepo.DeleteReference(ctx, d); err != nil {
			return fmt.Errorf("unable to delete repo %q: %w", d.Id, err)
		}
	}

	for _, u := range updated {
		if err := w.pgReferenceRepo.UpdateReference(ctx, w.createReference(u, repo)); err != nil {
			return fmt.Errorf("unable to update repo %q: %w", u.Id, err)
		}
		if err := w.UpdateRelations(ctx, w.createReference(u, repo)); err != nil {
			return err
		}
	}

	return nil
}

func (w *Service) UpdateRelations(ctx context.Context, m entity.ManifestReference) error {
	// get the old manifest
	oldManifest, err := w.pgReferenceRepo.GetReference(ctx, m.Id)
	if err != nil {
		return err
	}

	// updateRelation updates the relation between a selector value and the manifest
	updateRelation := func(ctx context.Context, n []string, m string, fn func(ctx context.Context, resourceID, manifestID string) error) error {
		if len(n) > 0 {
			for _, s := range n {
				if err := fn(ctx, s, m); err != nil {
					return fmt.Errorf("unable to create/delete relation between selector %q and manifest %q: %w", s, m, err)
				}
			}
		}
		return nil
	}

	// for each new relation needs to be created (exists in m but not in oldManifest)
	newNamespaceSelectors := substract1(m.NamespaceIDs, oldManifest.NamespaceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, newNamespaceSelectors, m.Id, func(ctx context.Context, namespaceID, manifestID string) error {
		if _, err := w.pgDeviceRepo.GetNamespace(ctx, namespaceID); err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return nil
			}
			return err
		}
		if err := w.pgReferenceRepo.CreateNamespaceRelation(ctx, namespaceID, manifestID); err != nil {
			if !errors.Is(err, common.ErrResourceAlreadyExists) {
				return err
			}
		}
		zap.S().Debugf("relation created between namespace %q and manifest %q", namespaceID, manifestID)
		return nil
	}); err != nil {
		return err
	}

	newSetSelectors := substract1(m.SetIDs, oldManifest.SetIDs, func(i string) string { return i })
	if err := updateRelation(ctx, newSetSelectors, m.Id, func(ctx context.Context, setID, manifestID string) error {
		if _, err := w.pgDeviceRepo.GetSet(ctx, setID); err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return nil
			}
			return err
		}
		if err := w.pgReferenceRepo.CreateSetRelation(ctx, setID, manifestID); err != nil {
			if !errors.Is(err, common.ErrResourceAlreadyExists) {
				return err
			}
		}
		zap.S().Debugf("relation created between set %q and manifest %q", setID, manifestID)
		return nil
	}); err != nil {
		return err
	}

	newDeviceSelectors := substract1(m.DeviceIDs, oldManifest.DeviceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, newDeviceSelectors, m.Id, func(ctx context.Context, deviceID, manifestID string) error {
		if _, err := w.pgDeviceRepo.GetDevice(ctx, deviceID); err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return nil
			}
			return err
		}
		if err := w.pgReferenceRepo.CreateDeviceRelation(ctx, deviceID, manifestID); err != nil {
			if !errors.Is(err, common.ErrResourceAlreadyExists) {
				return err
			}
		}
		zap.S().Debugf("relation created between device %q and manifest %q", deviceID, manifestID)
		return nil
	}); err != nil {
		return err
	}

	// remove the old ones
	oldNamespaceSelectors := substract1(oldManifest.NamespaceIDs, m.NamespaceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, oldNamespaceSelectors, m.Id, func(ctx context.Context, namespaceID, manifestID string) error {
		if _, err := w.pgDeviceRepo.GetNamespace(ctx, namespaceID); err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return nil
			}
		}
		err := w.pgReferenceRepo.DeleteNamespaceRelation(ctx, namespaceID, manifestID)
		if err != nil {
			return err
		}
		zap.S().Debugf("relation deleted between device %q and manifest %q", namespaceID, manifestID)
		return nil
	}); err != nil {
		return err
	}

	oldSetSelectors := substract1(oldManifest.SetIDs, m.SetIDs, func(i string) string { return i })
	if err := updateRelation(ctx, oldSetSelectors, m.Id, func(ctx context.Context, setID, manifestID string) error {
		if _, err := w.pgDeviceRepo.GetSet(ctx, setID); err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return nil
			}
		}
		err := w.pgReferenceRepo.DeleteSetRelation(ctx, setID, manifestID)
		if err != nil {
			return err
		}
		zap.S().Debugf("relation deleted between device %q and manifest %q", setID, manifestID)
		return nil
	}); err != nil {
		return err
	}

	oldDeviceSelectors := substract1(oldManifest.DeviceIDs, m.DeviceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, oldDeviceSelectors, m.Id, func(ctx context.Context, deviceID, manifestID string) error {
		if _, err := w.pgDeviceRepo.GetDevice(ctx, deviceID); err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return nil
			}
		}
		err := w.pgReferenceRepo.DeleteDeviceRelation(ctx, deviceID, manifestID)
		if err != nil {
			return err
		}
		zap.S().Debugf("relation deleted between device %q and manifest %q", deviceID, manifestID)
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *Service) CreateRelations(ctx context.Context, m entity.ManifestWork) error {
	for _, s := range m.Selectors {
		switch s.Type {
		case entity.NamespaceSelector:
			namespace, err := w.pgDeviceRepo.GetNamespace(ctx, s.Value)
			if err != nil {
				if errors.Is(err, common.ErrResourceNotFound) {
					zap.S().Warnw("unable to create relation. namespace does not exist", "namespace", s.Value)
					continue
				}
				return fmt.Errorf("unable to get namespace %q: %w", s.Value, err)
			}
			if contains(namespace.ManifestIDS, m.Id) {
				continue
			}
			if err := w.pgReferenceRepo.CreateNamespaceRelation(ctx, namespace.Name, m.Id); err != nil {
				return fmt.Errorf("unable to create namespace %q manifest %q relation: %w", namespace.Name, m.Id, err)
			}
		case entity.SetSelector:
			set, err := w.pgDeviceRepo.GetSet(ctx, s.Value)
			if err != nil {
				if errors.Is(err, common.ErrResourceNotFound) {
					zap.S().Warnw("unable to create relation. set does not exist", "set", s.Value)
					continue
				}
				return fmt.Errorf("unable to get set %q: %w", s.Value, err)
			}
			if contains(set.ManifestIDS, m.Id) {
				continue
			}
			if err := w.pgReferenceRepo.CreateSetRelation(ctx, set.Name, m.Id); err != nil {
				return fmt.Errorf("unable to create set %q manifest %q relation: %w", set.Name, m.Id, err)
			}
		case entity.DeviceSelector:
			device, err := w.pgDeviceRepo.GetDevice(ctx, s.Value)
			if err != nil {
				if errors.Is(err, common.ErrResourceNotFound) {
					zap.S().Warnw("unable to create relation. device does not exist", "device_id", s.Value)
					continue
				}
				return fmt.Errorf("unable to get device %q: %w", s.Value, err)
			}
			if err := w.pgReferenceRepo.CreateDeviceRelation(ctx, device.ID, m.Id); err != nil {
				return fmt.Errorf("unable to create device %q manifest %q relation: %w", device.ID, m.Id, err)
			}

		}
	}
	return nil
}

func (w *Service) deleteManifests(ctx context.Context, manifests []entity.ManifestReference) {
	for _, m := range manifests {
		if err := w.pgReferenceRepo.DeleteReference(ctx, m); err != nil {
			zap.S().Error("unable to delete manifest", "error", err, "manifest_id", m.Id, "manifest_repo", m.Repo.LocalPath)
			continue
		}
	}
}

func (w *Service) updateManifests(ctx context.Context, manifests []entity.ManifestReference) {
	for _, m := range manifests {
		if err := w.pgReferenceRepo.UpdateReference(ctx, m); err != nil {
			zap.S().Errorw("unable to update manifest", "error", err, "manifest_id", m.Id, "manifest_repo", m.Repo.LocalPath)
			continue
		}
	}
}

func (w *Service) createReference(manifest entity.ManifestWork, repo entity.Repository) entity.ManifestReference {
	ref := entity.ManifestReference{
		Id:           manifest.Id,
		Hash:         manifest.Hash,
		Path:         manifest.Path,
		Valid:        manifest.Valid,
		Repo:         repo,
		DeviceIDs:    make([]string, 0, len(manifest.Selectors)),
		SetIDs:       make([]string, 0, len(manifest.Selectors)),
		NamespaceIDs: make([]string, 0, len(manifest.Selectors)),
	}

	for _, s := range manifest.Selectors {
		switch s.Type {
		case entity.NamespaceSelector:
			ref.NamespaceIDs = append(ref.NamespaceIDs, s.Value)
		case entity.SetSelector:
			ref.SetIDs = append(ref.SetIDs, s.Value)
		case entity.DeviceSelector:
			ref.DeviceIDs = append(ref.DeviceIDs, s.Value)
		}
	}

	return ref
}
