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


Table: namespace
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] is_default                                     BOOL                 null: true   primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[ 2] configuration_id                               TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "BondLgwgbhAvjbhEuNoekyQZa",    "is_default": false,    "configuration_id": "ffRSTDwwqDJEhTsbtytivsJHB"}



*/

// Namespace struct is a row record of the namespace table in the tinyedge database
type Namespace struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] is_default                                     BOOL                 null: true   primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	IsDefault sql.NullBool `gorm:"column:is_default;type:BOOL;default:false;"`
	//[ 2] configuration_id                               TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ConfigurationID string `gorm:"column:configuration_id;type:TEXT;"`
}

var namespaceTableInfo = &TableInfo{
	Name: "namespace",
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
			Name:               "is_default",
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
			GoFieldName:        "IsDefault",
			GoFieldType:        "sql.NullBool",
			JSONFieldName:      "is_default",
			ProtobufFieldName:  "is_default",
			ProtobufType:       "bool",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "configuration_id",
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
			GoFieldName:        "ConfigurationID",
			GoFieldType:        "string",
			JSONFieldName:      "configuration_id",
			ProtobufFieldName:  "configuration_id",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (n *Namespace) TableName() string {
	return "namespace"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (n *Namespace) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (n *Namespace) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (n *Namespace) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (n *Namespace) TableInfo() *TableInfo {
	return namespaceTableInfo
}
