package database

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var database *Database

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

	//Get the table
	for _, table := range database.Tables {
		if table.Name == tableName {
			//Get the data
			data := GetDataFromATable(*database, table, key_int)

			//Return the data
			c.JSON(200, gin.H{
				"data": data,
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "error",
		"why":     "table not found",
	})
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
		for _, table := range database.Tables {
			if table.Name == tableName { //Check if the table exists
				//Get the data
				data := GetDataFromATable(*database, table, key_int)

				//Return the data
				c.JSON(200, gin.H{
					"data": data,
				})
				return
			}
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

	//Start server
	r.Run(db.Config.Host + ":" + fmt.Sprintf("%d", db.Config.Port))
}
