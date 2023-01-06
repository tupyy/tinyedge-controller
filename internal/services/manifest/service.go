package manifest

import (
	"context"
	"errors"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
	"go.uber.org/zap"
)

type relation struct {
	RelationType int
	ManifestID   string
	ObjectID     string
}

type Service struct {
	pgManifestRepo common.ManifestReaderWriter
	pgDeviceRepo   common.DeviceReader
}

func New(pgDeviceRepo common.DeviceReader, pgManifestRepo common.ManifestReaderWriter) *Service {
	return &Service{
		pgDeviceRepo:   pgDeviceRepo,
		pgManifestRepo: pgManifestRepo,
	}
}

func (w *Service) GetRepositories(ctx context.Context) ([]entity.Repository, error) {
	repos, err := w.pgManifestRepo.GetRepositories(ctx)
	if err != nil {
		return []entity.Repository{}, err
	}
	return repos, nil

}

func (w *Service) GetManifests(ctx context.Context, repo entity.Repository) ([]entity.ManifestWork, error) {
	return w.pgManifestRepo.GetRepoManifests(ctx, repo)
}

func (w *Service) UpdateRepository(ctx context.Context, r entity.Repository) error {
	if err := w.pgManifestRepo.UpdateRepo(ctx, r); err != nil {
		return err
	}

	return nil
}

func (w *Service) UpdateManifests(ctx context.Context, repo entity.Repository, gitManifests []entity.ManifestWork) (created []entity.ManifestWork, updated []entity.ManifestWork, deleted []entity.ManifestWork, err error) {
	created = make([]entity.ManifestWork, 0)
	deleted = make([]entity.ManifestWork, 0)
	updated = make([]entity.ManifestWork, 0)

	pgManifests, err := w.pgManifestRepo.GetRepoManifests(ctx, repo)
	if err != nil {
		return created, updated, deleted, err
	}

	created = substract(gitManifests, pgManifests, func(item entity.ManifestWork) string { return item.Id })
	deleted = substract(pgManifests, gitManifests, func(item entity.ManifestWork) string { return item.Id })
	gitUpdated := intersect(gitManifests, pgManifests, func(item entity.ManifestWork) string { return item.Id })
	pgUpdated := intersect(pgManifests, gitManifests, func(item entity.ManifestWork) string { return item.Id })

	// look for manifests which have been updated between git and pg.
	for _, m := range gitUpdated {
		for _, pm := range pgUpdated {
			if pm.Path == m.Path && pm.Hash != m.Hash {
				updated = append(updated, m)
				break
			}
		}
	}

	w.insertManifests(ctx, created)
	w.deleteManifests(ctx, deleted)
	w.updateManifests(ctx, updated)

	return created, updated, deleted, nil
}

