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


Table: sets_references
[ 0] device_set_id                                  TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] manifest_reference_id                          TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "device_set_id": "CqylhNHdEsrPjXhADdepCUMIx",    "manifest_reference_id": "rLfWWgkMTYAROaLpWXarNIkqn"}



*/

// SetsReferences struct is a row record of the sets_references table in the tinyedge database
type SetsReferences struct {
	//[ 0] device_set_id                                  TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	DeviceSetID string `gorm:"primary_key;column:device_set_id;type:TEXT;"`
	//[ 1] manifest_reference_id                          TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ManifestReferenceID string `gorm:"primary_key;column:manifest_reference_id;type:TEXT;"`
}

var sets_referencesTableInfo = &TableInfo{
	Name: "sets_references",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "device_set_id",
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
			GoFieldName:        "DeviceSetID",
			GoFieldType:        "string",
			JSONFieldName:      "device_set_id",
			ProtobufFieldName:  "device_set_id",
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
func (s *SetsReferences) TableName() string {
	return "sets_references"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (s *SetsReferences) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (s *SetsReferences) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (s *SetsReferences) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (s *SetsReferences) TableInfo() *TableInfo {
	return sets_referencesTableInfo
}
