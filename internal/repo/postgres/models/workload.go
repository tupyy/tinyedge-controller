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


Table: workload
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] path                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "qPAMftJVLtZCqdYiRDASUcCcq",    "path": "UVVGkdKeCgETEcLXKfxAtQoiY"}



*/

// Workload struct is a row record of the workload table in the tinyedge database
type Workload struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] path                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Path string `gorm:"column:path;type:TEXT;"`
}

var workloadTableInfo = &TableInfo{
	Name: "workload",
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
			Name:               "path",
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
			GoFieldName:        "Path",
			GoFieldType:        "string",
			JSONFieldName:      "path",
			ProtobufFieldName:  "path",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (w *Workload) TableName() string {
	return "workload"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (w *Workload) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (w *Workload) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (w *Workload) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (w *Workload) TableInfo() *TableInfo {
	return workloadTableInfo
}
