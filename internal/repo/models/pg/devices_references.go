package pg

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
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


Table: devices_references
[ 0] device_id                                      VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 1] reference_id                                   VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []


JSON Sample
-------------------------------------
{    "device_id": "QNxKpEcxKXYaskLFJYdMurNLH",    "reference_id": "MFNraYnoWjmOcHfOXILjSAZCT"}



*/

// DevicesReferences struct is a row record of the devices_references table in the tinyedge database
type DevicesReferences struct {
	//[ 0] device_id                                      VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	DeviceID string `gorm:"primary_key;column:device_id;type:VARCHAR;size:255;"`
	//[ 1] reference_id                                   VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	ReferenceID string `gorm:"primary_key;column:reference_id;type:VARCHAR;size:255;"`
}

var devices_referencesTableInfo = &TableInfo{
	Name: "devices_references",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "device_id",
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
			GoFieldName:        "DeviceID",
			GoFieldType:        "string",
			JSONFieldName:      "device_id",
			ProtobufFieldName:  "device_id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "reference_id",
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
			GoFieldName:        "ReferenceID",
			GoFieldType:        "string",
			JSONFieldName:      "reference_id",
			ProtobufFieldName:  "reference_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *DevicesReferences) TableName() string {
	return "devices_references"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *DevicesReferences) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *DevicesReferences) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *DevicesReferences) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *DevicesReferences) TableInfo() *TableInfo {
	return devices_referencesTableInfo
}
