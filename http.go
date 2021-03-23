package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "net/http"
)

func StartWebserver(app *appIndexes) {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))
	HandleTestRoutes(r)
	HandleIndexRoutes(r, app)
	HandleWebsocketRoutes(r, app)
	HandleClusterRoutes(r)
	r.Run(fmt.Sprintf("0.0.0.0:%d", *APIPort))
}
