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


Table: device_set
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] configuration_id                               TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] namespace_id                                   TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "acevdvejreMQGKlEldpTXBDCl",    "configuration_id": "bpLWPQePgBakbgRXUModgHrPg",    "namespace_id": "FdQmTniAuuvNHrSkFhKvGUaOg"}



*/

// DeviceSet struct is a row record of the device_set table in the tinyedge database
type DeviceSet struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] configuration_id                               TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ConfigurationID sql.NullString `gorm:"column:configuration_id;type:TEXT;"`
	//[ 2] namespace_id                                   TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	NamespaceID string `gorm:"column:namespace_id;type:TEXT;"`
}

var device_setTableInfo = &TableInfo{
	Name: "device_set",
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
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "namespace_id",
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
			GoFieldName:        "NamespaceID",
			GoFieldType:        "string",
			JSONFieldName:      "namespace_id",
			ProtobufFieldName:  "namespace_id",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *DeviceSet) TableName() string {
	return "device_set"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *DeviceSet) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *DeviceSet) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *DeviceSet) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *DeviceSet) TableInfo() *TableInfo {
	return device_setTableInfo
}
