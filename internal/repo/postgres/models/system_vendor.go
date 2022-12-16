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


Table: system_vendor
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "eZtkuxJQHNVOVbLUKKDhLNyXA"}



*/

// SystemVendor struct is a row record of the system_vendor table in the tinyedge database
type SystemVendor struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
}

var system_vendorTableInfo = &TableInfo{
	Name: "system_vendor",
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
	},
}

// TableName sets the insert table name for this struct type
func (s *SystemVendor) TableName() string {
	return "system_vendor"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (s *SystemVendor) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (s *SystemVendor) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (s *SystemVendor) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (s *SystemVendor) TableInfo() *TableInfo {
	return system_vendorTableInfo
}
