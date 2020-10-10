package main

import (
	"github.com/gin-gonic/gin"
	// "net/http"
)

func StartWebserver() {
	r := gin.Default()
	HandleTestRoutes(r)
	r.Run("localhost:8080")
}
