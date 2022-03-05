package database

import (
	"fmt"
	"strings"
)

func ParseRule(lines []string, start int) ([]RuleType, string, int) {
	var rules []RuleType
	var refeersTo string

	//Rule:
	//@rule
	//{
	//	@on_type
	//	{
	//		CAN_GLOBALLY_TAKE
	//		CAN_TAKE
	//		CAN_ADD
	//		CAN_REMOVE
	//		CAN_MODIFY
	//	}
	//} REFEERS TO [key (a columnt of the table that is "key" type)];

	//CAN_TAKE create a /api/v1/take/[tableName]/from/[startTable]/where/[key], if startTable = "-", no start table (POST)
	//CAN_ADD create a /api/v1/add/[tableName] (POST)
	//CAN_REMOVE create a /api/v1/remove/[tableName]/where/[key] (POST)
	//CAN_MODIFY create a /api/v1/modify/[tableName]/where/[key] (POST)
	//CAN_GLOBALLY_TAKE create a /api/v1/take/[tableName]/where/[key] (GET, NO KEY NEEDED)

	//Check if the line is a rule
	if !strings.HasPrefix(lines[start], "@rule") {
		return rules, refeersTo, start
	} else {
		if !strings.HasPrefix(lines[start+1], "{") {
			fmt.Println("Error: Rule has no body")
			return rules, refeersTo, start
		}

		start += 2
	}

	for ; start < len(lines); start++ {
		if strings.HasPrefix(lines[start], "@on_type") {
			start += 1
			if !strings.HasPrefix(lines[start], "{") {
				fmt.Println("Error: Rule has no body")
				return rules, refeersTo, start
			}

			start += 1

			var r RuleType

			for ; start < len(lines); start++ {
				if strings.HasPrefix(lines[start], "}") {
					ts := strings.TrimPrefix(lines[start], "}")
					ts = strings.TrimSpace(ts)
					ts = strings.TrimSuffix(ts, ";")

					if len(ts) > 0 {
						r.TableName = ts
					}

					start++
					break
				}

				if strings.HasPrefix(lines[start], "CAN_TAKE") {
					r.Can = append(r.Can, CAN_TAKE)
					continue
				} else if strings.HasPrefix(lines[start], "CAN_ADD") {
					r.Can = append(r.Can, CAN_ADD)
					continue
				} else if strings.HasPrefix(lines[start], "CAN_REMOVE") {
					r.Can = append(r.Can, CAN_REMOVE)
					continue
				} else if strings.HasPrefix(lines[start], "CAN_MODIFY") {
					r.Can = append(r.Can, CAN_MODIFY)
					continue
				} else if strings.HasPrefix(lines[start], "CAN_GLOBALLY_TAKE") {
					r.Can = append(r.Can, CAN_GLOBALLY_TAKE)
					continue
				} else {
					fmt.Println("Error: Unknown on_type")
				}

				fmt.Println(r)
			}
			rules = append(rules, r)

			continue
		}

		if strings.HasPrefix(lines[start-1], "}") {
			s := strings.TrimPrefix(lines[start-1], "}")
			if s == "" {
				fmt.Println("Error: Rule has no end")
			} else {
				//Remove "REFEERS TO"
				s = strings.TrimPrefix(s, " REFEERS TO ")

				//Remove ";"
				s = strings.TrimSuffix(s, ";")

				refeersTo = s
			}

			start++
			break
		}
	}
	return rules, refeersTo, start
}
