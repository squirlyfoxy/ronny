package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	dijkstra "github.com/squirlyfoxy/ronny/database/dijkstra"
)

//Saved as json in './db/database.json'
type Database struct {
	Name                      string          `json:"name"`
	Scripts                   []string        `json:"scripts"`           //Scripts path
	DatasFilesPaths           []string        `json:"datas_files_paths"` //Paths of the files containing the data (will be stored as jsons)
	Tables                    []Table         `json:"tables"`            //Tables
	WhereUSERACCESSKEYLocated string          `json:"-"`                 //Table name
	Config                    Config          `json:"-"`
	DijkstraRappresentation   *dijkstra.Graph `json:"-"`
}

func CreateTFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	file.Close()
}

func InitTFile(path string, t Table) {
	var t_data TableData
	t_data.Columns = t.Columns
	t_data.Data = make([][]interface{}, 0)

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	//Write the data
	json.NewEncoder(file).Encode(t_data)
}

func (d *Database) UpdateDB() {
	UnmarshallTablesIDsList()

	d.DijkstraRappresentation = dijkstra.NewGraph()

	for _, table := range d.Tables {
		//The script is parsed, now we need to update the DijkstraRappresentation (every node have weight of 1)
		d.DijkstraRappresentation.AddEdge(table.Name, "-", 1) //- is the root node
		//Subtables
		for _, subtable := range table.SubTables {
			d.DijkstraRappresentation.AddEdge(table.Name, subtable.Name, 1)
		}
		//1:N
		for _, ext_ty := range table.ExternalTypes {
			d.DijkstraRappresentation.AddEdge(table.Name, ext_ty.Type, 1)
		}
	}
}

func (d *Database) GetTable(tableName string) Table {
	for _, table := range d.Tables {
		if table.Name == tableName {
			return table
		}
	}

	return Table{}
}

