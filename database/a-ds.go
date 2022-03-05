package database

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func AddRoute(table Table, engine *gin.Engine) {
	//Route starts with /api/v1/
	//Methods:
	// - take
	// - add
	// - remove
	// - modify

	//How? POST, pass this:
	//id (URL PARAMATER)=the id (take, remove or modify)
	//data (BODY)=the data to be added, modified or removed

	//In case of funtions:
	// /api/v1/function/tableName/functionName

	//TODO: Add routes for every table
}

func StartADS(db *Database) {
	//Gin server
	r := gin.Default()

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

	for _, table := range db.Tables {
		//For every table, get the rule.
		//If the rule is empty, then the table is not accessible
		if table.Rule.RefeersTo == "" {
			continue
		} else {
			AddRoute(table, r)
		}
	}

	//Start server
	r.Run(db.Config.Host + ":" + fmt.Sprintf("%d", db.Config.Port))
}
