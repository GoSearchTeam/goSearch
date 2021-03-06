package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type QueryBody struct {
	Query      string   `json:"query"`
	Fields     []string `json:"fields"`
	BeginsWith bool     `json:"beginsWith"`
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

	r.GET("/index/listIndexes", func(c *gin.Context) {
		list := app.listIndexes()
		for listItem, _ := range list {
			list[listItem].Name = strings.ReplaceAll(list[listItem].Name, "\\.", ".")
		}
		c.JSON(200, list)
	})

	r.POST("/index/search", func(c *gin.Context) {
		// data, _ := ioutil.ReadAll(c.Request.Body)
		// jDat, _ := parseArbJSON(string(data))
		body := QueryBody{BeginsWith: false} // Default false if not included
		c.BindJSON(&body)
		var output []uint64 // temporary, will turn into documents later
		// query := jDat["query"].(string)
		query := body.Query
		fields := body.Fields
		bw := body.BeginsWith
		var responseJSON SearchResponse
		var res []uint64
		if fields != nil { // Field(s) specified
			res, responseJSON = app.search(query, fields, bw)
			output = append(output, res...)
		} else {
			res, responseJSON = app.search(query, make([]string, 0), bw)
			output = append(output, res...)
		}
		for responseItem, _ := range responseJSON.Items {
			responseJSON.Items[responseItem].Data, _ = nestJSON(responseJSON.Items[responseItem].Data)
		}
		c.JSON(200, responseJSON)
	})

	r.POST("/index/add", func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)
		jDat, _ := parseArbJSON(string(data))
		flatDat, _ := flattenJSON(jDat)
		docID := app.addIndex(flatDat, false)
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
			docIDs = append(docIDs, app.addIndex(item, false))
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
		err, created := app.updateIndex(jDat)
		if err != nil {
			log.Println("Error:", err)
			c.String(500, "Internal Error, check logs")
		} else if created {
			c.String(200, "Created new Index")
		} else {
			c.String(200, "Updated Index")
		}
	})
}
