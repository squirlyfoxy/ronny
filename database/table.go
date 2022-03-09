package database

type ColumnType int

const (
	ColType  ColumnType = 0
	KEY                 = 0
	STRING              = 1
	INT                 = 2
	FLOAT               = 3
	BOOL                = 4
	DATE                = 5
	EXTERNAL            = 6
)

//******************
//	  Column Type
//******************
var ColumnTypes = [7]string{
	"key",
	"string",
	"int",
	"float",
	"bool",
	"date",
	"external",
}

type ColumnRule int

const (
	CRule         ColumnRule = 0
	DEFAULT                  = 0
	AUTOINCREMENT            = 1
	UNIQUE                   = 2
	NOT_NULL                 = 3
	USERACCESSKEY            = 4
)

type OnType int

const (
	OType             OnType = 0
	CAN_GLOBALLY_TAKE        = 0
	CAN_TAKE                 = 1
	CAN_ADD                  = 2
	CAN_REMOVE               = 3
	CAN_MODIFY               = 4
)

type RuleType struct {
	Can       []OnType `json:"can"`
	TableName string   `json:"tableName"`
}

type Rule struct {
	RefeersTo string     `json:"refersTo"` //Name of the column that this rule is referring to
	RuleTypes []RuleType `json:"ruleTypes"`
}

type Column struct {
	Name    string     `json:"name"`
	Type    ColumnType `json:"type"`
	IsArray bool       `json:"isArray"` //If the column is an array (type starts with [])
	Rule    ColumnRule `json:"rule"`
}

type TableData struct {
	Columns []Column   `json:"columns"` //Will store the name of the columns
	Data    [][]string `json:"data"`    //Will store the data
}

type ExternalType struct {
	Type       string `json:"type"`
	ColumnName string `json:"columnName"`
}

type Table struct {
	Name          string         `json:"name"`
	SubTables     []Table        `json:"subTables"`
	Columns       []Column       `json:"columns"`
	Functions     []Function     `json:"functions"`
	Rule          Rule           `json:"rule"`
	ExternalTypes []ExternalType `json:"externalTypes"`
}
