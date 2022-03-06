package database

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func GetDataFromATable(database Database, table Table, key int) map[string]interface{} {
	//Check if CAN_GLOBALLY_TAKE
	contains := false

	for _, ts := range table.Rule.RuleTypes {
		if ContainsORtype(ts.Can, CAN_GLOBALLY_TAKE) {
			contains = true
		}
	}
	if !contains {
		return map[string]interface{}{
			"error": "You can't take this data",
		}
	}

	//Read the JSON (./db/data/[tableName].dat.json)
	file, err := os.Open("./db/data/" + table.Name + ".dat.json")
	if err != nil {
		fmt.Println(err)
		return map[string]interface{}{
			"error": "Error while opening the file",
		}
	}
	defer file.Close()

	//Read the data (TableData)
	var tableData TableData
	err = json.NewDecoder(file).Decode(&tableData)
	if err != nil {
		fmt.Println(err)
		return map[string]interface{}{}
	}

	//Check if the key is valid
	if key < 0 || key >= len(tableData.Data) {
		return map[string]interface{}{
			"error": "The key is not valid",
		}
	}

	data := tableData.Data[key]

	//Build the return ads JSON (get columns name from table)
	var ret map[string]interface{}
	//Initialize the map
	ret = make(map[string]interface{})

	var i int = 0

	for _, column := range table.Columns {
		current := data[i]

		//Check if the columnt is int or float
		if column.Type == INT || column.Type == FLOAT || (column.Type == KEY && column.Rule == AUTOINCREMENT) {
			//Is a JSON, [tName]:[value]

			//Convert current in a int
			var currentInt int
			currentInt, err = strconv.Atoi(current)
			if err != nil {
				fmt.Println(err)
				return map[string]interface{}{}
			}

			ret[column.Name] = currentInt
		} else {
			ret[column.Name] = current
		}

		i++
	}

	return ret
}
