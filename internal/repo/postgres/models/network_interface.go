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


Table: network_interface
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] hardware_id                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] name                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] mac_address                                    USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
[ 4] has_carrier                                    BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: []
[ 5] ip4                                            _INET                null: true   primary: false  isArray: false  auto: false  col: _INET           len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "WVRnxkPSAvtJJqKkOnliwRbYw",    "hardware_id": "PwLbcfOQCEhyKbWiwhnOWJqDw",    "name": "jhaXtKcNFUyXwOvlyKbJCmguX",    "mac_address": "nfOGMljLbEKoylhiKkFuQKWyQ",    "has_carrier": false}



*/

// NetworkInterface struct is a row record of the network_interface table in the tinyedge database
type NetworkInterface struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"`
	//[ 1] hardware_id                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	HardwareID sql.NullString `gorm:"column:hardware_id;type:TEXT;"`
	//[ 2] name                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Name string `gorm:"column:name;type:TEXT;"`
	//[ 3] mac_address                                    USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
	MacAddress string `gorm:"column:mac_address;type:VARCHAR;"`
	//[ 4] has_carrier                                    BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: []
	HasCarrier bool `gorm:"column:has_carrier;type:BOOL;"`
}

var network_interfaceTableInfo = &TableInfo{
	Name: "network_interface",
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
			Name:               "hardware_id",
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
			GoFieldName:        "HardwareID",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "hardware_id",
			ProtobufFieldName:  "hardware_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "name",
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
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "mac_address",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "USER_DEFINED",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "USER_DEFINED",
			ColumnLength:       -1,
			GoFieldName:        "MacAddress",
			GoFieldType:        "string",
			JSONFieldName:      "mac_address",
			ProtobufFieldName:  "mac_address",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "has_carrier",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "BOOL",
			DatabaseTypePretty: "BOOL",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "BOOL",
			ColumnLength:       -1,
			GoFieldName:        "HasCarrier",
			GoFieldType:        "bool",
			JSONFieldName:      "has_carrier",
			ProtobufFieldName:  "has_carrier",
			ProtobufType:       "bool",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (n *NetworkInterface) TableName() string {
	return "network_interface"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (n *NetworkInterface) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (n *NetworkInterface) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (n *NetworkInterface) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (n *NetworkInterface) TableInfo() *TableInfo {
	return network_interfaceTableInfo
}
