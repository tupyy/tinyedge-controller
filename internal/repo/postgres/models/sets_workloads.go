package models

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


Table: sets_workloads
[ 0] device_set_id                                  TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] workload_id                                    TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "device_set_id": "aHCnqtMRXgBXEsaNKexFoZamG",    "workload_id": "khuBpmNKeUpesLwSxXfXQCFpU"}



*/

// SetsWorkloads struct is a row record of the sets_workloads table in the tinyedge database
type SetsWorkloads struct {
	//[ 0] device_set_id                                  TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	DeviceSetID string `gorm:"primary_key;column:device_set_id;type:TEXT;"`
	//[ 1] workload_id                                    TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	WorkloadID string `gorm:"primary_key;column:workload_id;type:TEXT;"`
}

var sets_workloadsTableInfo = &TableInfo{
	Name: "sets_workloads",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "device_set_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "DeviceSetID",
			GoFieldType:        "string",
			JSONFieldName:      "device_set_id",
			ProtobufFieldName:  "device_set_id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "workload_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "WorkloadID",
			GoFieldType:        "string",
			JSONFieldName:      "workload_id",
			ProtobufFieldName:  "workload_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (s *SetsWorkloads) TableName() string {
	return "sets_workloads"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (s *SetsWorkloads) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (s *SetsWorkloads) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (s *SetsWorkloads) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (s *SetsWorkloads) TableInfo() *TableInfo {
	return sets_workloadsTableInfo
}
