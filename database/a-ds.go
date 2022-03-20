package database

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var database *Database

func GetTable(tableName string, key int) (map[string]interface{}, Table) {
	//Get the table
	for _, table := range database.Tables {
		if table.Name == tableName {
			//Get the data
			data := GetDataFromATable(*database, table, key)

			return data, table
		}
	}

	return nil, Table{}
}

func CheckTableRule(table Table, rule OnType) bool {
	//Check if the rule is on the table
	for _, r := range table.Rule.RuleTypes {
		for _, r2 := range r.Can {
			if r2 == rule {
				return true
			}
		}
	}

	return false
}

func AddRoute(c *gin.Context) {
	//Set json
	c.Header("Content-Type", "application/json")

	///api/v1/insert/[tableName]

	//Table name
	var data []interface{}
	var tRef Table
	tableName := c.Param("tableName")
	found := false

	for _, r := range database.Tables {
		if r.Name == tableName {
			found = true
			tRef = r

			if !CheckTableRule(r, CAN_ADD) {
				c.JSON(http.StatusForbidden, gin.H{
					"message": "you can't add this",
				})
				return
			}
		}
	}

	if !found {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "table not found",
		})
		return
	}

	data = make([]interface{}, len(tRef.Columns))

	//Get the data
	for i, r := range tRef.Columns {
		if (r.Type == KEY && r.Rule == AUTOINCREMENT) || (r.Type == KEY && r.Rule == USERACCESSKEY) {
			//If ID
			if r.Rule == AUTOINCREMENT {
				key := tRef.GetLatestAutoIncrementedKey()
				if key == -1 {
					c.JSON(http.StatusForbidden, gin.H{
						"message": "undefined error",
					})
					return
				}

				data[i] = strconv.Itoa(key)
			} else if r.Rule == USERACCESSKEY {
				data[i] = GenerateUserAccessKey()
			}

			continue
		}

		data[i] = append(data, c.PostForm(r.Name))

		if data[i] == "" && r.Rule != NOT_NULL {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing data in the " + r.Name + " column",
			})
			return
		}
	}

	//Add the data
	database.AddData(tRef, data)

	c.JSON(200, gin.H{
		"message": "data added successfully",
	})
}

func CanGloballyTakeRoute(c *gin.Context) {
	//Set json
	c.Header("Content-Type", "application/json")

	///api/v1/take/[tableName]/where/[key]
	//Get the table name
	tableName := c.Param("tableName")

	//Get the key
	key := c.Param("key")

	key_int, err := strconv.Atoi(key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "key is not a number",
		})
		return
	}

	data, t := GetTable(tableName, key_int)
	if !CheckTableRule(t, CAN_GLOBALLY_TAKE) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "you can't take this",
		})
		return
	}

	if data == nil {
		c.JSON(200, gin.H{
			"message": "error",
			"why":     "table not found",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"data": data,
		})
		return
	}
}

func CanTakeRoute(c *gin.Context) {
	//Set json
	c.Header("Content-Type", "application/json")

	///api/v1/take/[tableName]/from/[startTable]/where/[key]
	//Get the table name
	tableName := c.Param("tableName")

	//Get the start table
	startTable := c.Param("startTable")

	//Get the key
	key := c.Param("key")
	key_int, err := strconv.Atoi(key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "key is not a number",
		})
		return
	}

	if tableName == "-" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error",
		})
		return
	}

	//We need to determinate the path to the table (from the start)
	//If the path pass the root, error (no correlation between the two tables)

	//Get the table
	cost, path := database.DijkstraRappresentation.GetPath(startTable, tableName)

	if cost == -1 {
		c.JSON(200, gin.H{
			"message": "error",
			"why":     "no path found",
		})
		return
	}

	if startTable == "-" || cost == 0 { //Starts from the root
		//Get the table
		data, t := GetTable(tableName, key_int)
		if !CheckTableRule(t, CAN_TAKE) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "you can't take this",
			})
			return
		}

		if data != nil {
			c.JSON(200, gin.H{
				"data": data,
			})
			return
		}
	} else {
		if cost > 0 {
			//Check if in the path there is a "-", if so, error (no correlation between the two)
			for _, node := range path {
				if node == "-" {
					c.JSON(200, gin.H{
						"message": "error",
						"why":     "no correlation between the two types",
					})
					return
				}
			}

			//We need to follow the path (path) to get the data
			for _, table := range database.Tables {
				if table.Name == startTable { //Check if the table exists
					if !CheckTableRule(table, CAN_TAKE) {
						c.JSON(http.StatusForbidden, gin.H{
							"message": "you can't take this",
						})
						return
					}

					//Get the data
					data := GetDataFromATable(*database, table, key_int)

					var ret interface{} = data[path[1]]

					//Loop through the path
					for i := 2; i < len(path); i++ {
						ret = ret.(map[string]interface{})[path[i]]
					}

					//Return the data
					c.JSON(200, gin.H{
						"data": ret,
					})

					return
				}
			}

		}

		return
	}

	c.JSON(200, gin.H{
		"message": "error",
		"why":     "table not found",
	})
}

func StartADS(db *Database) {
	//Gin server
	r := gin.Default()

	database = db

	//Routes
	r.GET("/api/v1/", func(c *gin.Context) {
		if db.Config.DoDefaultRoute == 1 {
			c.JSON(200, gin.H{
				"ronny_version": "1.0",
				"database_name": db.Name,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "hi",
			})
		}
	})

	//api/v1/take/[tableName]/where/[key]
	r.GET("/api/v1/take/:tableName/where/:key", CanGloballyTakeRoute)
	///api/v1/take/[tableName]/from/[startTable]/where/[key]
	r.POST("/api/v1/take/:tableName/from/:startTable/where/:key", CanTakeRoute)
	///api/v1/insert/[tableName]
	//POST:
	//		data of the table, [key-value]
	r.POST("/api/v1/insert/:tableName", AddRoute)

	//Start server
	r.Run(db.Config.Host + ":" + fmt.Sprintf("%d", db.Config.Port))
}
