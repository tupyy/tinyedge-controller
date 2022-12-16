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


Table: devices_workloads
[ 0] device_id                                      TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] workload_id                                    TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "device_id": "QiZRSllwUjRMcBfSYKChkXxWH",    "workload_id": "LhsDtgTWluCBhesILKHBJmQvu"}



*/

// DevicesWorkloads struct is a row record of the devices_workloads table in the tinyedge database
type DevicesWorkloads struct {
	//[ 0] device_id                                      TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	DeviceID string `gorm:"primary_key;column:device_id;type:TEXT;"`
	//[ 1] workload_id                                    TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	WorkloadID string `gorm:"primary_key;column:workload_id;type:TEXT;"`
}

var devices_workloadsTableInfo = &TableInfo{
	Name: "devices_workloads",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "device_id",
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
			GoFieldName:        "DeviceID",
			GoFieldType:        "string",
			JSONFieldName:      "device_id",
			ProtobufFieldName:  "device_id",
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
func (d *DevicesWorkloads) TableName() string {
	return "devices_workloads"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *DevicesWorkloads) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *DevicesWorkloads) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *DevicesWorkloads) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *DevicesWorkloads) TableInfo() *TableInfo {
	return devices_workloadsTableInfo
}
