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


Table: secrets_manifests
[ 0] secret_id                                      VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 1] manifest_id                                    VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []


JSON Sample
-------------------------------------
{    "secret_id": "dcBYeGIpPBYhcpalevFXwwEtY",    "manifest_id": "WnuPZxgnnOgCrSmhurhUBNCvJ"}



*/

// SecretsManifests struct is a row record of the secrets_manifests table in the tinyedge database
type SecretsManifests struct {
	//[ 0] secret_id                                      VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	SecretID string `gorm:"primary_key;column:secret_id;type:VARCHAR;size:255;"`
	//[ 1] manifest_id                                    VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	ManifestID string `gorm:"primary_key;column:manifest_id;type:VARCHAR;size:255;"`
}

var secrets_manifestsTableInfo = &TableInfo{
	Name: "secrets_manifests",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "secret_id",
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
			GoFieldName:        "SecretID",
			GoFieldType:        "string",
			JSONFieldName:      "secret_id",
			ProtobufFieldName:  "secret_id",
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
func (s *SecretsManifests) TableName() string {
	return "secrets_manifests"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (s *SecretsManifests) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (s *SecretsManifests) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (s *SecretsManifests) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (s *SecretsManifests) TableInfo() *TableInfo {
	return secrets_manifestsTableInfo
}
