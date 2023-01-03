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


Table: configuration_cache
[ 0] device_id                                      TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] workload                                       BYTEA                null: true   primary: false  isArray: false  auto: false  col: BYTEA           len: -1      default: []


JSON Sample
-------------------------------------
{    "device_id": "uJqficNqEOjGcixESPauJpAYP",    "workload": "LQULdaxetZVfIVFkkWDwrKWKV"}



*/

// ConfigurationCache struct is a row record of the configuration_cache table in the tinyedge database
type ConfigurationCache struct {
	//[ 0] device_id                                      TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	DeviceID string `gorm:"primary_key;column:device_id;type:TEXT;"`
	//[ 1] workload                                       BYTEA                null: true   primary: false  isArray: false  auto: false  col: BYTEA           len: -1      default: []
	Workload sql.NullString `gorm:"column:workload;type:BYTEA;"`
}

var configuration_cacheTableInfo = &TableInfo{
	Name: "configuration_cache",
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
			Name:               "workload",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "BYTEA",
			DatabaseTypePretty: "BYTEA",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "BYTEA",
			ColumnLength:       -1,
			GoFieldName:        "Workload",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "workload",
			ProtobufFieldName:  "workload",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *ConfigurationCache) TableName() string {
	return "configuration_cache"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *ConfigurationCache) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *ConfigurationCache) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *ConfigurationCache) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *ConfigurationCache) TableInfo() *TableInfo {
	return configuration_cacheTableInfo
}
