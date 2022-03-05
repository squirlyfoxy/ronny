package database

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetDataFromATable(database Database, table Table, key int) []string {
	var data []string

	//Read the JSON (./db/data/[tableName].dat.json)
	file, err := os.Open("./db/data/" + table.Name + ".dat.json")
	if err != nil {
		fmt.Println(err)
		return data
	}
	defer file.Close()

	//Read the data (TableData)
	var tableData TableData
	err = json.NewDecoder(file).Decode(&tableData)
	if err != nil {
		fmt.Println(err)
		return data
	}

	//Check the lenght of the data
	if len(tableData.Data) < key || len(tableData.Data) == 0 {
		return data
	}

	//Return the data
	return tableData.Data[key]
}
