package database

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetDataFromATable(database Database, table Table, key int) []string {
	//Check if CAN_GLOBALLY_TAKE
	contains := false

	for _, ts := range table.Rule.RuleTypes {
		if ContainsORtype(ts.Can, CAN_GLOBALLY_TAKE) {
			contains = true
		}
	}
	if !contains {
		fmt.Println("You don't have access to this table")
		return []string{}
	}

	//Read the JSON (./db/data/[tableName].dat.json)
	file, err := os.Open("./db/data/" + table.Name + ".dat.json")
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	defer file.Close()

	//Read the data (TableData)
	var tableData TableData
	err = json.NewDecoder(file).Decode(&tableData)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	//Check the lenght of the data
	if len(tableData.Data) < key || len(tableData.Data) == 0 {
		return []string{}
	}

	//Return the data
	return tableData.Data[key]
}
