package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

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
	Name         string     `json:"name"`
	Type         ColumnType `json:"type"`
	TypeAsString string     `json:"typeAsString"`
	IsArray      bool       `json:"isArray"`  //If the column is an array (type starts with [])
	IsExtern     bool       `json:"isExtern"` //If
	Rule         ColumnRule `json:"rule"`
}

type TableData struct {
	Columns []Column        `json:"columns"` //Will store the name of the columns
	Data    [][]interface{} `json:"data"`    //Will store the data of the columns
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

//Saved in "./db/data/t_ids.json"
type TablesIDs struct {
	TableName string `json:"tableName"`
	LatestsID int    `json:"latestID"`
}

var TablesIDsList []TablesIDs = make([]TablesIDs, 0)

func UnmarshallTablesIDsList() {
	//Unmarshall the json file
	jsonFile, err := os.Open("./db/data/t_ids.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &TablesIDsList)
}

func SaveTablesIDsList() {
	josnData, _ := json.Marshal(TablesIDsList)

	//Write the json file
	err := ioutil.WriteFile("./db/data/t_ids.json", josnData, 0644)
	if err != nil {
		panic(err)
	}
}

func AddTableToIDsList(table string) {
	//Check if TablesIDsList is empty
	if len(TablesIDsList) == 0 {
		//Unmarshall the json file
		jsonFile, err := os.Open("./db/data/t_ids.json")
		if err != nil {
			panic(err)
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &TablesIDsList)
	}

	//Check if the table is already in the list
	for _, t := range TablesIDsList {
		if t.TableName == table {
			return
		}
	}

	//Add the table to the list
	TablesIDsList = append(TablesIDsList, TablesIDs{
		TableName: table,
	})
}

func (t Table) GetLatestAutoIncrementedKey() int {
	//Check if TablesIDsList is empty
	if len(TablesIDsList) == 0 {
		//Unmarshall the json file
		jsonFile, err := os.Open("./db/data/t_ids.json")
		if err != nil {
			panic(err)
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &TablesIDsList)
	}

	//Get the first TRS of the table, if the array is empty, return LatestsID
	for k, ts := range TablesIDsList {
		if ts.TableName == t.Name {
			ts.LatestsID += 1
			TablesIDsList[k] = ts

			//Save the new value
			SaveTablesIDsList()

			return ts.LatestsID
		}
	}

	return -1
}
