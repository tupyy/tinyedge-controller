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


Table: configuration
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] heartbeat_period_seconds                       INT2                 null: true   primary: false  isArray: false  auto: false  col: INT2            len: -1      default: [30]
[ 2] log_level                                      TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: [info]


JSON Sample
-------------------------------------
{    "id": "oxxdUCmbMJFnOafeLxFkvdRdU",    "heartbeat_period_seconds": 66,    "log_level": "tsZLTTvqAiuDXEkUofWgNqcVe"}



*/

// Configuration struct is a row record of the configuration table in the tinyedge database
type Configuration struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] heartbeat_period_seconds                       INT2                 null: true   primary: false  isArray: false  auto: false  col: INT2            len: -1      default: [30]
	HeartbeatPeriodSeconds sql.NullInt64 `gorm:"column:heartbeat_period_seconds;type:INT2;default:30;"`
	//[ 2] log_level                                      TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: [info]
	LogLevel sql.NullString `gorm:"column:log_level;type:TEXT;default:info;"`
}

var configurationTableInfo = &TableInfo{
	Name: "configuration",
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
			Name:               "heartbeat_period_seconds",
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
			GoFieldName:        "HeartbeatPeriodSeconds",
			GoFieldType:        "sql.NullInt64",
			JSONFieldName:      "heartbeat_period_seconds",
			ProtobufFieldName:  "heartbeat_period_seconds",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "log_level",
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
			GoFieldName:        "LogLevel",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "log_level",
			ProtobufFieldName:  "log_level",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *Configuration) TableName() string {
	return "configuration"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *Configuration) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *Configuration) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *Configuration) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *Configuration) TableInfo() *TableInfo {
	return configurationTableInfo
}
