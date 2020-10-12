package main

import (
	"github.com/gin-gonic/gin"
	// "net/http"
)

func StartWebserver(app *appIndexes) {
	r := gin.Default()
	HandleTestRoutes(r)
	HandleIndexRoutes(r, app)
	HandleWebsocketRoutes(r, app)
	r.Run("localhost:8080")
}
