package database

import (
	"fmt"
	"strings"
)

func GetColumnType(spl []string, col_n string, table *Table) ColumnType {
	//Get the column type

	//If the column type is a string
	if strings.HasPrefix(spl[1], "string") {
		return STRING
	}

	//If the column type is an int
	if strings.HasPrefix(spl[1], "int") {
		return INT
	}

	//If the column type is a float
	if strings.HasPrefix(spl[1], "float") {
		return FLOAT
	}

	//If the column type is a bool
	if strings.HasPrefix(spl[1], "bool") {
		return BOOL
	}

	//If the column type is a date
	if strings.HasPrefix(spl[1], "date") {
		return DATE
	}

	//If the column type is a key
	if strings.HasPrefix(spl[1], "key") {
		return KEY
	}

	var ts ExternalType
	ts.Type = spl[1]
	ts.ColumnName = col_n

	table.ExternalTypes = append(table.ExternalTypes, ts)

	return EXTERNAL
}

func GetColumnRule(spl []string, col_n string, table *Table) ColumnRule {
	//AUTOINCREMENT
	//UNIQUE
	//NOT_NULL
	//USERACCESSKEY

	//if spl[2] not set, return the default value
	if len(spl) < 3 {
		return DEFAULT
	}

	//if AUTOINCREMENT
	if strings.HasPrefix(spl[2], "AUTOINCREMENT") {
		//TODO: STARTS FROM?

		return AUTOINCREMENT
	}

	//if UNIQUE
	if strings.HasPrefix(spl[2], "UNIQUE") {
		return UNIQUE
	}

	//if NOT_NULL
	if strings.HasPrefix(spl[2], "NOT_NULL") {
		return NOT_NULL
	}

	//if USERACCESSKEY
	if strings.HasPrefix(spl[2], "USERACCESSKEY") {
		return USERACCESSKEY
	}

	//Not found
	fmt.Println("Rule not found")
	return DEFAULT
}

func ParseColumns(lines []string, table *Table, start int) ([]Column, int) {
	//Will return the columns and the index of the next line to be parsed

	//The columns will be created in the following way:
	//[columnName] [columnType]	[columnRule]

	var cols []Column
	var i = start

	lines = RemoveTabsFromLines(lines)

	for ; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "//") {
			continue
		}

		if lines[i] == "" || len(lines[i]) == 0 || strings.HasPrefix(lines[i], "{") || strings.HasPrefix(lines[i], "}") {
			break
		}

		//If the line starts with "@", "{" or "}", the columns creation process is finished
		if strings.HasPrefix(lines[i], "@") {
			i--
			break
		}

		var col Column

		//From the current line create an array that in each position will contains only one word
		//Only characters will remains
		spl := strings.FieldsFunc(lines[i], func(r rune) bool {
			return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || r == '[' || r == ']')
		})

		//If len(sql) = 0, split
		if len(spl) == 0 {
			break
		}

		//Get the column name
		col.Name = spl[0]
		//Get the column type
		if strings.HasPrefix(spl[1], "[]") {
			col.IsArray = true
			spl[1] = spl[1][2:]
		}

		col.Type = GetColumnType(spl, col.Name, table)
		//If col.Type is EXTERNAL
		if col.Type == EXTERNAL {
			col.IsExtern = true
		}

		//set the column type as string
		col.TypeAsString = spl[1]

		//Get the column rule
		col.Rule = GetColumnRule(spl, col.Name, table)

		//Check the name
		if col.Name == "" {
			fmt.Println("Column name is empty")
		}

		cols = append(cols, col)
	}

	return cols, (i - 1)
}

func Parse(lines []string) Table {
	//Every script can contains only a table (called type)

	//Structure of a script:
	//@type TableName
	//{
	//	/*Table Columns*/
	// 	/*Table Functions*/
	//	/*Table Rules*/
	//}

	//A table can contains subtables (subtypes) that can be assigned as a type of a row. Is the same as a table but with a different name.
	//How a columns is created:
	//...
	//[columnName] [columnType]	[columnRule]
	//...

	//See every ron file in "scripts" folder for more information

	//Parse the script

	var table Table

	//Loop through the lines
	columns_creation_process_finished := false

	//Remove tabs
	lines = RemoveTabsFromLines(lines)
	for i := 0; i < len(lines); i++ {
		//If line starts with "//" is a comment
		if strings.HasPrefix(lines[i], "//") || lines[i] == "" || lines[0] == "" {
			continue
		}

		//Get the father table name (if not yet set), if set (so the line starts with "@type " but the table is under construction, create a new subtable)
		if strings.HasPrefix(lines[i], "@type ") {
			//Get the name (after @type )

			if table.Name == "" {
				table.Name = strings.TrimPrefix(lines[i], "@type ")
			} else {
				//Create a new subtable
				table.SubTables = append(table.SubTables, Parse(lines[i:]))
			}
			//Check if in the following line there is only "{" character
			if lines[i+1] != "{" {
				fmt.Println("Expected '{' character")
				return Table{}
			}

			//Skip the next line
			i++
			continue
		}

		//If the line is a column so:
		//Not starts with "@"
		//Not starts with "{"
		//Not starts with "}".
		//This thing is made only one time, when a new column is found but the columns creation process is finisched, error.
		if !strings.HasPrefix(lines[i], "@") && !strings.HasPrefix(lines[i], "{") && !strings.HasPrefix(lines[i], "}") {
			if columns_creation_process_finished {
				//fmt.Println("Error: Columns creation process is already finished, but a new column was found.")
				//fmt.Println(lines[i])
				continue
			}

			//Call a method that will parse all the columns
			table.Columns, i = ParseColumns(lines, &table, i)
			columns_creation_process_finished = true

			continue
		}

		//Functions
		//If the line starts with @function, create a new function (@function (NameOfTheFunction))
		if strings.HasPrefix(lines[i], "@function") {
			var t Function
			t, i = ParseFunction(lines, i)

			table.Functions = append(table.Functions, t)

			continue
		}

		//Rules
		//If the line starts with @rule, create a new rule (@rule), the rule refeer to a table (REFEERS TO [NameOfTheTable] after "}")
		if strings.HasPrefix(lines[i], "@rule") {
			var t Rule
			r_types, refeers, new_i := ParseRule(lines, i)
			t.RefeersTo = refeers
			t.RuleTypes = r_types
			i = new_i

			table.Rule = t
			continue
		}
	}

	return table
}
