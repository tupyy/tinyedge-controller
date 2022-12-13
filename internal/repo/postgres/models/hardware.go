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


Table: hardware
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] os_information_id                              TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] system_vendor_id                               TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "AgaMbaeJtmHEfixBflvAcsRYj",    "os_information_id": "oBQFDMBJAvaTlerrAeRhsBfKA",    "system_vendor_id": "ZBVWuILYbfNwJBtAxivibRxVw"}



*/

// Hardware struct is a row record of the hardware table in the tinyedge database
type Hardware struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] os_information_id                              TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	OsInformationID sql.NullString `gorm:"column:os_information_id;type:TEXT;"`
	//[ 2] system_vendor_id                               TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	SystemVendorID sql.NullString `gorm:"column:system_vendor_id;type:TEXT;"`
}

var hardwareTableInfo = &TableInfo{
	Name: "hardware",
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
			Name:               "os_information_id",
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
			GoFieldName:        "OsInformationID",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "os_information_id",
			ProtobufFieldName:  "os_information_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "system_vendor_id",
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
			GoFieldName:        "SystemVendorID",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "system_vendor_id",
			ProtobufFieldName:  "system_vendor_id",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (h *Hardware) TableName() string {
	return "hardware"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (h *Hardware) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (h *Hardware) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (h *Hardware) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (h *Hardware) TableInfo() *TableInfo {
	return hardwareTableInfo
}
