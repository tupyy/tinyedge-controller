package reference

import (
	"context"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
	"go.uber.org/zap"
)

type Service struct {
	refReaderWriter ReferenceReaderWriter
	deviceReader    DeviceReader
	gitReader       GitReader
}

func New(deviceReader DeviceReader, ref ReferenceReaderWriter, git GitReader) *Service {
	return &Service{
		deviceReader:    deviceReader,
		refReaderWriter: ref,
		gitReader:       git,
	}
}

func (w *Service) GetReference(ctx context.Context, id string) (entity.ManifestReference, error) {
	return w.refReaderWriter.GetReference(ctx, id)
}

func (w *Service) GetReferences(ctx context.Context, repo entity.Repository) ([]entity.ManifestReference, error) {
	return w.refReaderWriter.GetReferences(ctx, repo)
}

func (w *Service) UpdateReferences(ctx context.Context, repo entity.Repository) error {
	references, err := w.refReaderWriter.GetReferences(ctx, repo)
	if err != nil {
		return fmt.Errorf("unable to read references of repo %q: %w", repo.Id, err)
	}

	manifests, err := w.gitReader.GetReferences(ctx, repo)
	if err != nil {
		return fmt.Errorf("unable to read manifest from repo %q: %w", repo.Id, err)
	}

	created := substract(manifests, references, func(item entity.ManifestReference) string { return item.Id })
	deleted := substract(references, manifests, func(item entity.ManifestReference) string { return item.Id })
	updated := intersect(manifests, references,
		func(item entity.ManifestReference) string { return item.Id },
		func(manifest entity.ManifestReference, ref entity.ManifestReference) bool {
			return ref.Hash != manifest.Hash
		},
	)

	for _, c := range created {
		if err := w.refReaderWriter.InsertReference(ctx, c); err != nil && !errService.IsResourceAlreadyExists(err) {
			return fmt.Errorf("unable to insert reference %q: %w", c.Id, err)
		}
		if err := w.UpdateRelations(ctx, c); err != nil {
			return err
		}
	}

	for _, d := range deleted {
		if err := w.refReaderWriter.DeleteReference(ctx, d); err != nil {
			return fmt.Errorf("unable to delete reference %q: %w", d.Id, err)
		}
	}

	for _, u := range updated {
		if err := w.refReaderWriter.UpdateReference(ctx, u); err != nil {
			return fmt.Errorf("unable to update reference %q: %w", u.Id, err)
		}
		if err := w.UpdateRelations(ctx, u); err != nil {
			return err
		}
	}

	return nil
}

