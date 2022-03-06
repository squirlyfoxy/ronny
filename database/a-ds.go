package database

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
		name_to_lower := strings.ToLower(table.Name)

		if name_to_lower == tableName {
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
	//tableName := c.Param("tableName")

	//Get the start table
	//startTable := c.Param("startTable")

	//Get the key
	//key := c.Param("key")
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
