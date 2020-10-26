package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
)

func HandleTestRoutes(r *gin.Engine) {
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"yer": "ye",
		})
	})

	r.POST("/test", func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)
		jDat, _ := parseArbJSON(string(data))
		log.Println(jDat["test"])
		c.String(200, "ok")
	})
}
