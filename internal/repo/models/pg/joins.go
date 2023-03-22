package pg

import (
	"database/sql"
)

type DeviceJoin struct {
	Device
	ManifestId string `gorm:"column_name:manifest_id;type:TEXT"`
}

type SetJoin struct {
	DeviceSet
	DeviceId                            string         `gorm:"column_name:device_id;type:TEXT"`
	ManifestId                          string         `gorm:"column_name:manifest_id;type:TEXT"`
	ConfigurationID                     string         `gorm:"column_name:configuration_id;type:TEXT"`
	ConfigurationHeartbeatPeriodSeconds sql.NullInt64  `gorm:"column:configuration_heartbeat_period_seconds;type:INT2;default:30;"`
	ConfigurationLogLevel               sql.NullString `gorm:"column:configuration_log_level;type:TEXT;default:info;"`
}

type NamespaceJoin struct {
	Namespace
	DeviceId                            string         `gorm:"column_name:device_id;type:TEXT"`
	SetId                               string         `gorm:"column_name:set_id;type:TEXT"`
	ManifestId                          string         `gorm:"column_name:manifest_id;type:TEXT"`
	ConfigurationID                     string         `gorm:"column_name:configuration_id;type:TEXT"`
	ConfigurationHeartbeatPeriodSeconds sql.NullInt64  `gorm:"column:configuration_heartbeat_period_seconds;type:INT2;default:30;"`
	ConfigurationLogLevel               sql.NullString `gorm:"column:configuration_log_level;type:TEXT;default:info;"`
}

type ReferenceJoin struct {
	Reference
	DeviceId               string         `gorm:"column_name:device_id;type:TEXT"`
	SetId                  string         `gorm:"column_name:set_id;type:TEXT"`
	NamespaceId            string         `gorm:"column_name:namespace_id;type:TEXT"`
	Repo_ID                string         `gorm:"primary_key;column:repo_id;type:TEXT;"`
	Repo_URL               string         `gorm:"column:repo_url;type:TEXT;"`
	Repo_Branch            sql.NullString `gorm:"column:repo_branch;type:TEXT;"`
	Repo_LocalPath         sql.NullString `gorm:"column:repo_local_path;type:TEXT;"`
	Repo_CurrentHeadSha    sql.NullString `gorm:"column:repo_current_head_sha;type:TEXT;"`
	Repo_TargetHeadSha     sql.NullString `gorm:"column:repo_target_head_sha;type:TEXT;"`
	Repo_PullPeriodSeconds sql.NullInt64  `gorm:"column:repo_pull_period_seconds;type:INT2;default:20;"`
}
