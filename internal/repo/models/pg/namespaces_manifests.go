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


Table: namespaces_manifests
[ 0] namespace_id                                   VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 1] manifest_id                                    VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []


JSON Sample
-------------------------------------
{    "namespace_id": "FnJxCkrSmvtMRZlavJScnSkln",    "manifest_id": "VgBSmbtAHMWwkKvSYPTDDQYDe"}



*/

// NamespacesManifests struct is a row record of the namespaces_manifests table in the tinyedge database
type NamespacesManifests struct {
	//[ 0] namespace_id                                   VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	NamespaceID string `gorm:"primary_key;column:namespace_id;type:VARCHAR;size:255;"`
	//[ 1] manifest_id                                    VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	ManifestID string `gorm:"primary_key;column:manifest_id;type:VARCHAR;size:255;"`
}

var namespaces_manifestsTableInfo = &TableInfo{
	Name: "namespaces_manifests",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "namespace_id",
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
			GoFieldName:        "NamespaceID",
			GoFieldType:        "string",
			JSONFieldName:      "namespace_id",
			ProtobufFieldName:  "namespace_id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "manifest_id",
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
			GoFieldName:        "ManifestID",
			GoFieldType:        "string",
			JSONFieldName:      "manifest_id",
			ProtobufFieldName:  "manifest_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (n *NamespacesManifests) TableName() string {
	return "namespaces_manifests"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (n *NamespacesManifests) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (n *NamespacesManifests) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (n *NamespacesManifests) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (n *NamespacesManifests) TableInfo() *TableInfo {
	return namespaces_manifestsTableInfo
}
