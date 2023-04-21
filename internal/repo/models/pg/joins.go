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
	DeviceId               string         `gorm:"column:device_id;type:TEXT"`
	SetId                  string         `gorm:"column:set_id;type:TEXT"`
	NamespaceId            string         `gorm:"column:namespace_id;type:TEXT"`
	Repo_ID                string         `gorm:"primary_key;column:repo_id;type:TEXT;"`
	Repo_URL               string         `gorm:"column:repo_url;type:TEXT;"`
	Repo_Branch            sql.NullString `gorm:"column:repo_branch;type:TEXT;"`
	Repo_LocalPath         sql.NullString `gorm:"column:repo_local_path;type:TEXT;"`
	Repo_CurrentHeadSha    sql.NullString `gorm:"column:repo_current_head_sha;type:TEXT;"`
	Repo_TargetHeadSha     sql.NullString `gorm:"column:repo_target_head_sha;type:TEXT;"`
	Repo_PullPeriodSeconds sql.NullInt64  `gorm:"column:repo_pull_period_seconds;type:INT2;default:20;"`
}
