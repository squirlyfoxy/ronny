package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	db "github.com/squirlyfoxy/ronny/database"
)

func StartCMD(dt db.Database) {
	//Start command parser
	fmt.Println("Welcome to Ronny Database Version 1.0")
	fmt.Println("Type 'help' for a list of commands")
	fmt.Println("--------------------------------------")
	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("ronny> ")

		//Read input, space is not a newline
		var input string
		input, _ = reader.ReadString('\n')
		//CRLF to LF
		input = strings.Replace(input, "\n", "", -1)

		//Commands
		//'help' - Prints the list of commands
		//'exit' - Exits the program
		//'take [typeName] from [startTypeName] where [startTypeKey] = [value]' - execute a simple query
		//'exec [functionName] from [typeName]' - execute a function
		//'reload' - reloads the database

		//Check if input is 'help'
		if strings.HasPrefix(input, "help") {
			fmt.Println("'help' - Prints the list of commands")
			fmt.Println("'exit' - Exits from Ronny")
			fmt.Println("'wami' - Show the name of the database")
			fmt.Println("'chme [newName]' - Change the name of the database")
			fmt.Println("'tlist' - show all the tables")
			fmt.Println("'shtable [tableID]' - show a table structure")
			fmt.Println("'save' - save the database")
			fmt.Println("'serve' - start the API web server")
			fmt.Println("'realod scripts' - reload the database resetting the data structures")
		} else if strings.HasPrefix(input, "exit") {
			fmt.Println("Bye!")
			break
		} else if strings.HasPrefix(input, "reload scripts") {
			//Remove all the scripts
			dt.Scripts = []string{}
			//Remove all the tables
			dt.Tables = []db.Table{}

			//Get all the files in the './db/scripts' folder
			files, err := ioutil.ReadDir("./db/scripts")
			if err != nil {
				fmt.Println(err)
				continue
			}

			//For each file
			for _, f := range files {
				dt.ReadScript(string("./db/scripts/" + f.Name()))
			}

			fmt.Println()

			//Write "Database reloaded!" in green
			fmt.Print("\033[32mDatabase reloaded!\033[0m\n")
		} else if strings.HasPrefix(input, "tlist") {
			if len(dt.Tables) == 0 {
				fmt.Println("No tables found")
				continue
			}

			//Print all the table names
			fmt.Println("Tables:")
			fmt.Println("--------------------------------------")
			fmt.Println("|  ID          |  Name               |")
			fmt.Println("--------------------------------------")
			var i int
			for _, t := range dt.Tables {
				fmt.Printf("|  %d           |  %s             |\n", i, t.Name)
				i++
			}
			fmt.Println("--------------------------------------")
			fmt.Println("")
			fmt.Println("Use ID in this environment to access tables fastly.")
		} else if strings.HasPrefix(input, "shtable") {
			//Get the table ID (shtable [tableID]). convert it to an int. Remove \r at the end
			rs := strings.Replace(input, "shtable ", "", -1)
			rs = strings.Replace(rs, "\r", "", -1)
			tableID, err := strconv.Atoi(rs)
			if err != nil {
				fmt.Println("Invalid table ID, use tlist to see all the avviable tables")
				continue
			}

			if len(dt.Tables) <= tableID {
				fmt.Println("Invalid table ID, use tlist to see all the avviable tables")
				continue
			}

			if len(dt.Tables) == 0 {
				fmt.Println("No tables found")
				continue
			}

			//Get the table
			t := dt.Tables[tableID]

			//Print the table structure
			fmt.Println("Table: " + t.Name)

			//Print the table structure
			//Name and type should be alighed to the same length
			fmt.Printf("%-15s\n", "--------------------------------------")
			fmt.Printf("|  %-15s |  %-15s |\n", "Name", "Type")
			fmt.Printf("%-15s\n", "--------------------------------------")
			for _, c := range t.Columns {
				var ts string
				ts = db.ColumnTypes[c.Type]

				fmt.Printf("|  %-15s |  %-15s |", c.Name, ts)

				//If autoincrement, is the primary key ("<-")

				if c.Rule == db.AUTOINCREMENT {
					fmt.Printf("  <- PRIMARY_KEY")
				} else if c.Rule == db.USERACCESSKEY {
					fmt.Printf("  <- USER_ACCESS_KEY")
				}
				fmt.Printf("\n")
			}
			fmt.Printf("%-15s\n", "--------------------------------------")
			fmt.Println("")
		} else if strings.HasPrefix(input, "wami") {
			if dt.Name == "" {
				fmt.Println("Your database has no name YeT.")
				continue
			}

			fmt.Println("Database: " + dt.Name)
		} else if strings.HasPrefix(input, "chme") {
			//Get the new name
			rs := strings.Replace(input, "chme ", "", -1)
			rs = strings.Replace(rs, "\r", "", -1)

			if rs == "" {
				fmt.Println("You need to specify a new name for the database")
				continue
			}

			//Change the name
			dt.Name = rs
		} else if strings.HasPrefix(input, "save") {
			//Save the database
			dt.Save()

			fmt.Println("Database saved!")
		} else if strings.HasPrefix(input, "serve") {
			//TODO: Start the API web server
		} else {
			fmt.Println("No commands with this path")
		}
	}
}
