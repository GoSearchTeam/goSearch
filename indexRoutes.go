package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
)

type QueryBody struct {
	Query      string   `json:"query"`
	Fields     []string `json:"fields"`
	BeginsWith bool     `json:"beginsWith"`
}

type SearchResponse struct {
	DocIDs    []uint64      `json:"docIDs"`
	Documents []interface{} `json:"documents"`
}

type AddMultipleIndex struct {
	Items []map[string]interface{} `json:"items"`
}

type DeleteDocumentRequestBody struct {
	DocID uint64 `json:"docID"`
}

type UpdateDocumentRequestBody struct {
	DocID uint64 `json:"docID"`
}

func HandleIndexRoutes(r *gin.Engine, app *appIndexes) {
	r.GET("/index/listItems", func(c *gin.Context) {
		list := app.listIndexItems()
		// stringed := make([]string, 0)
		// for _, i := range list {
		// 	log.Println(i)
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
		body := QueryBody{BeginsWith: false} // Default false if not included
		c.BindJSON(&body)
		var output []uint64 // temporary, will turn into documents later
		documents := make([]string, 0)
		// query := jDat["query"].(string)
		query := body.Query
		fields := body.Fields
		bw := body.BeginsWith
		if fields != nil { // Field(s) specified
			res, docs := app.search(query, fields, bw)
			output = append(output, res...)
			documents = append(documents, docs...)
		} else {
			res, docs := app.search(query, make([]string, 0), bw)
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
		docID := app.addIndex(jDat)
		c.JSON(200, gin.H{
			"msg":   "Added Index",
			"docID": docID,
		})
	})

	r.POST("/index/addMultiple", func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)
		jDat := AddMultipleIndex{}
		json.Unmarshal([]byte(data), &jDat)
		docIDs := make([]uint64, 0)
		for _, item := range jDat.Items {
			docIDs = append(docIDs, app.addIndex(item))
		}
		c.JSON(200, gin.H{
			"msg":    "Added Indexes",
			"docIDs": docIDs,
		})
	})

	r.POST("/index/delete", func(c *gin.Context) {
		body := DeleteDocumentRequestBody{}
		c.BindJSON(&body)
		err := app.deleteIndex(body.DocID)
		if err != nil {
			log.Printf("%v\n", err)
			c.String(500, fmt.Sprintf("Internal Error: %v", err))
		} else {
			c.String(200, fmt.Sprintf("Deleted document %v", body.DocID))
		}
	})

	r.POST("/index/update", func(c *gin.Context) {
		d := json.NewDecoder(c.Request.Body)
		d.UseNumber()
		var jDat map[string]interface{}
		if err := d.Decode(&jDat); err != nil {
			log.Fatal(err)
		}
		err := app.updateIndex(jDat)
		if err != nil {
			log.Printf("%v\n", err)
			c.String(500, fmt.Sprintf("Internal Error: %v", err))
		} else {
			c.String(200, "Updated Index")
		}
	})
}
