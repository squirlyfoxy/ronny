package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	//Get the text
	//redo:
	scanner := bufio.NewScanner(file)
	lines := []string{}
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	//Parse the text
	err = json.Unmarshal([]byte(strings.Join(lines, "\n")), &tableData)

	if err != nil {
		fmt.Println(err)
		return map[string]interface{}{}
	}

	//Build the return ads JSON (get columns name from table)
	var ret map[string]interface{}
	//Initialize the map
	ret = make(map[string]interface{})

	var data map[string]interface{} = make(map[string]interface{})
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

			for i, col := range table.Columns {
				data[col.Name] = v[i]
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

		//IF the column is a foreign key, get the data from the other table
		if column.IsExtern {
			if column.IsArray {
				var array []interface{} = make([]interface{}, 0)
				for _, v := range currentat.([]interface{}) {
					v_int, _ := strconv.Atoi(v.(string))
					array = append(array, GetDataFromATable(database, database.GetTable(column.TypeAsString), v_int))
				}
				ret[column.Name] = array
			} else {
				currentat_int, _ := strconv.Atoi(currentat.(string))
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
