package main

import (
	"os"

	db "github.com/squirlyfoxy/ronny/database"
)

func main() {
	//Read databases
	database, err := db.ReadDatabase()
	if err != nil {
		panic(err)
	}

	//Read the configuration file
	database.Config = db.ReadConfiguration()

	//Check if param is -serve
	if len(os.Args) == 2 && os.Args[1] == "-serve" {
		//Start server
		db.StartADS(&database)
	} else {
		//Start command parser
		StartCMD(database)
	}
}
