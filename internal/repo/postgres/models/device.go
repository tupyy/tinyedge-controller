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


Table: device
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] namespace                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] enroled_at                                     TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[ 3] registered_at                                  TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[ 4] enroled                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[ 5] registered                                     BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[ 6] configuration_id                               TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 7] hardware_id                                    TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "JeGPLqrODXjXAmguBngnuxIig",    "namespace": "GslSBjwWCJuFbxxnSjuEHwmSm",    "enroled_at": "2129-12-01T09:26:26.624054247+01:00",    "registered_at": "2238-08-03T04:46:31.849414346+02:00",    "enroled": true,    "registered": true,    "configuration_id": "thodponxyEiLhFkOYTSiUGQLG",    "hardware_id": "MMyhAtvHKRZsoisJCiLQsUxpK"}



*/

// Device struct is a row record of the device table in the tinyedge database
type Device struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] namespace                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Namespace string `gorm:"column:namespace;type:TEXT;"`
	//[ 2] enroled_at                                     TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
	EnroledAt time.Time `gorm:"column:enroled_at;type:TIMESTAMP;"`
	//[ 3] registered_at                                  TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
	RegisteredAt time.Time `gorm:"column:registered_at;type:TIMESTAMP;"`
	//[ 4] enroled                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	Enroled bool `gorm:"column:enroled;type:BOOL;default:false;"`
	//[ 5] registered                                     BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	Registered bool `gorm:"column:registered;type:BOOL;default:false;"`
	//[ 6] configuration_id                               TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ConfigurationID string `gorm:"column:configuration_id;type:TEXT;"`
	//[ 7] hardware_id                                    TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	HardwareID string `gorm:"column:hardware_id;type:TEXT;"`
}

var deviceTableInfo = &TableInfo{
	Name: "device",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
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
			GoFieldName:        "ID",
			GoFieldType:        "string",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "namespace",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "Namespace",
			GoFieldType:        "string",
			JSONFieldName:      "namespace",
			ProtobufFieldName:  "namespace",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "enroled_at",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TIMESTAMP",
			DatabaseTypePretty: "TIMESTAMP",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMP",
			ColumnLength:       -1,
			GoFieldName:        "EnroledAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "enroled_at",
			ProtobufFieldName:  "enroled_at",
			ProtobufType:       "uint64",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "registered_at",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TIMESTAMP",
			DatabaseTypePretty: "TIMESTAMP",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMP",
			ColumnLength:       -1,
			GoFieldName:        "RegisteredAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "registered_at",
			ProtobufFieldName:  "registered_at",
			ProtobufType:       "uint64",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "enroled",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "BOOL",
			DatabaseTypePretty: "BOOL",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "BOOL",
			ColumnLength:       -1,
			GoFieldName:        "Enroled",
			GoFieldType:        "bool",
			JSONFieldName:      "enroled",
			ProtobufFieldName:  "enroled",
			ProtobufType:       "bool",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "registered",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "BOOL",
			DatabaseTypePretty: "BOOL",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "BOOL",
			ColumnLength:       -1,
			GoFieldName:        "Registered",
			GoFieldType:        "bool",
			JSONFieldName:      "registered",
			ProtobufFieldName:  "registered",
			ProtobufType:       "bool",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "configuration_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "ConfigurationID",
			GoFieldType:        "string",
			JSONFieldName:      "configuration_id",
			ProtobufFieldName:  "configuration_id",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "hardware_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "HardwareID",
			GoFieldType:        "string",
			JSONFieldName:      "hardware_id",
			ProtobufFieldName:  "hardware_id",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *Device) TableName() string {
	return "device"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *Device) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *Device) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *Device) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *Device) TableInfo() *TableInfo {
	return deviceTableInfo
}
