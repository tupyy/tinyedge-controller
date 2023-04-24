package manifest

import (
	"context"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
	"go.uber.org/zap"
)

type Service struct {
	manifestReaderWriter ManifestReaderWriter
	deviceReader         DeviceReader
	gitReader            GitReader
}

func New(deviceReader DeviceReader, rw ManifestReaderWriter, git GitReader) *Service {
	return &Service{
		deviceReader:         deviceReader,
		gitReader:            git,
		manifestReaderWriter: rw,
	}
}

func (w *Service) GetManifests(ctx context.Context, repo entity.Repository) ([]entity.Manifest, error) {
	return nil, nil
}

func (w *Service) GetManifest(ctx context.Context, id string) (entity.Manifest, error) {
	return nil, nil
}

func (w *Service) UpdateManifests(ctx context.Context, repo entity.Repository) error {
	pgManifests, err := w.manifestReaderWriter.GetManifests(ctx, repo)
	if err != nil {
		return fmt.Errorf("unable to read manifests of repo %q: %w", repo.Id, err)
	}

	gitManifests, err := w.gitReader.GetManifests(ctx, repo)
	if err != nil {
		return fmt.Errorf("unable to read manifest from repo %q: %w", repo.Id, err)
	}

	created := substract(gitManifests, pgManifests, func(m entity.Manifest) string { return m.GetID() })
	deleted := substract(pgManifests, gitManifests, func(m entity.Manifest) string { return m.GetID() })
	updated := intersect(gitManifests, pgManifests, func(m entity.Manifest) string { return m.GetID() }, func(m1, m2 entity.Manifest) bool { return m1.GetHash() == m1.GetHash() })

	for _, c := range created {
		if err := w.manifestReaderWriter.InsertManifest(ctx, c); err != nil && !errService.IsResourceAlreadyExists(err) {
			return fmt.Errorf("unable to insert manifest %q: %w", c.GetID(), err)
		}
		if err := w.UpdateRelations(ctx, c); err != nil {
			return err
		}
	}

	for _, d := range deleted {
		if err := w.manifestReaderWriter.DeleteManifest(ctx, d.GetID()); err != nil {
			return fmt.Errorf("unable to delete manifest %q: %w", d.GetID(), err)
		}
	}

	for _, u := range updated {
		if err := w.UpdateRelations(ctx, u); err != nil {
			return err
		}
		if err := w.manifestReaderWriter.UpdateManifest(ctx, u); err != nil {
			return fmt.Errorf("unable to update manifest %q: %w", u.GetID(), err)
		}
	}

	return nil
}

