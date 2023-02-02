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


Table: manifest_reference
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] repo_id                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] valid                                          BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: []
[ 3] hash                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 4] path_manifest_reference                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "tYIXwpNPEjqnKBlssQYTAUpom",    "repo_id": "ZFHXfSVyZVjsUAIEJKLiijkry",    "valid": false,    "hash": "uKVhYKScbCMurkFGfOxCwiYGm",    "path_manifest_reference": "hVkUVvaWCesNcfwNQwQsnnMFo"}



*/

// ManifestReference struct is a row record of the manifest_reference table in the tinyedge database
type ManifestReference struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] repo_id                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	RepoID string `gorm:"column:repo_id;type:TEXT;"`
	//[ 2] valid                                          BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: []
	Valid bool `gorm:"column:valid;type:BOOL;"`
	//[ 3] hash                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Hash string `gorm:"column:hash;type:TEXT;"`
	//[ 4] path_manifest_reference                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	PathManifestReference string `gorm:"column:path_manifest_reference;type:TEXT;"`
}

var manifest_referenceTableInfo = &TableInfo{
	Name: "manifest_reference",
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
			Name:               "repo_id",
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
			GoFieldName:        "RepoID",
			GoFieldType:        "string",
			JSONFieldName:      "repo_id",
			ProtobufFieldName:  "repo_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
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
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
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
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "path_manifest_reference",
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
			GoFieldName:        "PathManifestReference",
			GoFieldType:        "string",
			JSONFieldName:      "path_manifest_reference",
			ProtobufFieldName:  "path_manifest_reference",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (m *ManifestReference) TableName() string {
	return "manifest_reference"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (m *ManifestReference) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (m *ManifestReference) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (m *ManifestReference) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (m *ManifestReference) TableInfo() *TableInfo {
	return manifest_referenceTableInfo
}
