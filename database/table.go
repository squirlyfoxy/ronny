package database

type ColumnType int

const (
	ColType ColumnType = 0
	KEY                = 0
	STRING             = 1
	INT                = 2
	FLOAT              = 3
	BOOL               = 4
	DATE               = 5
)

//******************
//	  Column Type
//******************
var ColumnTypes = [6]string{
	"key",
	"string",
	"int",
	"float",
	"bool",
	"date",
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
	OType      OnType = 0
	CAN_TAKE          = 0
	CAN_ADD           = 1
	CAN_REMOVE        = 2
	CAN_MODIFY        = 3
)

type RuleType struct {
	Can       []OnType
	TableName string
}

type Rule struct {
	//Name of the column
	RefeersTo string
	RuleTypes []RuleType
}

type Column struct {
	Name string
	Type ColumnType
	Rule ColumnRule
}

type Table struct {
	Name      string
	SubTables []Table
	Columns   []Column
	Functions []Function
	Rule      Rule
}
