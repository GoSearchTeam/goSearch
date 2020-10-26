package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	// "net/http"
)

func StartWebserver(app *appIndexes) {
	logFile, _ := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))
	HandleTestRoutes(r)
	HandleIndexRoutes(r, app)
	HandleWebsocketRoutes(r, app)
	r.Run("localhost:8080")
}
