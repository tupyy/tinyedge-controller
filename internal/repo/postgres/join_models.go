package postgres

import "github.com/tupyy/tinyedge-controller/internal/repo/postgres/models"

type SetJoin struct {
	models.DeviceSet
	DeviceId   string `gorm:"column_name:device_id;type:TEXT"`
	ManifestId string `gorm:"column_name:manifest_id;type:TEXT"`
}

type NamespaceJoin struct {
	models.Namespace
	DeviceId   string `gorm:"column_name:device_id;type:TEXT"`
	SetId      string `gorm:"column_name:set_id;type:TEXT"`
	ManifestId string `gorm:"column_name:manifest_id;type:TEXT"`
}
