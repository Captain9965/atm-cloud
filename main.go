package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"time"

	"cloud/mqttApp"

)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// instance with logger and recovery middleware
	r := gin.Default()

	// Ping pong
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		//user password pairs
		"lenny":  "nullpass",
		"Robert": "123",
	}))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok", "user": user})
		}
	})

	return r
}

func publish() {
	client := mqttApi.Connect("lenny", "broker.hivemq.com:1883", "lenny", "Lenny123")
	timer := time.NewTicker(1 * time.Second)
	for t := range timer.C {
		client.Publish("me", 0, false, t.String())
	}
}

func runMqtt(){
	mqttApi.Listen("broker.hivemq.com:1883", "me", "client", "lenny", "Lenny123")
	publish()
}

func main() {
	//mqtt go routine
	go runMqtt()
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
