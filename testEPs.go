package main

import (
	"github.com/gin-gonic/gin"
)

func HandleTestRoutes(r *gin.Engine) {
	r.GET("/hc", func(c *gin.Context) {
		c.Status(200)
	})
}