func (w *Service) UpdateRelations(ctx context.Context, gitManifest entity.Manifest) error {
	// get the old pgManifest
	pgManifest, err := w.manifestReaderWriter.GetManifest(ctx, gitManifest.GetID())
	if err != nil {
		return err
	}

	// updateRelation updates the relation between a selector value and the manifest
	updateRelation := func(ctx context.Context, resources []string, m string, fn func(ctx context.Context, resourceID, manifestID string) error) error {
		if len(resources) > 0 {
			for _, resource := range resources {
				if err := fn(ctx, resource, m); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// for each new relation needs to be created (exists in m but not in oldManifest)
	namespaces := substract(gitManifest.GetSelectors().ExtractType(entity.NamespaceSelector), pgManifest.GetNamespaces(), func(i string) string { return i })
	if err := updateRelation(ctx, namespaces, gitManifest.GetID(), func(ctx context.Context, namespaceID, manifestID string) error {
		if _, err := w.deviceReader.GetNamespace(ctx, namespaceID); err != nil {
			if errService.IsResourceNotFound(err) {
				return nil
			}
			return fmt.Errorf("unable to create relation between namespace %q and manifest %q: %w", namespaceID, manifestID, err)
		}
		if err := w.manifestReaderWriter.CreateRelation(ctx, entity.NewNamespaceRelation(namespaceID, manifestID)); err != nil {
			if !errService.IsResourceAlreadyExists(err) {
				return fmt.Errorf("unable to create relation between namespace %q and manifest %q: %w", namespaceID, manifestID, err)
			}
		}
		zap.S().Debugf("relation created between namespace %q and manifest %q", namespaceID, manifestID)
		return nil
	}); err != nil {
		return err
	}

	sets := substract(gitManifest.GetSelectors().ExtractType(entity.SetSelector), pgManifest.GetSets(), func(i string) string { return i })
	if err := updateRelation(ctx, sets, gitManifest.GetID(), func(ctx context.Context, setID, manifestID string) error {
		if _, err := w.deviceReader.GetSet(ctx, setID); err != nil {
			if errService.IsResourceNotFound(err) {
				return nil
			}
			return fmt.Errorf("unable to create relation between set %q and manifest %q: %w", setID, manifestID, err)
		}
		if err := w.manifestReaderWriter.CreateRelation(ctx, entity.NewSetRelation(setID, manifestID)); err != nil {
			if !errService.IsResourceAlreadyExists(err) {
				return fmt.Errorf("unable to create relation between set %q and manifest %q: %w", setID, manifestID, err)
			}
		}
		zap.S().Debugf("relation created between set %q and manifest %q", setID, manifestID)
		return nil
	}); err != nil {
		return err
	}

	devices := substract(gitManifest.GetSelectors().ExtractType(entity.DeviceSelector), pgManifest.GetDevices(), func(i string) string { return i })
	if err := updateRelation(ctx, devices, gitManifest.GetID(), func(ctx context.Context, deviceID, manifestID string) error {
		if _, err := w.deviceReader.GetDevice(ctx, deviceID); err != nil {
			if errService.IsResourceNotFound(err) {
				return nil
			}
			return fmt.Errorf("unable to create relation between device %q and manifest %q: %w", deviceID, manifestID, err)
		}
		if err := w.manifestReaderWriter.CreateRelation(ctx, entity.NewDeviceRelation(deviceID, manifestID)); err != nil {
			if !errService.IsResourceAlreadyExists(err) {
				return fmt.Errorf("unable to create relation between device %q and manifest %q: %w", deviceID, manifestID, err)
			}
		}
		zap.S().Debugf("relation created between device %q and manifest %q", deviceID, manifestID)
		return nil
	}); err != nil {
		return err
	}

	// remove the old ones
	namespaces = substract(pgManifest.GetNamespaces(), gitManifest.GetSelectors().ExtractType(entity.NamespaceSelector), func(i string) string { return i })
	if err := updateRelation(ctx, namespaces, gitManifest.GetID(), func(ctx context.Context, namespaceID, manifestID string) error {
		if _, err := w.deviceReader.GetNamespace(ctx, namespaceID); err != nil {
			if errService.IsResourceNotFound(err) {
				return nil
			}
		}
		err := w.manifestReaderWriter.DeleteRelation(ctx, entity.NewNamespaceRelation(namespaceID, manifestID))
		if err != nil {
			return fmt.Errorf("unable to delete  between namespace %q and manifest %q: %w", namespaceID, manifestID, err)
		}
		zap.S().Debugf("relation deleted between device %q and manifest %q", namespaceID, manifestID)
		return nil
	}); err != nil {
		return err
	}

	sets = substract(pgManifest.GetSelectors().ExtractType(entity.SetSelector), gitManifest.GetSelectors().ExtractType(entity.SetSelector), func(i string) string { return i })
	if err := updateRelation(ctx, sets, gitManifest.GetID(), func(ctx context.Context, setID, manifestID string) error {
		if _, err := w.deviceReader.GetSet(ctx, setID); err != nil {
			if errService.IsResourceNotFound(err) {
				return nil
			}
		}
		err := w.manifestReaderWriter.DeleteRelation(ctx, entity.NewSetRelation(setID, manifestID))
		if err != nil {
			return fmt.Errorf("unable to delete  between set %q and manifest %q: %w", setID, manifestID, err)
		}
		zap.S().Debugf("relation deleted between device %q and manifest %q", setID, manifestID)
		return nil
	}); err != nil {
		return err
	}

	devices = substract(pgManifest.GetDevices(), gitManifest.GetSelectors().ExtractType(entity.DeviceSelector), func(i string) string { return i })
	if err := updateRelation(ctx, devices, gitManifest.GetID(), func(ctx context.Context, deviceID, manifestID string) error {
		if _, err := w.deviceReader.GetDevice(ctx, deviceID); err != nil {
			if errService.IsResourceNotFound(err) {
				return nil
			}
		}
		err := w.manifestReaderWriter.DeleteRelation(ctx, entity.NewDeviceRelation(deviceID, manifestID))
		if err != nil {
			return fmt.Errorf("unable to delete  between device %q and manifest %q: %w", deviceID, manifestID, err)
		}
		zap.S().Debugf("relation deleted between device %q and manifest %q", deviceID, manifestID)
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *Service) CreateRelations(ctx context.Context, m entity.Manifest) error {
	for _, s := range m.GetSelectors() {
		var r entity.Relation
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
			r = entity.NewNamespaceRelation(namespace.Name, m.GetID())
		case entity.SetSelector:
			set, err := w.deviceReader.GetSet(ctx, s.Value)
			if err != nil {
				if errService.IsResourceNotFound(err) {
					zap.S().Warnw("unable to create relation. set does not exist", "set", s.Value)
					continue
				}
				return fmt.Errorf("unable to get set %q: %w", s.Value, err)
			}
			r = entity.NewSetRelation(set.Name, m.GetID())
		case entity.DeviceSelector:
			device, err := w.deviceReader.GetDevice(ctx, s.Value)
			if err != nil {
				if errService.IsResourceNotFound(err) {
					zap.S().Warnw("unable to create relation. device does not exist", "device_id", s.Value)
					continue
				}
				return fmt.Errorf("unable to get device %q: %w", s.Value, err)
			}
			r = entity.NewDeviceRelation(device.ID, m.GetID())
		}

		if err := w.manifestReaderWriter.CreateRelation(ctx, r); err != nil {
			if errService.IsResourceAlreadyExists(err) {
				zap.S().Info("Relation %q between %s and manifest %s already exists", r.Type, r.ResourceID, r.ManifestID)
				continue
			}
			return fmt.Errorf("unable to create relation between resource %q and manifest %q: %w", r.ResourceID, r.ResourceID, err)
		}
	}
	return nil
}

func (w *Service) deleteManifests(ctx context.Context, manifests []entity.Manifest) {
	for _, m := range manifests {
		if err := w.manifestReaderWriter.DeleteManifest(ctx, m.GetID()); err != nil {
			zap.S().Error("unable to delete manifest", "error", err, "manifest_id", m.GetID())
			continue
		}
	}
}

func (w *Service) updateManifests(ctx context.Context, manifests []entity.Manifest) {
	for _, m := range manifests {
		if err := w.manifestReaderWriter.UpdateManifest(ctx, m); err != nil {
			zap.S().Errorw("unable to update manifest", "error", err, "manifest_id", m.GetID())
			continue
		}
	}
}