func (w *Service) UpdateRelations(ctx context.Context, m entity.ManifestReference) error {
	// get the old manifest
	oldManifest, err := w.refReaderWriter.GetReference(ctx, m.Id)
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
	newNamespaceSelectors := substract(m.NamespaceIDs, oldManifest.NamespaceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, newNamespaceSelectors, m.Id, func(ctx context.Context, namespaceID, referenceID string) error {
		if _, err := w.deviceReader.GetNamespace(ctx, namespaceID); err != nil {
			if errService.IsResourceNotFound(err) {
				return nil
			}
			return fmt.Errorf("unable to create relation between namespace %q and reference %q: %w", namespaceID, referenceID, err)
		}
		if err := w.refReaderWriter.CreateRelation(ctx, entity.NewNamespaceRelation(namespaceID, referenceID)); err != nil {
			if !errService.IsResourceAlreadyExists(err) {
				return fmt.Errorf("unable to create relation between namespace %q and reference %q: %w", namespaceID, referenceID, err)
			}
		}
		zap.S().Debugf("relation created between namespace %q and reference %q", namespaceID, referenceID)
		return nil
	}); err != nil {
		return err
	}

	newSetSelectors := substract(m.SetIDs, oldManifest.SetIDs, func(i string) string { return i })
	if err := updateRelation(ctx, newSetSelectors, m.Id, func(ctx context.Context, setID, referenceID string) error {
		if _, err := w.deviceReader.GetSet(ctx, setID); err != nil {
			if errService.IsResourceNotFound(err) {
				return nil
			}
			return fmt.Errorf("unable to create relation between set %q and reference %q: %w", setID, referenceID, err)
		}
		if err := w.refReaderWriter.CreateRelation(ctx, entity.NewSetRelation(setID, referenceID)); err != nil {
			if !errService.IsResourceAlreadyExists(err) {
				return fmt.Errorf("unable to create relation between set %q and reference %q: %w", setID, referenceID, err)
			}
		}
		zap.S().Debugf("relation created between set %q and reference %q", setID, referenceID)
		return nil
	}); err != nil {
		return err
	}

	newDeviceSelectors := substract(m.DeviceIDs, oldManifest.DeviceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, newDeviceSelectors, m.Id, func(ctx context.Context, deviceID, referenceID string) error {
		if _, err := w.deviceReader.GetDevice(ctx, deviceID); err != nil {
			if errService.IsResourceNotFound(err) {
				return nil
			}
			return fmt.Errorf("unable to create relation between device %q and reference %q: %w", deviceID, referenceID, err)
		}
		if err := w.refReaderWriter.CreateRelation(ctx, entity.NewDeviceRelation(deviceID, referenceID)); err != nil {
			if !errService.IsResourceAlreadyExists(err) {
				return fmt.Errorf("unable to create relation between device %q and reference %q: %w", deviceID, referenceID, err)
			}
		}
		zap.S().Debugf("relation created between device %q and reference %q", deviceID, referenceID)
		return nil
	}); err != nil {
		return err
	}

	// remove the old ones
	oldNamespaceSelectors := substract(oldManifest.NamespaceIDs, m.NamespaceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, oldNamespaceSelectors, m.Id, func(ctx context.Context, namespaceID, referenceID string) error {
		if _, err := w.deviceReader.GetNamespace(ctx, namespaceID); err != nil {
			if errService.IsResourceNotFound(err) {
				return nil
			}
		}
		err := w.refReaderWriter.DeleteRelation(ctx, entity.NewNamespaceRelation(namespaceID, referenceID))
		if err != nil {
			return fmt.Errorf("unable to delete  between namespace %q and reference %q: %w", namespaceID, referenceID, err)
		}
		zap.S().Debugf("relation deleted between device %q and reference %q", namespaceID, referenceID)
		return nil
	}); err != nil {
		return err
	}

	oldSetSelectors := substract(oldManifest.SetIDs, m.SetIDs, func(i string) string { return i })
	if err := updateRelation(ctx, oldSetSelectors, m.Id, func(ctx context.Context, setID, referenceID string) error {
		if _, err := w.deviceReader.GetSet(ctx, setID); err != nil {
			if errService.IsResourceNotFound(err) {
				return nil
			}
		}
		err := w.refReaderWriter.DeleteRelation(ctx, entity.NewSetRelation(setID, referenceID))
		if err != nil {
			return fmt.Errorf("unable to delete  between set %q and reference %q: %w", setID, referenceID, err)
		}
		zap.S().Debugf("relation deleted between device %q and reference %q", setID, referenceID)
		return nil
	}); err != nil {
		return err
	}

	oldDeviceSelectors := substract(oldManifest.DeviceIDs, m.DeviceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, oldDeviceSelectors, m.Id, func(ctx context.Context, deviceID, referenceID string) error {
		if _, err := w.deviceReader.GetDevice(ctx, deviceID); err != nil {
			if errService.IsResourceNotFound(err) {
				return nil
			}
		}
		err := w.refReaderWriter.DeleteRelation(ctx, entity.NewDeviceRelation(deviceID, referenceID))
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
			namespace, err := w.deviceReader.GetNamespace(ctx, s.Value)
			if err != nil {
				if errService.IsResourceNotFound(err) {
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
			set, err := w.deviceReader.GetSet(ctx, s.Value)
			if err != nil {
				if errService.IsResourceNotFound(err) {
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
			device, err := w.deviceReader.GetDevice(ctx, s.Value)
			if err != nil {
				if errService.IsResourceNotFound(err) {
					zap.S().Warnw("unable to create relation. device does not exist", "device_id", s.Value)
					continue
				}
				return fmt.Errorf("unable to get device %q: %w", s.Value, err)
			}
			r = entity.NewDeviceRelation(device.ID, m.Id)
		}
		if err := w.refReaderWriter.CreateRelation(ctx, r); err != nil {
			return fmt.Errorf("unable to create relation between resource %q and manifest %q: %w", r.ResourceID, r.ResourceID, err)
		}
	}
	return nil
}

func (w *Service) deleteManifests(ctx context.Context, manifests []entity.ManifestReference) {
	for _, m := range manifests {
		if err := w.refReaderWriter.DeleteReference(ctx, m); err != nil {
			zap.S().Error("unable to delete manifest", "error", err, "manifest_id", m.Id, "manifest_repo", m.Repo.LocalPath)
			continue
		}
	}
}

func (w *Service) updateManifests(ctx context.Context, manifests []entity.ManifestReference) {
	for _, m := range manifests {
		if err := w.refReaderWriter.UpdateReference(ctx, m); err != nil {
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
