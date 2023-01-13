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
	secretReader    common.SecretReader
}

func New(pgDeviceRepo common.DeviceReader, pgManifestRepo common.ReferenceReaderWriter, gitRepo common.GitReader, secretReader common.SecretReader) *Service {
	return &Service{
		pgDeviceRepo:    pgDeviceRepo,
		pgReferenceRepo: pgManifestRepo,
		gitRepo:         gitRepo,
		secretReader:    secretReader,
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

	// for each secret in the manifest get the value from vault
	for i := 0; i < len(manifest.Secrets); i++ {
		secret := &manifest.Secrets[i]
		s, err := w.secretReader.GetSecret(ctx, secret.Path, secret.Key)
		if err != nil {
			return entity.ManifestWork{}, fmt.Errorf("unable to read secret %q from vault: %w", secret.Path, err)
		}
		secret.Value = s.Value
		secret.Hash = s.Hash
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
			return fmt.Errorf("unable to insert reference %q: %w", c.Id, err)
		}
		if err := w.UpdateRelations(ctx, ref); err != nil {
			return err
		}
	}

	for _, d := range deleted {
		if err := w.pgReferenceRepo.DeleteReference(ctx, d); err != nil {
			return fmt.Errorf("unable to delete reference %q: %w", d.Id, err)
		}
	}

	for _, u := range updated {
		if err := w.pgReferenceRepo.UpdateReference(ctx, w.createReference(u, repo)); err != nil {
			return fmt.Errorf("unable to update reference %q: %w", u.Id, err)
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
	updateRelation := func(ctx context.Context, n []string, m string, fn func(ctx context.Context, resourceID, referenceID string) error) error {
		if len(n) > 0 {
			for _, s := range n {
				if err := fn(ctx, s, m); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// for each new relation needs to be created (exists in m but not in oldManifest)
	newNamespaceSelectors := substract1(m.NamespaceIDs, oldManifest.NamespaceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, newNamespaceSelectors, m.Id, func(ctx context.Context, namespaceID, referenceID string) error {
		if _, err := w.pgDeviceRepo.GetNamespace(ctx, namespaceID); err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return nil
			}
			return fmt.Errorf("unable to create relation between namespace %q and reference %q: %w", namespaceID, referenceID, err)
		}
		if err := w.pgReferenceRepo.CreateRelation(ctx, entity.NewDeviceRelation(namespaceID, referenceID)); err != nil {
			if !errors.Is(err, common.ErrResourceAlreadyExists) {
				return fmt.Errorf("unable to create relation between namespace %q and reference %q: %w", namespaceID, referenceID, err)
			}
		}
		zap.S().Debugf("relation created between namespace %q and reference %q", namespaceID, referenceID)
		return nil
	}); err != nil {
		return err
	}

	newSetSelectors := substract1(m.SetIDs, oldManifest.SetIDs, func(i string) string { return i })
	if err := updateRelation(ctx, newSetSelectors, m.Id, func(ctx context.Context, setID, referenceID string) error {
		if _, err := w.pgDeviceRepo.GetSet(ctx, setID); err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return nil
			}
			return fmt.Errorf("unable to create relation between set %q and reference %q: %w", setID, referenceID, err)
		}
		if err := w.pgReferenceRepo.CreateRelation(ctx, entity.NewSetRelation(setID, referenceID)); err != nil {
			if !errors.Is(err, common.ErrResourceAlreadyExists) {
				return fmt.Errorf("unable to create relation between set %q and reference %q: %w", setID, referenceID, err)
			}
		}
		zap.S().Debugf("relation created between set %q and reference %q", setID, referenceID)
		return nil
	}); err != nil {
		return err
	}

	newDeviceSelectors := substract1(m.DeviceIDs, oldManifest.DeviceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, newDeviceSelectors, m.Id, func(ctx context.Context, deviceID, referenceID string) error {
		if _, err := w.pgDeviceRepo.GetDevice(ctx, deviceID); err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return nil
			}
			return fmt.Errorf("unable to create relation between device %q and reference %q: %w", deviceID, referenceID, err)
		}
		if err := w.pgReferenceRepo.CreateRelation(ctx, entity.NewDeviceRelation(deviceID, referenceID)); err != nil {
			if !errors.Is(err, common.ErrResourceAlreadyExists) {
				return fmt.Errorf("unable to create relation between device %q and reference %q: %w", deviceID, referenceID, err)
			}
		}
		zap.S().Debugf("relation created between device %q and reference %q", deviceID, referenceID)
		return nil
	}); err != nil {
		return err
	}

	// remove the old ones
	oldNamespaceSelectors := substract1(oldManifest.NamespaceIDs, m.NamespaceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, oldNamespaceSelectors, m.Id, func(ctx context.Context, namespaceID, referenceID string) error {
		if _, err := w.pgDeviceRepo.GetNamespace(ctx, namespaceID); err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return nil
			}
		}
		err := w.pgReferenceRepo.DeleteRelation(ctx, entity.NewNamespaceRelation(namespaceID, referenceID))
		if err != nil {
			return fmt.Errorf("unable to delete  between namespace %q and reference %q: %w", namespaceID, referenceID, err)
		}
		zap.S().Debugf("relation deleted between device %q and reference %q", namespaceID, referenceID)
		return nil
	}); err != nil {
		return err
	}

	oldSetSelectors := substract1(oldManifest.SetIDs, m.SetIDs, func(i string) string { return i })
	if err := updateRelation(ctx, oldSetSelectors, m.Id, func(ctx context.Context, setID, referenceID string) error {
		if _, err := w.pgDeviceRepo.GetSet(ctx, setID); err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return nil
			}
		}
		err := w.pgReferenceRepo.DeleteRelation(ctx, entity.NewSetRelation(setID, referenceID))
		if err != nil {
			return fmt.Errorf("unable to delete  between set %q and reference %q: %w", setID, referenceID, err)
		}
		zap.S().Debugf("relation deleted between device %q and reference %q", setID, referenceID)
		return nil
	}); err != nil {
		return err
	}

	oldDeviceSelectors := substract1(oldManifest.DeviceIDs, m.DeviceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, oldDeviceSelectors, m.Id, func(ctx context.Context, deviceID, referenceID string) error {
		if _, err := w.pgDeviceRepo.GetDevice(ctx, deviceID); err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return nil
			}
		}
		err := w.pgReferenceRepo.DeleteRelation(ctx, entity.NewDeviceRelation(deviceID, referenceID))
		if err != nil {
			return fmt.Errorf("unable to delete  between device %q and reference %q: %w", deviceID, referenceID, err)
		}
		zap.S().Debugf("relation deleted between device %q and reference %q", deviceID, referenceID)
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *Service) CreateRelations(ctx context.Context, m entity.ManifestWork) error {
	for _, s := range m.Selectors {
		var r entity.ReferenceRelation
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
			r = entity.NewNamespaceRelation(namespace.Name, m.Id)
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
			r = entity.NewSetRelation(set.Name, m.Id)
		case entity.DeviceSelector:
			device, err := w.pgDeviceRepo.GetDevice(ctx, s.Value)
			if err != nil {
				if errors.Is(err, common.ErrResourceNotFound) {
					zap.S().Warnw("unable to create relation. device does not exist", "device_id", s.Value)
					continue
				}
				return fmt.Errorf("unable to get device %q: %w", s.Value, err)
			}
			r = entity.NewDeviceRelation(device.ID, m.Id)
		}
		if err := w.pgReferenceRepo.CreateRelation(ctx, r); err != nil {
			return fmt.Errorf("unable to create relation between resource %q and manifest %q: %w", r.ResourceID, r.ResourceID, err)
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
