package main

import (
	db "github.com/squirlyfoxy/ronny/database"
)

func main() {
	//Read databases
	database, err := db.ReadDatabase()
	if err != nil {
		panic(err)
	}

	//Start command parser
	StartCMD(database)
}
