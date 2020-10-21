package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type QueryBody struct {
	Query  string   `json:query`
	Fields []string `json:fields`
}

type SearchResponse struct {
	DocIDs    []uint64      `json:docIDs`
	Documents []interface{} `json:documents`
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
		var body QueryBody
		c.BindJSON(&body)
		var output []uint64 // temporary, will turn into documents later
		documents := make([]string, 0)
		// query := jDat["query"].(string)
		query := body.Query
		fields := body.Fields
		if fields != nil { // Field(s) specified
			res, docs := app.search(query, fields)
			output = append(output, res...)
			documents = append(documents, docs...)
		} else {
			res, docs := app.search(query, make([]string, 0))
			output = append(output, res...)
			documents = append(documents, docs...)
		}
		jsonDocs := make([]interface{}, 0)
		for _, doc := range documents {
			parsed, _ := parseArbJSON(doc)
			jsonDocs = append(jsonDocs, parsed)
		}
		c.JSON(200, SearchResponse{output, jsonDocs})
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
