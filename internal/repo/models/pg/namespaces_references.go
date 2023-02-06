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


Table: namespaces_references
[ 0] namespace_id                                   TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] manifest_reference_id                          TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "namespace_id": "PiMmafSFYZQBeHZJPbRLAWaVR",    "manifest_reference_id": "WYYlNyvtLmCnveHhhOveIGsbn"}



*/

// NamespacesReferences struct is a row record of the namespaces_references table in the tinyedge database
type NamespacesReferences struct {
	//[ 0] namespace_id                                   TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	NamespaceID string `gorm:"primary_key;column:namespace_id;type:TEXT;"`
	//[ 1] manifest_reference_id                          TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ManifestReferenceID string `gorm:"primary_key;column:manifest_reference_id;type:TEXT;"`
}

var namespaces_referencesTableInfo = &TableInfo{
	Name: "namespaces_references",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "namespace_id",
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
			GoFieldName:        "NamespaceID",
			GoFieldType:        "string",
			JSONFieldName:      "namespace_id",
			ProtobufFieldName:  "namespace_id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "manifest_reference_id",
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
			GoFieldName:        "ManifestReferenceID",
			GoFieldType:        "string",
			JSONFieldName:      "manifest_reference_id",
			ProtobufFieldName:  "manifest_reference_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (n *NamespacesReferences) TableName() string {
	return "namespaces_references"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (n *NamespacesReferences) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (n *NamespacesReferences) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (n *NamespacesReferences) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (n *NamespacesReferences) TableInfo() *TableInfo {
	return namespaces_referencesTableInfo
}
