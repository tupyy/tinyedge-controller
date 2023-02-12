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


Table: reference
[ 0] id                                             VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 1] ref_type                                       USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
[ 2] name                                           VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 3] repo_id                                        VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 4] valid                                          BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: []
[ 5] hash                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 6] path_reference                                 TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "JYsicbPebRRBhdiQKyLdmAafS",    "ref_type": "vWkSWBVWQheQVRvFBsiVpHDnN",    "name": "AafOATNqiuDleHJcUYqgoTncT",    "repo_id": "YEAyHFAFLVejjDcHNdKfjWGMQ",    "valid": true,    "hash": "JKFdEnVtrDmIhSZVEEbWIExAG",    "path_reference": "ukVdYBkFaDLRcWIafbWDNrcbi"}



*/

// Reference struct is a row record of the reference table in the tinyedge database
type Reference struct {
	//[ 0] id                                             VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	ID string `gorm:"primary_key;column:id;type:VARCHAR;size:255;"`
	//[ 1] ref_type                                       USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
	RefType string `gorm:"column:ref_type;type:VARCHAR;"`
	//[ 2] name                                           VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	Name string `gorm:"column:name;type:VARCHAR;size:255;"`
	//[ 3] repo_id                                        VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	RepoID string `gorm:"column:repo_id;type:VARCHAR;size:255;"`
	//[ 4] valid                                          BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: []
	Valid bool `gorm:"column:valid;type:BOOL;"`
	//[ 5] hash                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Hash string `gorm:"column:hash;type:TEXT;"`
	//[ 6] path_reference                                 TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	PathReference string `gorm:"column:path_reference;type:TEXT;"`
}

var referenceTableInfo = &TableInfo{
	Name: "reference",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
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
			GoFieldName:        "ID",
			GoFieldType:        "string",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "ref_type",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "USER_DEFINED",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "USER_DEFINED",
			ColumnLength:       -1,
			GoFieldName:        "RefType",
			GoFieldType:        "string",
			JSONFieldName:      "ref_type",
			ProtobufFieldName:  "ref_type",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "repo_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "RepoID",
			GoFieldType:        "string",
			JSONFieldName:      "repo_id",
			ProtobufFieldName:  "repo_id",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "valid",
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
			GoFieldName:        "Valid",
			GoFieldType:        "bool",
			JSONFieldName:      "valid",
			ProtobufFieldName:  "valid",
			ProtobufType:       "bool",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "hash",
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
			GoFieldName:        "Hash",
			GoFieldType:        "string",
			JSONFieldName:      "hash",
			ProtobufFieldName:  "hash",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "path_reference",
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
			GoFieldName:        "PathReference",
			GoFieldType:        "string",
			JSONFieldName:      "path_reference",
			ProtobufFieldName:  "path_reference",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *Reference) TableName() string {
	return "reference"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *Reference) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *Reference) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *Reference) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *Reference) TableInfo() *TableInfo {
	return referenceTableInfo
}
