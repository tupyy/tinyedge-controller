package postgres

import (
	"context"

	"gorm.io/gorm"
)

func namespaceQuery(db *gorm.DB) *gorm.DB {
	workloadSubQuery := db.Table("manifest").
		Select(`manifest.id as workload_id,manifest.path as workload_path, 
		namespaces_manifests.namespace_id, 
		repo.*
		`).
		Joins("JOIN repo on repo.id = manifest.repo_id").
		Joins("JOIN namespaces_manifests on namespaces_manifests.manifest_id = manifest.id")

	configurationSubQuery := db.Table("manifest").
		Select(`manifest.id as conf_id, manifest.path as conf_path, repo.*
		`).
		Joins("JOIN repo on repo.id = manifest.repo_id")

	return db.Table("namespace").
		Select(`namespace.*,
			device_set.id as set_id,
			device.id as device_id, 
			c.conf_id as conf_id,c.conf_path as conf_path, c.local_path as conf_local_path,
			w.local_path as workload_repo_local_path,w.workload_path as workload_path,w.workload_id as workload_id`).
		Joins("LEFT JOIN device ON device.namespace_id = namespace.id").
		Joins("LEFT JOIN device_set ON device_set.namespace_id = namespace.id").
		Joins("LEFT JOIN (?) as c ON c.conf_id = namespace.configuration_id", configurationSubQuery).
		Joins("LEFT JOIN (?) as w ON w.namespace_id = namespace.id", workloadSubQuery)
}

func setQuery(db *gorm.DB) *gorm.DB {
	workloadSubQuery := db.Table("manifest").
		Select(`manifest.id as workload_id,manifest.path as workload_path, 
		sets_manifests.device_set_id as set_id, 
		repo.*
		`).
		Joins("JOIN repo on repo.id = manifest.repo_id").
		Joins("JOIN sets_manifests on sets_manifests.manifest_id = manifest.id")

	configurationSubQuery := db.Table("manifest").
		Select(`manifest.id as conf_id, manifest.path as conf_path, repo.*
		`).
		Joins("JOIN repo on repo.id = manifest.repo_id")

	return db.Table("device_set").
		Select(`device_set.*,
			device.id as device_id, 
			c.conf_id as conf_id,c.conf_path as conf_path, c.local_path as conf_local_path,
			w.local_path as workload_repo_local_path,w.workload_path as workload_path,w.workload_id as workload_id`).
		Joins("LEFT JOIN device ON device.device_set_id = device_set.id").
		Joins("LEFT JOIN namespace ON namespace.id = device_set.namespace_id").
		Joins("LEFT JOIN (?) as c ON c.conf_id = device_set.configuration_id", configurationSubQuery).
		Joins("LEFT JOIN (?) as w ON w.set_id = device_set.id", workloadSubQuery)
}

func deviceQuery(db *gorm.DB) *gorm.DB {
	workloadSubQuery := db.Table("manifest").
		Select(`manifest.id as workload_id,manifest.path as workload_path, 
		devices_manifests.device_id as device_id,
		repo.*
		`).
		Joins("JOIN repo on repo.id = manifest.repo_id").
		Joins("JOIN devices_manifests on devices_manifests.manifest_id = manifest.id")

	configurationSubQuery := db.Table("manifest").
		Select(`manifest.id as conf_id, manifest.path as conf_path, repo.*
		`).
		Joins("JOIN repo on repo.id = manifest.repo_id")

	return db.Table("device").
		Select(`device.*,
			device_set.id as set_id, 
			c.conf_id as conf_id,c.conf_path as conf_path, c.local_path as conf_local_path,
			w.local_path as workload_repo_local_path,w.workload_path as workload_path,w.workload_id as workload_id`).
		Joins("LEFT JOIN device_set ON device_set.id = device.device_set_id").
		Joins("LEFT JOIN (?) as c ON c.conf_id = device.configuration_id", configurationSubQuery).
		Joins("LEFT JOIN (?) as w ON w.device_id = device.id", workloadSubQuery)
}

type manifestQueryBuilder struct {
	tx *gorm.DB
}

func newManifestQuery(ctx context.Context, db *gorm.DB) *manifestQueryBuilder {
	tx := db.Session(&gorm.Session{SkipHooks: true}).WithContext(ctx).Table("manifest").
		Select(`manifest.*, devices_manifests.device_id as device_id, sets_manifests.device_set_id as set_id, namespaces_manifests.namespace_id as namespace_id,
		repo.id as repo_id, repo.url as repo_url, repo.branch as repo_branch, repo.local_path as repo_local_path,
		repo.auth_type as repo_auth_type, repo.auth_secret_path as repo_auth_secret_path,
		repo.current_head_sha as repo_current_head_sha, repo.target_head_sha as repo_target_head_sha,
		repo.pull_period_seconds as repo_pull_period_seconds`).
		Joins("LEFT JOIN namespaces_manifests ON namespaces_manifests.manifest_id = manifest.id").
		Joins("LEFT JOIN sets_manifests ON sets_manifests.manifest_id = manifest.id").
		Joins("LEFT JOIN devices_manifests ON devices_manifests.manifest_id = manifest.id").
		Joins("JOIN repo ON repo.id = manifest.repo_id")
	return &manifestQueryBuilder{tx}
}

func (mm *manifestQueryBuilder) WithRepoId(id string) *manifestQueryBuilder {
	mm.tx.Where("repo_id = ?", id)
	return mm
}

func (mm *manifestQueryBuilder) WithReferenceID(id string) *manifestQueryBuilder {
	mm.tx.Where("manifest.id = ?", id)
	return mm
}

func (mm *manifestQueryBuilder) WithNamespaceID(id string) *manifestQueryBuilder {
	mm.tx.Where("namespaces_manifests.namespace_id = ?", id)
	return mm
}

func (mm *manifestQueryBuilder) WithDeviceID(id string) *manifestQueryBuilder {
	mm.tx.Where("devices_manifests.device_id = ?", id)
	return mm
}

func (mm *manifestQueryBuilder) WithSetID(id string) *manifestQueryBuilder {
	mm.tx.Where("sets_manifests.device_set_id = ?", id)
	return mm
}

func (mm *manifestQueryBuilder) Build() *gorm.DB {
	return mm.tx
}
