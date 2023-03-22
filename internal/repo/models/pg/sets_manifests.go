package pg

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

/*
DB Table Details
-------------------------------------


Table: sets_manifests
[ 0] device_set_id                                  VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 1] workload_manifest_id                           VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []


JSON Sample
-------------------------------------
{    "device_set_id": "OvQflvsYYLmPcPVQjywrTWjdW",    "workload_manifest_id": "AYmrYGRxpDKpFTXuBrjnBuGSc"}



*/

// SetsManifests struct is a row record of the sets_manifests table in the tinyedge database
type SetsManifests struct {
	//[ 0] device_set_id                                  VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	DeviceSetID string `gorm:"primary_key;column:device_set_id;type:VARCHAR;size:255;"`
	//[ 1] workload_manifest_id                           VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	WorkloadManifestID string `gorm:"primary_key;column:workload_manifest_id;type:VARCHAR;size:255;"`
}

var sets_manifestsTableInfo = &TableInfo{
	Name: "sets_manifests",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "device_set_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "DeviceSetID",
			GoFieldType:        "string",
			JSONFieldName:      "device_set_id",
			ProtobufFieldName:  "device_set_id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "workload_manifest_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "WorkloadManifestID",
			GoFieldType:        "string",
			JSONFieldName:      "workload_manifest_id",
			ProtobufFieldName:  "workload_manifest_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (s *SetsManifests) TableName() string {
	return "sets_manifests"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (s *SetsManifests) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (s *SetsManifests) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (s *SetsManifests) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (s *SetsManifests) TableInfo() *TableInfo {
	return sets_manifestsTableInfo
}
