package main

import (
	"github.com/gin-gonic/gin"
	// "net/http"
)

func StartWebserver() {
	r := gin.Default()
	app := initApp("Example App")
	HandleTestRoutes(r)
	HandleIndexRoutes(r, app)
	r.Run("localhost:8080")
}
