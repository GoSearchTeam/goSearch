package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type listResponseJSON struct {
	indexName string   `json:indexName`
	values    []string `json:indexValues`
}

func HandleIndexRoutes(r *gin.Engine, app *appIndexes) {
	r.GET("/index/listItems", func(c *gin.Context) {
		list := app.listIndexItems()
		// stringed := make([]string, 0)
		// for _, i := range list {
		// 	fmt.Println(i)
		// }
		c.JSON(200, list)
	})

	r.GET("/index/listIndexes", func(c *gin.Context) {
		list := app.listIndexes()
		c.JSON(200, list)
	})

	r.POST("/index/add", func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)
		jDat, _ := parseArbJSON(string(data))
		app.addIndex(jDat)
		for k, v := range jDat {
			fmt.Printf("Key: %s Value: %s\n", k, v)
		}
		c.String(200, "Added Index")
	})
}
