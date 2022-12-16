package models

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
[ 1] hardware_profile_scope                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: [full]
[ 2] hardware_profile_include                       BOOL                 null: true   primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [true]
[ 3] heartbeat_period_seconds                       INT2                 null: true   primary: false  isArray: false  auto: false  col: INT2            len: -1      default: [30]


JSON Sample
-------------------------------------
{    "id": "btIWSRMneAIkMPeLaawcqMWHc",    "hardware_profile_scope": "FuhWaQAmmGPVaSyAAGBymJiiU",    "hardware_profile_include": false,    "heartbeat_period_seconds": 36}



*/

// Configuration struct is a row record of the configuration table in the tinyedge database
type Configuration struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] hardware_profile_scope                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: [full]
	HardwareProfileScope sql.NullString `gorm:"column:hardware_profile_scope;type:TEXT;default:full;"`
	//[ 2] hardware_profile_include                       BOOL                 null: true   primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [true]
	HardwareProfileInclude sql.NullBool `gorm:"column:hardware_profile_include;type:BOOL;default:true;"`
	//[ 3] heartbeat_period_seconds                       INT2                 null: true   primary: false  isArray: false  auto: false  col: INT2            len: -1      default: [30]
	HeartbeatPeriodSeconds sql.NullInt64 `gorm:"column:heartbeat_period_seconds;type:INT2;default:30;"`
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
			Name:               "hardware_profile_scope",
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
			GoFieldName:        "HardwareProfileScope",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "hardware_profile_scope",
			ProtobufFieldName:  "hardware_profile_scope",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "hardware_profile_include",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "BOOL",
			DatabaseTypePretty: "BOOL",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "BOOL",
			ColumnLength:       -1,
			GoFieldName:        "HardwareProfileInclude",
			GoFieldType:        "sql.NullBool",
			JSONFieldName:      "hardware_profile_include",
			ProtobufFieldName:  "hardware_profile_include",
			ProtobufType:       "bool",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
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
			ProtobufPos:        4,
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
