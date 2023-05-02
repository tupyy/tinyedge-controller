package pg

import (
	"database/sql"
)

type DeviceJoin struct {
	Device
	ConfigurationID        string `gorm:"column:conf_id;type:TEXT"`
	ConfigurationPath      string `gorm:"column:conf_path;type:TEXT"`
	ConfigurationLocalPath string `gorm:"column:conf_local_path;type:TEXT"`
	WorkloadID             string `gorm:"column:workload_id;type:TEXT"`
	WorkloadRepoLocalPath  string `gorm:"column:workload_repo_local_path;type:TEXT"`
	WorkloadPath           string `gorm:"column:workload_path;type:TEXT"`
}

type SetJoin struct {
	DeviceSet
	DeviceId               string `gorm:"column:device_id;type:TEXT"`
	ConfigurationID        string `gorm:"column:conf_id;type:TEXT"`
	ConfigurationPath      string `gorm:"column:conf_path;type:TEXT"`
	ConfigurationLocalPath string `gorm:"column:conf_local_path;type:TEXT"`
	WorkloadID             string `gorm:"column:workload_id;type:TEXT"`
	WorkloadRepoLocalPath  string `gorm:"column:workload_repo_local_path;type:TEXT"`
	WorkloadPath           string `gorm:"column:workload_path;type:TEXT"`
}

type NamespaceJoin struct {
	Namespace
	DeviceId               string `gorm:"column:device_id;type:TEXT"`
	SetId                  string `gorm:"column:set_id;type:TEXT"`
	ConfigurationID        string `gorm:"column:conf_id;type:TEXT"`
	ConfigurationPath      string `gorm:"column:conf_path;type:TEXT"`
	ConfigurationLocalPath string `gorm:"column:conf_local_path;type:TEXT"`
	WorkloadID             string `gorm:"column:workload_id;type:TEXT"`
	WorkloadRepoLocalPath  string `gorm:"column:workload_repo_local_path;type:TEXT"`
	WorkloadPath           string `gorm:"column:workload_path;type:TEXT"`
}

type ManifestJoin struct {
	Manifest
	DeviceId              string         `gorm:"column:device_id;type:TEXT"`
	SetId                 string         `gorm:"column:set_id;type:TEXT"`
	NamespaceId           string         `gorm:"column:namespace_id;type:TEXT"`
	RepoID                string         `gorm:"column:repo_id;type:TEXT;"`
	RepoURL               string         `gorm:"column:repo_url;type:TEXT;"`
	RepoBranch            sql.NullString `gorm:"column:repo_branch;type:TEXT;"`
	RepoLocalPath         sql.NullString `gorm:"column:repo_local_path;type:TEXT;"`
	RepoAuth              sql.NullString `gorm:"column:repo_auth_type;type:TEXT;"`
	RepoAuthSecret        sql.NullString `gorm:"column:repo_auth_secret_path;type:TEXT;"`
	RepoCurrentHeadSha    sql.NullString `gorm:"column:repo_current_head_sha;type:TEXT;"`
	RepoTargetHeadSha     sql.NullString `gorm:"column:repo_target_head_sha;type:TEXT;"`
	RepoPullPeriodSeconds sql.NullInt64  `gorm:"column:repo_pull_period_seconds;type:INT2;default:20;"`
}
