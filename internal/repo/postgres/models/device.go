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
[ 1] enroled_at                                     TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[ 2] registered_at                                  TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[ 3] enroled                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: [not_enroled]
[ 4] registered                                     BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[ 5] certificate_sn                                 TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 6] namespace_id                                   TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 7] device_set_id                                  TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 8] configuration_id                               TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "skXHfpWkwcFYJDMvemESFZNhK",    "enroled_at": "2305-04-30T08:27:44.61348959+02:00",    "registered_at": "2198-10-10T01:24:06.252049751+02:00",    "enroled": "NokleHhEpjTeVWLqFenQkhoKH",    "registered": true,    "certificate_sn": "MMgJWhhMrmWHEsxDShoyKSegt",    "namespace_id": "JMjuIxjefGNmNkBwmVJorCIqe",    "device_set_id": "fbtoQARguvhHchGWyPdaEcQdS",    "configuration_id": "AoWwLKuLbgyWyDDPKQeOioPoU"}



*/

// Device struct is a row record of the device table in the tinyedge database
type Device struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] enroled_at                                     TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
	EnroledAt time.Time `gorm:"column:enroled_at;type:TIMESTAMP;"`
	//[ 2] registered_at                                  TIMESTAMP            null: true   primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
	RegisteredAt time.Time `gorm:"column:registered_at;type:TIMESTAMP;"`
	//[ 3] enroled                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: [not_enroled]
	Enroled string `gorm:"column:enroled;type:TEXT;default:not_enroled;"`
	//[ 4] registered                                     BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	Registered bool `gorm:"column:registered;type:BOOL;default:false;"`
	//[ 5] certificate_sn                                 TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	CertificateSn sql.NullString `gorm:"column:certificate_sn;type:TEXT;"`
	//[ 6] namespace_id                                   TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	NamespaceID sql.NullString `gorm:"column:namespace_id;type:TEXT;"`
	//[ 7] device_set_id                                  TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	DeviceSetID sql.NullString `gorm:"column:device_set_id;type:TEXT;"`
	//[ 8] configuration_id                               TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ConfigurationID sql.NullString `gorm:"column:configuration_id;type:TEXT;"`
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
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
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
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "enroled",
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
			GoFieldName:        "Enroled",
			GoFieldType:        "string",
			JSONFieldName:      "enroled",
			ProtobufFieldName:  "enroled",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
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
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "certificate_sn",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "CertificateSn",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "certificate_sn",
			ProtobufFieldName:  "certificate_sn",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "namespace_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "NamespaceID",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "namespace_id",
			ProtobufFieldName:  "namespace_id",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "device_set_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "DeviceSetID",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "device_set_id",
			ProtobufFieldName:  "device_set_id",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "configuration_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "ConfigurationID",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "configuration_id",
			ProtobufFieldName:  "configuration_id",
			ProtobufType:       "string",
			ProtobufPos:        9,
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
