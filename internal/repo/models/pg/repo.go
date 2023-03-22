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


Table: repo
[ 0] id                                             VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 1] url                                            TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] branch                                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] local_path                                     TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 4] auth_type                                      VARCHAR(20)          null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 20      default: []
[ 5] auth_secret_path                               VARCHAR(20)          null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 20      default: []
[ 6] current_head_sha                               TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 7] target_head_sha                                TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 8] pull_period_seconds                            INT2                 null: true   primary: false  isArray: false  auto: false  col: INT2            len: -1      default: [20]


JSON Sample
-------------------------------------
{    "id": "brtNBGgbvkvDsplHtlqcolWDW",    "url": "cYEOKLfWRKIxWeiykimTLdHJa",    "branch": "tbQACFPwqlGRlpUSJwqrnYOqy",    "local_path": "hRHnYAEElDBqYgeKWlvTYEtSS",    "auth_type": "usBYqRRZauMOvdNTTGXCFRujD",    "auth_secret_path": "QLsWBdUNqjaMUQcHZncQRcJhd",    "current_head_sha": "NwRPyfPPVuffvXyOVVLaeTlob",    "target_head_sha": "BvykNGkCpXVUqoTvDqRjTSXNx",    "pull_period_seconds": 33}



*/

// Repo struct is a row record of the repo table in the tinyedge database
type Repo struct {
	//[ 0] id                                             VARCHAR(255)         null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	ID string `gorm:"primary_key;column:id;type:VARCHAR;size:255;"`
	//[ 1] url                                            TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	URL string `gorm:"column:url;type:TEXT;"`
	//[ 2] branch                                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Branch sql.NullString `gorm:"column:branch;type:TEXT;"`
	//[ 3] local_path                                     TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	LocalPath sql.NullString `gorm:"column:local_path;type:TEXT;"`
	//[ 4] auth_type                                      VARCHAR(20)          null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 20      default: []
	AuthType sql.NullString `gorm:"column:auth_type;type:VARCHAR;size:20;"`
	//[ 5] auth_secret_path                               VARCHAR(20)          null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 20      default: []
	AuthSecretPath sql.NullString `gorm:"column:auth_secret_path;type:VARCHAR;size:20;"`
	//[ 6] current_head_sha                               TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	CurrentHeadSha sql.NullString `gorm:"column:current_head_sha;type:TEXT;"`
	//[ 7] target_head_sha                                TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	TargetHeadSha sql.NullString `gorm:"column:target_head_sha;type:TEXT;"`
	//[ 8] pull_period_seconds                            INT2                 null: true   primary: false  isArray: false  auto: false  col: INT2            len: -1      default: [20]
	PullPeriodSeconds sql.NullInt64 `gorm:"column:pull_period_seconds;type:INT2;default:20;"`
}

var repoTableInfo = &TableInfo{
	Name: "repo",
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
			Name:               "url",
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
			GoFieldName:        "URL",
			GoFieldType:        "string",
			JSONFieldName:      "url",
			ProtobufFieldName:  "url",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "branch",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "Branch",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "branch",
			ProtobufFieldName:  "branch",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "local_path",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "LocalPath",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "local_path",
			ProtobufFieldName:  "local_path",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "auth_type",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(20)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       20,
			GoFieldName:        "AuthType",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "auth_type",
			ProtobufFieldName:  "auth_type",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "auth_secret_path",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(20)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       20,
			GoFieldName:        "AuthSecretPath",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "auth_secret_path",
			ProtobufFieldName:  "auth_secret_path",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "current_head_sha",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "CurrentHeadSha",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "current_head_sha",
			ProtobufFieldName:  "current_head_sha",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "target_head_sha",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "TargetHeadSha",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "target_head_sha",
			ProtobufFieldName:  "target_head_sha",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "pull_period_seconds",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "INT2",
			DatabaseTypePretty: "INT2",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT2",
			ColumnLength:       -1,
			GoFieldName:        "PullPeriodSeconds",
			GoFieldType:        "sql.NullInt64",
			JSONFieldName:      "pull_period_seconds",
			ProtobufFieldName:  "pull_period_seconds",
			ProtobufType:       "int32",
			ProtobufPos:        9,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *Repo) TableName() string {
	return "repo"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *Repo) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *Repo) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *Repo) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *Repo) TableInfo() *TableInfo {
	return repoTableInfo
}