func (d *Database) AddData(table Table, data []interface{}) {
	//Get the path of the file
	path := "./db/data/" + table.Name + ".dat.json"

	//Read the file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Unmarshal the data
	var t_data TableData
	err = json.Unmarshal(file, &t_data)
	if err != nil {
		fmt.Println(err)
		return
	}

	//data is an array of interface{}, we need to append something like that: "", "", []...
	var new_data []interface{}
	for i, d := range data {
		if d == nil {
			continue
		}

		//If is an array, the data to appenmd is ["el1", "el2", ...]"]
		switch d.(type) {
		case []interface{}:
			//Check if the column at i is an array, if not..
			if t_data.Columns[i].IsArray == false {
				//There is a problem...
				switch d.([]interface{})[len(d.([]interface{}))-1].(type) {
				case string:
					new_data = append(new_data, d.([]interface{})[len(d.([]interface{}))-1].(string))
					break
				case int:
					new_data = append(new_data, d.([]interface{})[len(d.([]interface{}))-1].(int))
					break
				case float64:
					new_data = append(new_data, d.([]interface{})[len(d.([]interface{}))-1].(float64))
					break
				case bool:
					new_data = append(new_data, d.([]interface{})[len(d.([]interface{}))-1].(bool))
					break
				}

				continue
			}

			//Take the last d.([]interface{})[len(d.([]interface{}))-1] that is a string and convert it to a []interface
			to_convert := d.([]interface{})[len(d.([]interface{}))-1]
			var to_append []interface{}

			//Remove from to_convert [ and ]
			to_convert = to_convert.(string)[1 : len(to_convert.(string))-1]
			//Split by ,
			ts := strings.Split(to_convert.(string), ",")

			for _, el := range ts {
				switch el {
				case "":
					continue
				case "null":
					to_append = append(to_append, nil)
					break
				case "true":
					to_append = append(to_append, true)
					break
				case "false":
					to_append = append(to_append, false)
					break
				default:
					to_append = append(to_append, el)
					break
				}
			}

			new_data = append(new_data, to_append)

			break
		case int:
			new_data = append(new_data, fmt.Sprintf("%d", d.(int)))
			break
		case float64:
			new_data = append(new_data, fmt.Sprintf("%f", d.(float64)))
			break
		case bool:
			new_data = append(new_data, d.(bool))
			break
		case string:
			new_data = append(new_data, d.(string))
			break
		}
	}
	//Append the new data
	t_data.Data = append(t_data.Data, new_data)

	//Marshal the data
	file, err = json.Marshal(t_data)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Write the data
	err = ioutil.WriteFile(path, file, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (d *Database) ReadScript(path string) {
	d.Scripts = append(d.Scripts, path)

	fmt.Println("Reading script: " + path)

	//Read a the file, get the lines and parse them
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	//Get the hash of the script
	hash, err := md5sum(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	if HashAlreadyContained(hash) {
		return
	}

	AddHashToDB(hash, path)

	var lines []string

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	table := Parse(lines)

	//Parse the lines
	d.Tables = append(d.Tables, table)

	//Update the DijkstraRappresentation
	d.UpdateDB()

	AddTableToIDsList(table.Name)

	//If "./db/data/" + d.Tables[len(d.Tables)-1].Name + ".dat.json" exists, skip it
	if _, err := os.Stat("./db/data/" + d.Tables[len(d.Tables)-1].Name + ".dat.json"); err == nil {
		return
	}

	//Create the file that will contains the data (./db/data/[TableName].dat.json)
	//Create the folder if it doesn't exist
	if _, err := os.Stat("./db/data"); os.IsNotExist(err) {
		os.Mkdir("./db/data", 0777)
	}

	//Create the file (check before if already exists, if so request user for permission to recreate)
	if _, err := os.Stat("./db/data/" + d.Tables[len(d.Tables)-1].Name + ".dat.json"); os.IsNotExist(err) {
		CreateTFile("./db/data/" + d.Tables[len(d.Tables)-1].Name + ".dat.json")

		InitTFile("./db/data/"+d.Tables[len(d.Tables)-1].Name+".dat.json", table)

		d.DatasFilesPaths = append(d.DatasFilesPaths, "./db/data/"+d.Tables[len(d.Tables)-1].Name+".dat.json")
	}
}

func ReadDatabase() (Database, error) {
	var database Database
	//Check if file exists, if not create it
	if _, err := ioutil.ReadFile("./db/database.json"); err != nil {
		ioutil.WriteFile("./db/database.json", []byte("{}"), 0644)
	}

	file, err := ioutil.ReadFile("./db/database.json")
	if err != nil {
		return database, err
	}
	err = json.Unmarshal(file, &database)

	for _, table := range database.Tables {
		located := false

		for _, column := range table.Columns {
			if column.Type == USERACCESSKEY {
				database.WhereUSERACCESSKEYLocated = table.Name
				located = true
				break
			}
		}

		if located {
			break
		}
	}

	//Check if scripts are in the database
	for _, scr := range database.Scripts {
		if _, err := os.Stat(scr); os.IsNotExist(err) {
			//Print error in red: Script [path] not found, remove it from the database will erase the data. (y/n)
		redo:
			fmt.Printf("\033[31mScript [%s] not found, remove it from the database will erase the data. (y/n) -> \033[0m", scr)
			var answer string

			fmt.Scanln(&answer)

			if answer == "y" {
				//Remove the script from the database
				database.Scripts = Remove(database.Scripts, scr)

				//Remove the data file
				database.DatasFilesPaths = Remove(database.DatasFilesPaths, "./db/data/"+database.Tables[len(database.Tables)-1].Name+".dat.json")
				err := os.Remove("./db/data/" + database.Tables[len(database.Tables)-1].Name + ".dat.json")
				if err != nil {
					fmt.Println(err)
				}

				//Print in green: Script [path] removed
				fmt.Println("")
				fmt.Printf("\033[32mScript [%s] removed\033[0m\n", scr)
			} else if err != nil {
				goto redo
			}
		}
	}

	//Check if data files are in the database
	for _, data := range database.DatasFilesPaths {
		if _, err := os.Stat(data); os.IsNotExist(err) {
			//Print error in red: Data file [path] not found.
			fmt.Printf("\033[31mData file [%s] not found.\033[0m\n", data)

			//Remove the data file
			database.DatasFilesPaths = Remove(database.DatasFilesPaths, data)
		}
	}

	database.UpdateDB()

	return database, nil
}

func (d *Database) AddScript(scr string) {
	//If this file is already in the database, skip it
	for _, path := range d.Scripts {
		if path == scr {
			return
		}
	}

	d.ReadScript(scr)

	d.Scripts = append(d.Scripts, scr)

	d.Save()
}

func (d *Database) Save() error {
	file, err := json.Marshal(d)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./db/database.json", file, 0644)
	if err != nil {
		return err
	}
	return nil
}
