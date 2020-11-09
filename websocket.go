package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var jwts []jwt.Token

var jwtSecret string = "thisisanexamplesecret"

type WebsocketResponse struct {
	TimeNS    int64            `json:"timeNS"`
	Documents []DocumentObject `json:"documents"`
}

type WebsocketAuthValidResponse struct {
	JWT string `json:"token"`
}

func UpgradeToWebsocket(w http.ResponseWriter, r *http.Request, c *gin.Context, app *appIndexes) {
	// Validate ws connection
	jwtQuery := c.Request.URL.Query()["token"][0]

	token, err := jwt.Parse(jwtQuery, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check exists

	} else if int64(claims["exp"].(float32)) <= time.Now().Unix() { // Check expired
		c.JSON(401, gin.H{
			"msg": "Token Expired",
		})
		return
	} else {
		c.JSON(401, gin.H{
			"msg": "Token Invalid",
		})
		return
	}

	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade: %v", err)
		return
	}
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		start := time.Now()
		flatJSON, _ := parseArbJSON(string(msg))
		body := QueryBody{BeginsWith: false}
		var output []uint64
		documents := make([]DocumentObject, 0)
		var responseJSON SearchResponse
		fmt.Println(responseJSON)
		var res []uint64
		if flatJSON["query"] != nil {
			json.Unmarshal(msg, &body)
			query := body.Query
			fields := body.Fields
			bw := body.BeginsWith
			if fields != nil { // Field(s) specified
				res, responseJSON = app.search(query, fields, bw)
				output = append(output, res...)
			} else {
				res, responseJSON = app.search(query, make([]string, 0), bw)
				output = append(output, res...) // TODO: Send documents as well
			}
			// Convert to array of strings
			out2 := make([]string, len(output))
			for _, item := range output {
				//fmt.Sprintf("%d", item);
				out2 = append(out2, fmt.Sprintf("%v", item))
			}
			// out3 := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(output)), ","), "[]")
			end := time.Now()
			resStruct := &WebsocketResponse{
				TimeNS:    end.Sub(start).Nanoseconds(),
				Documents: documents,
			}
			resBody, _ := json.Marshal(resStruct)
			conn.WriteMessage(t, []byte(resBody))
		} else {
			conn.WriteMessage(t, msg)
		}
	}
}

func HandleWebsocketRoutes(r *gin.Engine, app *appIndexes) {
	r.GET("/ws", func(c *gin.Context) {
		UpgradeToWebsocket(c.Writer, c.Request, c, app)
	})

	r.GET("/ws/auth", func(c *gin.Context) {
		newUUID, _ := uuid.NewRandom()
		userID := fmt.Sprintf("user#%v", newUUID)
		jwtClaims := jwt.MapClaims{}
		var expireMin time.Duration
		expireMin, err := time.ParseDuration(os.Getenv("JWT_EXP"))
		fmt.Println(os.Getenv("JWT_EXP"))
		if err != nil {
			expireMin = 20 * time.Minute
		}
		jwtClaims["exp"] = time.Now().Add(expireMin).Unix() // 20 minutes default
		jwtClaims["iat"] = time.Now().Unix()
		jwtClaims["userID"] = userID
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
		token, _ := accessToken.SignedString([]byte(jwtSecret))
		c.JSON(200, WebsocketAuthValidResponse{token})
	})
}
