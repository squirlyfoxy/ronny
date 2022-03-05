package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Saved as json in './db/database.json'
type Database struct {
	Name            string   `json:"name"`
	Scripts         []string `json:"scripts"`           //Scripts path
	DatasFilesPaths []string `json:"datas_files_paths"` //Paths of the files containing the data (will be stored as jsons)
	Tables          []Table  `json:"tables"`            //Tables
	Config          Config   `json:"-"`
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
	t_data.Data = [][]string{}

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	//Write the data
	json.NewEncoder(file).Encode(t_data)
}

func (d *Database) ReadScript(path string) {
	//If this file is already in the database, skip it
	for _, scr := range d.Scripts {
		if scr == path {
			return
		}
	}

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
	}
	d.DatasFilesPaths = append(d.DatasFilesPaths, "./db/data/"+d.Tables[len(d.Tables)-1].Name+".dat.json")
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

	//If data file in the database repeats, remove the repeated ones
	for i := 0; i < len(database.DatasFilesPaths); i++ {
		for j := i + 1; j < len(database.DatasFilesPaths); j++ {
			if database.DatasFilesPaths[i] == database.DatasFilesPaths[j] {
				database.DatasFilesPaths = Remove(database.DatasFilesPaths, database.DatasFilesPaths[j])
			}
		}
	}

	return database, nil
}

func (d *Database) AddScript(scr string) {
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
