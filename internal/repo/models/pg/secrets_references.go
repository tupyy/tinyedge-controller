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


Table: secrets_references
[ 0] secret_id                                      VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 1] reference_id                                   VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []


JSON Sample
-------------------------------------
{    "secret_id": "yHjMsXxaHqMrBwNIeieunFUeQ",    "reference_id": "SkRgwPjmFYRaVNkeepxlOZJGI"}



*/

// SecretsReferences struct is a row record of the secrets_references table in the tinyedge database
type SecretsReferences struct {
	//[ 0] secret_id                                      VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	SecretID string `gorm:"primary_key;column:secret_id;type:VARCHAR;size:255;"`
	//[ 1] reference_id                                   VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	ReferenceID string `gorm:"primary_key;column:reference_id;type:VARCHAR;size:255;"`
}

var secrets_referencesTableInfo = &TableInfo{
	Name: "secrets_references",
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
			Name:               "reference_id",
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
			GoFieldName:        "ReferenceID",
			GoFieldType:        "string",
			JSONFieldName:      "reference_id",
			ProtobufFieldName:  "reference_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (s *SecretsReferences) TableName() string {
	return "secrets_references"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (s *SecretsReferences) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (s *SecretsReferences) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (s *SecretsReferences) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (s *SecretsReferences) TableInfo() *TableInfo {
	return secrets_referencesTableInfo
}
