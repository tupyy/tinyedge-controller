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


Table: manifest_work
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] repo_id                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] path_manifest_work                             TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] content                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "CnpCCWsFxkLUdEylqmfvCuWDf",    "repo_id": "XAKifvaIOUESQmFwGULIToNZe",    "path_manifest_work": "rCwilSwNEttwhPclWavAibdBo",    "content": "gQuGgkrPlxCTbClPDfTaqmAXJ"}



*/

// ManifestWork struct is a row record of the manifest_work table in the tinyedge database
type ManifestWork struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] repo_id                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	RepoID string `gorm:"column:repo_id;type:TEXT;"`
	//[ 2] path_manifest_work                             TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	PathManifestWork string `gorm:"column:path_manifest_work;type:TEXT;"`
	//[ 3] content                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Content string `gorm:"column:content;type:TEXT;"`
}

var manifest_workTableInfo = &TableInfo{
	Name: "manifest_work",
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
			Name:               "path_manifest_work",
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
			GoFieldName:        "PathManifestWork",
			GoFieldType:        "string",
			JSONFieldName:      "path_manifest_work",
			ProtobufFieldName:  "path_manifest_work",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "content",
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
			GoFieldName:        "Content",
			GoFieldType:        "string",
			JSONFieldName:      "content",
			ProtobufFieldName:  "content",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (m *ManifestWork) TableName() string {
	return "manifest_work"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (m *ManifestWork) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (m *ManifestWork) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (m *ManifestWork) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (m *ManifestWork) TableInfo() *TableInfo {
	return manifest_workTableInfo
}
