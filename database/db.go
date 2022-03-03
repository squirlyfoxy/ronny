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

	var lines []string

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	//Parse the lines
	d.Tables = append(d.Tables, Parse(lines))

	//Create the file that will contains the data (./db/data/[TableName].dat.json)
	//Create the folder if it doesn't exist
	if _, err := os.Stat("./db/data"); os.IsNotExist(err) {
		os.Mkdir("./db/data", 0777)
	}

	//Create the file (check before if already exists, if so request user for permission to recreate)
	if _, err := os.Stat("./db/data/" + d.Tables[len(d.Tables)-1].Name + ".dat.json"); os.IsNotExist(err) {
		file, err := os.Create("./db/data/" + d.Tables[len(d.Tables)-1].Name + ".dat.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		//Prepare the file with a basic structure
		file.WriteString("{\"columns\":[],\"data\":[]}")

		file.Close()
	} else {
	redo:
		fmt.Printf("Do you want to recreate the data stored in Ronny for this table? No can be dangerous if you have changed the data structure (y/n) -> ")
		var answer string

		fmt.Scanln(&answer)

		if answer == "y" {
			file, err := os.Create("./db/data/" + d.Tables[len(d.Tables)-1].Name + ".dat.json")
			if err != nil {
				fmt.Println(err)
				return
			}
			//Prepare the file with a basic structure
			file.WriteString("{\"columns\":[],\"data\":[]}")

			file.Close()
		} else if answer != "n" {
			goto redo
		}
	}

	//Save the database
	d.Save()
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
