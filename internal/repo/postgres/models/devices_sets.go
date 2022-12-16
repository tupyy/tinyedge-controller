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


Table: devices_sets
[ 0] device_id                                      TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] device_set_id                                  TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "device_id": "OZDBDBZtrGgRbcOnhGmQiXiDG",    "device_set_id": "HWuBdJRnhwnBbdiAwRieXEmok"}



*/

// DevicesSets struct is a row record of the devices_sets table in the tinyedge database
type DevicesSets struct {
	//[ 0] device_id                                      TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	DeviceID string `gorm:"primary_key;column:device_id;type:TEXT;"`
	//[ 1] device_set_id                                  TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	DeviceSetID string `gorm:"primary_key;column:device_set_id;type:TEXT;"`
}

var devices_setsTableInfo = &TableInfo{
	Name: "devices_sets",
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
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *DevicesSets) TableName() string {
	return "devices_sets"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *DevicesSets) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *DevicesSets) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *DevicesSets) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *DevicesSets) TableInfo() *TableInfo {
	return devices_setsTableInfo
}
