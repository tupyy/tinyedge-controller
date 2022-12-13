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


Table: os_information
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "IuoowJfycNPeVcfMDyIHrviHh"}



*/

// OsInformation struct is a row record of the os_information table in the tinyedge database
type OsInformation struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
}

var os_informationTableInfo = &TableInfo{
	Name: "os_information",
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
func (o *OsInformation) TableName() string {
	return "os_information"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (o *OsInformation) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (o *OsInformation) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (o *OsInformation) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (o *OsInformation) TableInfo() *TableInfo {
	return os_informationTableInfo
}
