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

	//Build the return ads JSON (get columns name from table)
	var ret map[string]interface{}
	//Initialize the map
	ret = make(map[string]interface{})

	var i int = 0
	var data []string
	var pr_key_pos int = 0
	var entered bool = false

	for _, col := range table.Columns {
		if col.Type == KEY && col.Rule == AUTOINCREMENT {
			break
		}
		pr_key_pos++
	}

	for _, v := range tableData.Data {
		if v[pr_key_pos] == strconv.Itoa(key) {
			entered = true
			data = v
			break
		}
	}

	if !entered {
		return map[string]interface{}{
			"error": "No data found",
		}
	}

	for _, column := range table.Columns {
		current := data[i]

		//Check if the columnt is int or float
		if column.Type == INT || column.Type == FLOAT || (column.Type == KEY && column.Rule == AUTOINCREMENT) {
			//Convert current in a int
			var currentInt int
			currentInt, err = strconv.Atoi(current)
			if err != nil {
				fmt.Println(err)
				return map[string]interface{}{}
			}

			ret[column.Name] = currentInt
		} else {
			entered := false

			for _, ext_type := range table.ExternalTypes {
				if ext_type.ColumnName == column.Name {
					//Get the data from the table "ext_type.Type" where the key is "current"
					var ext_data map[string]interface{}

					for _, tables := range database.Tables {
						if tables.Name == ext_type.Type {
							current_int, err := strconv.Atoi(current)
							if err != nil {
								fmt.Println(err)
								return map[string]interface{}{}
							}

							ext_data = GetDataFromATable(database, tables, current_int)

							ret[column.Name] = ext_data
							entered = true

							break
						}
					}

					break
				}
			}

			if entered {
				i++
				continue
			}

			ret[column.Name] = current
		}

		i++
	}

	return ret
}
