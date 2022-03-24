package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func ExecuteFunction(
	table Table,
	function Function,
) (interface{}, error) {
	//Execute function and return result
	//Read the script of the table (and get from that only the rows of the function)

	var lines []string
	for x, tb := range database.Tables {
		exit := false

		if tb.Name == table.Name {
			for _, fn := range table.Functions {
				if fn.Name == function.Name {
					exit = true

					//Read the script of the table
					script_path := database.Scripts[x]
					content, err := ioutil.ReadFile(script_path)
					if err != nil {
						return nil, err
					}

					//Get only the rows of the function
					lines = strings.Split(string(content), "\n")
					lines = RemoveTabsFromLines(lines) //Remove tabs from lines

					lines = lines[fn.StartLine : fn.EndLine-1]
					break
				}
			}
		}

		if exit {
			break
		}
	}

	//To execute the function, we need to transform the script into a golang script, run it and return the result

	//************************
	//Transpiler

	we_are_in_a_loop := false

	types_to_transpile := []string{}
	types_list := []string{}
	transpiled_function := ""

	transpiled_script := "package main\n"
	transpiled_script += "import (\n"
	transpiled_script += "\t\"encoding/json\"\n"
	transpiled_script += ")\n"

	//packages_to_import := []string{} //Packages to import
	restart_index := 0
	for i, line := range lines {
		if i < restart_index {
			continue
		}

		//if line is @var, then we need to add the variable to the transpiled script
		if strings.HasPrefix(line, "@var") {
			type_to_transpile := ""

			//There are two whays to add a variable:
			// 1.
			// @var (type)
			// {
			//    Rows to select
			// } name;
			// 2.
			// @var (type) name;

			//We will create a type that contains all the rows to select
			//If the type is not a table, it will be a string, a float or a int

			var name string

			//split
			line_split := strings.Split(line, " ")
			if len(line_split) == 2 { //first way
				//Get the type
				type_ := line_split[1]
				type_ = AlphanumericOnly(type_)
				i++

				//Get the list of rows to select
				rows_to_select := []string{}
				for j := i + 1; j < len(lines); j++ {
					if strings.HasPrefix(lines[j], "}") {
						restart_index = j
						spl := strings.Split(lines[j], " ")
						name = AlphanumericOnly(spl[1])
						name = strings.Replace(name, ";", "", -1)
						break
					}

					lines[j] = AlphanumericOnly(lines[j])

					rows_to_select = append(rows_to_select, lines[j])
				}

				if !Contains(types_list, type_+"_"+name) {
					//Check if the type is a table
					is_table := false
					var types_to_add_to_type map[string]string = map[string]string{}
					for _, tb := range database.Tables {
						if tb.Name == type_ {
							for _, cols := range tb.Columns {
								//key: ColumnName, value: TypeAsString
								if Contains(rows_to_select, cols.Name) {
									types_to_add_to_type[cols.Name] = cols.TypeAsString
								}
							}

							is_table = true
							break
						}
					}
					if !is_table { //Is not a table but we are trying to manage it as a table, error
						return nil, fmt.Errorf("The type of the variable is not a table 2")
					}

					//Create the type
					type_to_transpile = "type " + type_ + "_" + name + " struct {\n"
					for col, type_ := range types_to_add_to_type {
						type_to_transpile += "    " + col + " " + type_ + "\n"
					}
					type_to_transpile += "}"

					//Add the type to the list of types to transpile
					types_to_transpile = append(types_to_transpile, type_to_transpile)
					types_list = append(types_list, type_+"_"+name)

					transpiled_script += type_to_transpile + "\n"
				}
			} else if len(line_split) == 3 { //second way
				//Get the type
				type_ := line_split[1]
				type_ = AlphanumericOnly(type_)
				name = line_split[2]
				name = strings.Replace(name, ";", "", -1)

				//Check if the type is a table
				is_table := false

				if !Contains(types_list, type_+"_"+name+"_all") {
					var types_to_add_to_type map[string]string = map[string]string{}
					for _, tb := range database.Tables {
						if tb.Name == type_ {
							for _, cols := range tb.Columns {
								//key: ColumnName, value: TypeAsString
								types_to_add_to_type[cols.Name] = cols.TypeAsString
							}

							is_table = true
							break
						}
					}
					if !is_table { //Is not a table but we are trying to manage it as a table, error
						return nil, fmt.Errorf("The type of the variable is not a table 1")
					}

					//Create the type
					type_to_transpile = "type " + type_ + " struct {\n"
					for col, type_ := range types_to_add_to_type {
						type_to_transpile += "    " + col + " " + type_ + "\n"
					}
					type_to_transpile += "}"

					//Add the type to the list of types to transpile
					types_to_transpile = append(types_to_transpile, type_to_transpile)
					types_list = append(types_list, type_+"_"+name+"_all")

					transpiled_script += type_to_transpile + "\n"
				}
			} else {
				return nil, fmt.Errorf("Error in the line")
			}

			//Add the variable to the transpiled function
			transpiled_function += "	var " + name + " []" + types_list[len(types_list)-1] + "\n"
			continue
		}

		//If the line starts with "loop", we need to get all the data of the table before (in the transpiler will be []map[string]interface{}))
		if strings.HasPrefix(line, "loop") {
			//Split
			line_split := strings.Split(line, " ")
			// loop on [table] as [rowName]
			if len(line_split) == 5 {
				//If [table] == *, we are looping in the current table
				line_split[4] = AlphanumericOnly(line_split[4])

				t_to_split := line_split[2]
				if t_to_split == "*" {
					t_to_split = table.Name
				}

				t_to_get := database.GetTable(t_to_split)

				//Get the table data
				data := GetAllData(database, t_to_get)
				data_json, err := json.Marshal(data)
				if err != nil {
					return nil, err
				}

				//data_json to a string
				data_json_str := string(data_json)

				//Add the data to the transpiled function as []map[string]interface{} (variable)
				transpiled_function += "	" + t_to_split + "_data_json" + " := []byte(`" + data_json_str + "`)\n"
				//Unmarshal t_to_split + "_data" to []map[string]interface{}
				transpiled_function += "	" + t_to_split + "_data := " + "make([]map[string]interface{}, 0)\n"
				transpiled_function += "	json.Unmarshal(" + t_to_split + "_data_json, &" + t_to_split + "_data)\n"

				//Prepare the loop
				transpiled_function += "	for _, " + line_split[4] + " := range " + t_to_split + "_data {\n"
				we_are_in_a_loop = true //WHAT IF WE HAVE A LOOP IN A LOOP?

				restart_index = i + 1

				continue
			} else {
				return nil, fmt.Errorf("Error in the line [", i, "]")
			}
		}

		//If starts with "if"
		if strings.HasPrefix(line, "if") {
			//write the if in the transpiled function
			split_line := strings.Split(line, " ")
			split_line[3] = AlphanumericOnly(split_line[3])
			transpiled_function += split_line[0] + " " + split_line[1] + split_line[2] + split_line[3] + " {\n"
			//We are in an if
			for j := i + 1; j < len(lines); j++ {
				line = lines[j]
				if strings.HasPrefix(line, "{") {
					continue
				}

				if strings.HasPrefix(line, "}") {
					//We are out of the if
					restart_index = j + 1
					break
				}

				//If variable.Add(value), append
				//Split by .
				if strings.Contains(line, ".Add") {
					line_split := strings.Split(line, ".")
					if len(line_split) == 2 {
						line_split[1] = strings.Replace(line_split[1], ";", "", -1)
						line_split[1] = strings.Replace(line_split[1], "Add(", "", -1)
						line_split[1] = strings.Replace(line_split[1], ")", "", -1)
						line_split[1] = AlphanumericOnly(line_split[1])

						//Append the value to the variable
						transpiled_function += "	" + line_split[0] + " = append(" + line_split[0] + ", " + line_split[1] + ")\n"
					}
				}

				//restart_index = j + 1
			}
			continue
		}

		//If starts with "return"
		if strings.HasPrefix(line, "return") {
			//Remove ;
			line = strings.Replace(line, ";", "", -1)

			//Split
			line_split := strings.Split(line, " ")

			transpiled_function += "	" + line_split[0] + " " + line_split[1] + "\n"
		}
	}

	if we_are_in_a_loop {
		transpiled_function += "	}\n"
		we_are_in_a_loop = false
	}

	transpiled_function += "}\n"
	transpiled_script += "func " + function.Name + "() {\n" + transpiled_function + "}\n"

	//Main function
	transpiled_script += "func main() {\n"
	transpiled_script += "	" + function.Name + "()\n"
	transpiled_script += "}\n"

	//Save the transpiled function
	ioutil.WriteFile("./transpiled_functions/"+function.Name+".go", []byte(transpiled_script), 0644)

	//TODO: COMPILE THE CREATED FILE

	//TODO: RUN THE COMPILED FILE

	//TODO: GET THE RESULT

	return nil, nil
}

func ToString(row map[string]interface{}) string {
	return fmt.Sprintf("%v", row)
}
