package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func UpgradeToWebsocket(w http.ResponseWriter, r *http.Request, app *appIndexes) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		flatJSON, _ := parseArbJSON(string(msg))
		var body QueryBody
		var output []uint32
		if flatJSON["query"] != nil {
			fmt.Println("Beginning search...")
			json.Unmarshal(msg, &body)
			query := body.Query
			fields := body.Fields
			if fields != nil { // Field(s) specified
				for _, i := range fields {
					fmt.Println(i)
					res := app.searchByField(query, i)
					output = append(output, res...)
				}
			} else {
				res := app.search(query, make([]string, 0))
				output = append(output, res...)
			}
			fmt.Println(output)
			// Convert to array of strings
			out2 := make([]string, len(output))
			for _, item := range output {
				out2 = append(out2, fmt.Sprintf("%v", item))
			}
			out3 := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(output)), ","), "[]")
			conn.WriteMessage(t, []byte(out3))
		} else {
			conn.WriteMessage(t, msg)
		}
	}
}

func HandleWebsocketRoutes(r *gin.Engine, app *appIndexes) {
	r.GET("/ws", func(c *gin.Context) {
		UpgradeToWebsocket(c.Writer, c.Request, app)
	})
}
