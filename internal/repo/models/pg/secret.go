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


Table: secret
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] path                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] current_hash                                   TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] target_hash                                    TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "oAxVNLONKhGyQVJNUvLofBgcT",    "path": "NhPQBBEpJAexicTLKmgxSObeN",    "current_hash": "AbIimyTZSQfiFkNTfifcZvNFd",    "target_hash": "bLXKuiJaUCyqFTUuwjmIHkaZh"}



*/

// Secret struct is a row record of the secret table in the tinyedge database
type Secret struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] path                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Path string `gorm:"column:path;type:TEXT;"`
	//[ 2] current_hash                                   TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	CurrentHash string `gorm:"column:current_hash;type:TEXT;"`
	//[ 3] target_hash                                    TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	TargetHash string `gorm:"column:target_hash;type:TEXT;"`
}

var secretTableInfo = &TableInfo{
	Name: "secret",
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

		&ColumnInfo{
			Index:              2,
			Name:               "current_hash",
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
			GoFieldName:        "CurrentHash",
			GoFieldType:        "string",
			JSONFieldName:      "current_hash",
			ProtobufFieldName:  "current_hash",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "target_hash",
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
			GoFieldName:        "TargetHash",
			GoFieldType:        "string",
			JSONFieldName:      "target_hash",
			ProtobufFieldName:  "target_hash",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (s *Secret) TableName() string {
	return "secret"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (s *Secret) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (s *Secret) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (s *Secret) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (s *Secret) TableInfo() *TableInfo {
	return secretTableInfo
}
