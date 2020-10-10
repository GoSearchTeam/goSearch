package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type queryBody struct {
	Query  string   `json:query`
	Fields []string `json:fields`
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

	r.POST("/index/search", func(c *gin.Context) {
		// data, _ := ioutil.ReadAll(c.Request.Body)
		// jDat, _ := parseArbJSON(string(data))
		var body queryBody
		c.BindJSON(&body)
		fmt.Println(body)
		var output []uint32 // temporary, will turn into documents later
		// query := jDat["query"].(string)
		query := body.Query
		fields := body.Fields
		if fields != nil { // Field(s) specified
			for _, i := range fields {
				fmt.Println(i)
				res := app.searchByField(query, i)
				output = append(output, res...)
			}
		} else {
			res := app.search(query, make([]string, 0))
			output = append(output, res...)
		}
		c.JSON(200, output)
	})

	r.POST("/index/add", func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)
		jDat, _ := parseArbJSON(string(data))
		app.addIndex(jDat)
		// for k, v := range jDat {
		// 	fmt.Printf("Key: %s Value: %s\n", k, v)
		// }
		c.String(200, "Added Index")
	})
}
