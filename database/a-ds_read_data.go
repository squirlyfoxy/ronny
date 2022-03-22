package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Buffers
var primary_key_positions_buffer map[string]int = make(map[string]int)

func GetDataFromATable(database Database, table Table, key int) map[string]interface{} {
	//Check if CAN_GLOBALLY_TAKE
	if !CheckTableRule(table, CAN_GLOBALLY_TAKE) {
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

	//Read the data (TableData)
	var tableData TableData

	//Get the text
	//redo:
	scanner := bufio.NewScanner(file)
	lines := []string{}
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()

	//Parse the text
	err = json.Unmarshal([]byte(strings.Join(lines, "\n")), &tableData)
	if err != nil {
		fmt.Println(err)
		return map[string]interface{}{}
	}

	var ret map[string]interface{} = make(map[string]interface{})
	var data map[string]interface{} = make(map[string]interface{})

	var pr_key_pos int = 0
	var entered bool = false

	//Position of the primary key
	//Check if primary_key_positions_buffer contains the table name as a key
	if _, ok := primary_key_positions_buffer[table.Name]; ok {
		pr_key_pos = primary_key_positions_buffer[table.Name]
	} else {
		for _, col := range table.Columns {
			if col.Type == KEY && col.Rule == AUTOINCREMENT {
				break
			}
			pr_key_pos++
		}

		primary_key_positions_buffer[table.Name] = pr_key_pos
	}

	for _, v := range tableData.Data { //For each row
		if v[pr_key_pos] == strconv.Itoa(key) { //If the primary key is the same as the key we are looking for
			entered = true

			for i, v2 := range v {
				data[table.Columns[i].Name] = v2
			}

			break
		}
	}
	if !entered {
		return map[string]interface{}{
			"error": "No data found",
		}
	}

	for key, currentat := range data {
		var pos int = 0
		for _, col := range table.Columns {
			if col.Name == key {
				break
			}
			pos++
		}
		column := table.Columns[pos]

		if currentat == nil {
			ret[column.Name] = nil
			continue
		}

		//IF the column is a foreign key, get the data from the other table
		if column.IsExtern {
			//N:N
			if column.IsArray {
				var array []interface{} = make([]interface{}, 0)

				for _, v := range currentat.([]interface{}) {
					v_int, err := strconv.Atoi(v.(string))
					if err != nil {
						ret[column.Name] = nil
						continue
					}

					array = append(array, GetDataFromATable(database, database.GetTable(column.TypeAsString), v_int))
				}
				ret[column.Name] = array
			} else {
				//1:N
				currentat_int, err := strconv.Atoi(currentat.(string))
				if err != nil {
					ret[column.Name] = nil
					continue
				}
				ret[column.Name] = GetDataFromATable(database, database.GetTable(column.TypeAsString), currentat_int)
			}

			continue
		}

		if column.Type == INT || (column.Type == KEY && column.Rule == AUTOINCREMENT) {
			currentat_int, _ := strconv.Atoi(currentat.(string))
			ret[column.Name] = currentat_int

			continue
		} else if column.Type == FLOAT {
			currentat_float, _ := strconv.ParseFloat(currentat.(string), 64)
			ret[column.Name] = currentat_float

			continue
		}

		ret[column.Name] = currentat
	}

	return ret
}
