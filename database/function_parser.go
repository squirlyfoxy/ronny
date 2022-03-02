package database

import (
	"fmt"
	"strings"
)

func ParseFunction(lines []string, current_line int) (Function, int) {
	//Function name: @function (FunctionName)

	//@function (FunctionName)
	// {
	// 	//Function body
	//  //Here the row could be:
	//  //		- Variables (@var (Type) Name; or :
	//  //			@var (Type)
	//  //			{
	//  //			     //Columns that i want in this type
	//	//			} Name;
	//  //		- Loops (loop on [TableName] as [RowName])
	//  //          //If TableName = *, loop will be on the table where the function is called
	//  //			//If TableName = [TableName], loop will be on the table with the name [TableName]
	//  //			//RowName is the name that will be used to access the row
	//  //      - If statements (if [condition])
	//  //          //If condition is true, the area delimitated by { and } will be executed

	// In a function, can exists variables in loops and if statements.
	// If statements and loops are delimited by { and }
	// }

	var variables []string
	var func_name string
	start_line := current_line

r:
	if strings.HasPrefix(lines[current_line], "\t") {
		lines[current_line] = lines[current_line][1:]
		goto r
	}

	//Get the function name
	if strings.HasPrefix(lines[current_line], "@function") {
		func_name = strings.TrimPrefix(lines[current_line], "@function ")
		//From func_name remove "(" and ")"
		func_name = strings.TrimPrefix(func_name, "(")
		func_name = strings.TrimSuffix(func_name, ")")
		current_line++

	rs:
		if strings.HasPrefix(lines[current_line], "\t") {
			lines[current_line] = lines[current_line][1:]
			goto rs
		}

		if strings.HasPrefix(lines[current_line], "{") {
			current_line++
		} else {
			fmt.Println("Error: Function " + func_name + " has no body")
			return Function{}, current_line
		}
	}

	for ; current_line < len(lines); current_line++ {
	redo:
		if strings.HasPrefix(lines[current_line], "\t") {
			lines[current_line] = lines[current_line][1:]
			goto redo
		}

		//The goal is to find the end of the function
		if lines[current_line] == "" || strings.HasPrefix(lines[current_line], "//") || len(lines[current_line]) == 0 {
			continue
		}

		//If var, check if is declared as specified top
		if strings.HasPrefix(lines[current_line], "@var") {
			//Check if the variable is declared as specified top (@var (Type) Name; or :
			//  //			@var (Type)
			//  //			{
			//  //			     //Columns that i want in this type
			//	//			} Name;

			//If after the var there is " (" and ")", and in this delimitated area there is something and, after that, there is the name that we want to assign at the variable with a ";", skip

			//Split the line by spaces
			spl := strings.FieldsFunc(lines[current_line], func(r rune) bool {
				return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || r == '(' || r == ')')
			})

			if len(spl) == 3 {
				//Check if the variable is declared as specified: @var (Type) Name;
				if strings.HasSuffix(lines[current_line], ";") {
					//Check if a type is specified
					if strings.HasPrefix(spl[1], "(") && strings.HasSuffix(spl[1], ")") {
						//Check if the type is specified
						if spl[1][1:len(spl[1])-1] != "" {
							//Check if the name is specified
							if spl[2] != "" {
								//Add the variable to the list
								variables = append(variables, spl[2])
								continue
							} else {
								fmt.Println("Name for @var nor specified")
							}
						} else {
							fmt.Println("Type for @var nor specified")
						}
					} else {
						fmt.Println("Type for @var nor specified")
					}
				} else {
					fmt.Println("Line should end with ;")
				}

				continue
			} else {
				if len(spl) == 2 {
					//  			@var (Type)
					//  			{
					//  			     //Columns that i want in this type
					//				} Name;
					if strings.HasPrefix(spl[1], "(") && strings.HasSuffix(spl[1], ")") {
						//Check if the type is specified
						if spl[1][1:len(spl[1])-1] != "" {
							//Check if the name is specified
							if spl[1] != "" {
								//Add the variable to the list
								variables = append(variables, spl[1])
								continue
							} else {
								fmt.Println("Name for @var nor specified")
							}
						} else {
							fmt.Println("Type for @var nor specified")
						}
					} else {
						fmt.Println("Type for @var nor specified")
					}

					//Check if in the following lines there is a { and }
					for ; current_line < len(lines); current_line++ {
						if strings.HasPrefix(lines[current_line], "{") && strings.HasSuffix(lines[current_line], "}") {
							break
						}
					}

					continue
				}
			}
		}

		//If loop, check if is declared as specified top
		if strings.HasPrefix(lines[current_line], "loop") {
			//A loop not ends with ";"
			//Check if the loop is declared as specified top (loop on [TableName] as [RowName])
			//  If TableName = *, loop will be on the table where the function is called
			//  If TableName = [TableName], loop will be on the table with the name [TableName]
			//  RowName is the name that will be used to access the row
			//  If the loop is declared as specified top, the area delimitated by { and } will be executed

			//Split the line by spaces
			spl := strings.FieldsFunc(lines[current_line], func(r rune) bool {
				return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || r == '(' || r == ')' || r == '*')
			})

			if len(spl) == 5 {
				//Check if the table is specified
				if spl[2] != "" {
					//Check if the row is specified
					if spl[4] != "" {
						//Check if the following lines there is a { (current_line+1) and } (end)
					rts:
						if strings.HasPrefix(lines[current_line+1], "\t") {
							lines[current_line+1] = lines[current_line+1][1:]
							goto rts
						}

						if strings.HasPrefix(lines[current_line+1], "{") {
							current_line++
							for ; current_line < len(lines); current_line++ {
							rtrs:
								if strings.HasPrefix(lines[current_line], "\t") {
									lines[current_line] = lines[current_line][1:]
									goto rtrs
								}
								if strings.HasPrefix(lines[current_line], "}") {
									break
								}
							}

							continue
						} else {
							fmt.Println("Loop not starts with }")
						}

						continue
					} else {
						fmt.Println("Row for loop not specified")
					}
				} else {
					fmt.Println("Table for loop not specified")
				}
			} else {
				fmt.Println("Loop not specified")
			}

			continue
		}

		//If if, check if is declared as specified top
		if strings.HasPrefix(lines[current_line], "if") {
			//Check if the if is declared as specified top (if [Condition])
			//  If Condition is true, the area delimitated by { and } will be executed
			//  If Condition is false, the area delimitated by { and } will not be executed

			//Split the line by spaces
			spl := strings.FieldsFunc(lines[current_line], func(r rune) bool {
				return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || r == '(' || r == ')')
			})

			if len(spl) == 3 {
				//Check if the condition is specified
				if spl[2] != "" {
					//Check if the following lines there is a { (current_line+1) and } (end)
					if strings.HasPrefix(lines[current_line+1], "{") {
						current_line++
						for ; current_line < len(lines); current_line++ {
						rtvs:
							if strings.HasPrefix(lines[current_line], "\t") {
								lines[current_line] = lines[current_line][1:]
								goto rtvs
							}

							if strings.HasPrefix(lines[current_line], "}") {
								break
							}
						}

						continue
					} else {
						fmt.Println("If not ends with }")
					}

					continue
				} else {
					fmt.Println("Condition for if not specified")
				}
			} else {
				fmt.Println("If not specified")
			}

			continue
		}

		//Check if the function has a return statement
		if strings.HasPrefix(lines[len(lines)-1], "return") {
			//return ends with the name of the variable returned and ;
			//Split the line by spaces
			spl := strings.FieldsFunc(lines[len(lines)-1], func(r rune) bool {
				return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || r == '(' || r == ')')
			})

			if len(spl) == 2 {
				//Check if the variable is specified
				if spl[1] != "" {
					//Check if the variable is declared
					//Remove from spl[1] the ";"
					if strings.HasSuffix(spl[1], ";") {
						spl[1] = spl[1][:len(spl[1])-1]
						if Contains(variables, spl[1]) {
							//Check if in the following lines there is a }
							if strings.HasPrefix(lines[current_line+1], "}") {
								break
							} else {
								fmt.Println("Function not ends with }")
							}
						} else {
							fmt.Println("Variable not declared")
						}
					} else {
						fmt.Println("Return not ends with ;")
					}
				}
			} else {
				fmt.Println("Return not specified")
			}
		}
	}

	return Function{
		Name:      func_name,
		StartLine: start_line,
	}, current_line
}
