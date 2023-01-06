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
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] url                                            TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] branch                                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] local_path                                     TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 4] current_head_sha                               TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 5] target_head_sha                                TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 6] pull_period_seconds                            INT2                 null: true   primary: false  isArray: false  auto: false  col: INT2            len: -1      default: [20]


JSON Sample
-------------------------------------
{    "id": "QigWQgfqXTPrslqjaIHLZaCyf",    "url": "otHygOxfDLmKhgWCLFIhIyiBG",    "branch": "pttMJUwZmhPoeJjQCnJNNpZYC",    "local_path": "OTJBpMMeaYjQixpVKgwJpWINB",    "current_head_sha": "QTwPltowFQnHAVHKMLIQJINxx",    "target_head_sha": "fljgSpahFojKPuEaGattyeFQD",    "pull_period_seconds": 18}



*/

// Repo struct is a row record of the repo table in the tinyedge database
type Repo struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] url                                            TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	URL string `gorm:"column:url;type:TEXT;"`
	//[ 2] branch                                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Branch sql.NullString `gorm:"column:branch;type:TEXT;"`
	//[ 3] local_path                                     TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	LocalPath sql.NullString `gorm:"column:local_path;type:TEXT;"`
	//[ 4] current_head_sha                               TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	CurrentHeadSha sql.NullString `gorm:"column:current_head_sha;type:TEXT;"`
	//[ 5] target_head_sha                                TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	TargetHeadSha sql.NullString `gorm:"column:target_head_sha;type:TEXT;"`
	//[ 6] pull_period_seconds                            INT2                 null: true   primary: false  isArray: false  auto: false  col: INT2            len: -1      default: [20]
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
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
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
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
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
			ProtobufPos:        7,
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
