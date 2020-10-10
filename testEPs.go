package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
		fmt.Println(jDat["test"])
		c.String(200, "ok")
	})
}