func (w *Service) UpdateRelations(ctx context.Context, m entity.ManifestWork) error {
	// get the old manifest
	oldManifest, err := w.pgManifestRepo.GetManifest(ctx, m.Id)
	if err != nil {
		return err
	}

	// updateRelation updates the relation between a selector value and the manifest
	updateRelation := func(ctx context.Context, n []string, m string, fn func(ctx context.Context, n, m string) error, getFn func(ctx context.Context, id string) error) error {
		if len(n) > 0 {
			for _, s := range n {
				err := getFn(ctx, s)
				if err != nil && errors.Is(err, common.ErrResourceNotFound) {
					continue
				}
				if err := fn(ctx, s, m); err != nil {
					return fmt.Errorf("unable to create/delete relation between selector %q and manifest %q: %w", s, m, err)
				}
			}
		}
		return nil
	}

	// getIds returns a list of selector values based on selector type
	getIds := func(selectors []entity.Selector, selectorType entity.SelectorType) []string {
		list := make([]string, 0, len(selectors))
		for _, s := range selectors {
			if s.Type == selectorType {
				list = append(list, s.Value)
			}
		}
		return list
	}

	// for each new relation needs to be created (exists in m but not in oldManifest)
	newNamespaceSelectors := substract(getIds(m.Selectors, entity.NamespaceSelector), oldManifest.NamespaceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, newNamespaceSelectors, m.Id, w.pgManifestRepo.CreateNamespaceRelation, func(ctx context.Context, id string) error {
		_, err := w.pgDeviceRepo.GetNamespace(ctx, id)
		return err
	}); err != nil {
		return err
	}

	newSetSelectors := substract(getIds(m.Selectors, entity.SetSelector), oldManifest.SetIDs, func(i string) string { return i })
	if err := updateRelation(ctx, newSetSelectors, m.Id, w.pgManifestRepo.CreateSetRelation, func(ctx context.Context, id string) error {
		_, err := w.pgDeviceRepo.GetSet(ctx, id)
		return err
	}); err != nil {
		return err
	}

	newDeviceSelectors := substract(getIds(m.Selectors, entity.DeviceSelector), oldManifest.DeviceIDs, func(i string) string { return i })
	if err := updateRelation(ctx, newDeviceSelectors, m.Id, w.pgManifestRepo.CreateDeviceRelation, func(ctx context.Context, id string) error {
		_, err := w.pgDeviceRepo.GetDevice(ctx, id)
		return err
	}); err != nil {
		return err
	}

	// remove the old ones
	oldNamespaceSelectors := substract(oldManifest.NamespaceIDs, getIds(m.Selectors, entity.NamespaceSelector), func(i string) string { return i })
	if err := updateRelation(ctx, oldNamespaceSelectors, m.Id, w.pgManifestRepo.DeleteNamespaceRelation, func(ctx context.Context, id string) error {
		_, err := w.pgDeviceRepo.GetNamespace(ctx, id)
		return err
	}); err != nil {
		return err
	}

	oldSetSelectors := substract(oldManifest.SetIDs, getIds(m.Selectors, entity.SetSelector), func(i string) string { return i })
	if err := updateRelation(ctx, oldSetSelectors, m.Id, w.pgManifestRepo.DeleteSetRelation, func(ctx context.Context, id string) error {
		_, err := w.pgDeviceRepo.GetSet(ctx, id)
		return err
	}); err != nil {
		return err
	}

	oldDeviceSelectors := substract(oldManifest.DeviceIDs, getIds(m.Selectors, entity.DeviceSelector), func(i string) string { return i })
	if err := updateRelation(ctx, oldDeviceSelectors, m.Id, w.pgManifestRepo.DeleteDeviceRelation, func(ctx context.Context, id string) error {
		_, err := w.pgDeviceRepo.GetDevice(ctx, id)
		return err
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
			if err := w.pgManifestRepo.CreateNamespaceRelation(ctx, namespace.Name, m.Id); err != nil {
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
			if err := w.pgManifestRepo.CreateSetRelation(ctx, set.Name, m.Id); err != nil {
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
			if err := w.pgManifestRepo.CreateDeviceRelation(ctx, device.ID, m.Id); err != nil {
				return fmt.Errorf("unable to create device %q manifest %q relation: %w", device.ID, m.Id, err)
			}

		}
	}
	return nil
}

func (w *Service) insertManifests(ctx context.Context, manifests []entity.ManifestWork) {
	for _, m := range manifests {
		if err := w.pgManifestRepo.InsertManifest(ctx, m); err != nil {
			zap.S().Errorw("unable to insert manifest", "error", err, "manifest_id", m.Id, "manifest_repo", m.Repo.LocalPath)
			continue
		}
	}
}

func (w *Service) deleteManifests(ctx context.Context, manifests []entity.ManifestWork) {
	for _, m := range manifests {
		if err := w.pgManifestRepo.DeleteManifest(ctx, m); err != nil {
			zap.S().Error("unable to delete manifest", "error", err, "manifest_id", m.Id, "manifest_repo", m.Repo.LocalPath)
			continue
		}
	}
}

func (w *Service) updateManifests(ctx context.Context, manifests []entity.ManifestWork) {
	for _, m := range manifests {
		if err := w.pgManifestRepo.UpdateManifest(ctx, m); err != nil {
			zap.S().Errorw("unable to update manifest", "error", err, "manifest_id", m.Id, "manifest_repo", m.Repo.LocalPath)
			continue
		}
	}
}

// substract return all elements of a which are not found in b
func substract[T any, slice []T, S func(elem T) string](a slice, b slice, fn S) slice {
	m1 := make(map[string]T)
	m2 := make(map[string]T)

	limit := len(a)
	if limit < len(b) {
		limit = len(b)
	}

	for i := 0; i < limit; i++ {
		if i < len(a) {
			id := fn(a[i])
			m1[id] = a[i]
		}

		if i < len(b) {
			id := fn(b[i])
			m2[id] = b[i]
		}
	}

	res := make([]T, 0, len(a))
	for id, v := range m1 {
		if _, found := m2[id]; !found {
			res = append(res, v)
		}
	}

	return res
}

func intersect[T any, slice []T, S func(elem T) string](a slice, b slice, fn S) slice {
	m1 := make(map[string]T)
	m2 := make(map[string]T)

	limit := len(a)
	if limit < len(b) {
		limit = len(b)
	}

	for i := 0; i < limit; i++ {
		if i < len(a) {
			id := fn(a[i])
			m1[id] = a[i]
		}

		if i < len(b) {
			id := fn(b[i])
			m2[id] = b[i]
		}
	}

	res := make([]T, 0, len(a))
	for id, v := range m1 {
		if _, found := m2[id]; found {
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
