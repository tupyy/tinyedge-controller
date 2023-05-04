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


Table: manifest
[ 0] id                                             VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 1] version                                        VARCHAR(30)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 30      default: []
[ 2] repo_id                                        VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 3] path                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "hXiyjXhJDXOKlotoiaEtAsVLY",    "version": "wcDKIyvUhLEZZthKtpEpAwWBu",    "repo_id": "rQeFyuUZaibZXGJRquBrvxCAw",    "path": "ryYaRrUQiVXvdAUkKkCBgKlbe"}



*/

// Manifest struct is a row record of the manifest table in the tinyedge database
type Manifest struct {
	//[ 0] id                                             VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	ID string `gorm:"primary_key;column:id;type:VARCHAR;size:255;"`
	//[ 1] version                                        VARCHAR(30)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 30      default: []
	Version string `gorm:"column:version;type:VARCHAR;size:30;"`
	//[ 2] repo_id                                        VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	RepoID string `gorm:"column:repo_id;type:VARCHAR;size:255;"`
	//[ 3] path                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Path string `gorm:"column:path;type:TEXT;"`
}

var manifestTableInfo = &TableInfo{
	Name: "manifest",
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
			Name:               "version",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(30)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       30,
			GoFieldName:        "Version",
			GoFieldType:        "string",
			JSONFieldName:      "version",
			ProtobufFieldName:  "version",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
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
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
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
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (m *Manifest) TableName() string {
	return "manifest"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (m *Manifest) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (m *Manifest) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (m *Manifest) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (m *Manifest) TableInfo() *TableInfo {
	return manifestTableInfo
}
