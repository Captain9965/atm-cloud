package main

import (
	"bytes"
	dbapp "cloud/dbApp"
	"cloud/mqttApp"
	stkapp "cloud/stkApp"
	"cloud/ussdApp"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	r.GET("/register_urls", func(c *gin.Context) {
		registerURL()
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

	r.POST("confirmation", func(c *gin.Context) {
		// Parse JSON
		var confirmation struct {
			Amount int `json:"TransAmount"`
		}

		if c.Bind(&confirmation) == nil {
			log.Printf("Received kshs %d", confirmation.Amount)
			uid := uuid.New()
			amount := strconv.Itoa(confirmation.Amount)
			payload := fmt.Sprintf(`{"ev":"pay","ms":"45667","a":"%s","u":"%s"}`, amount, uid.String())
			mqttApi.Publish("v/p/30004a-31385105-38353933", payload)

			c.JSON(http.StatusOK, gin.H{"ResultCode": 0, "ResultDesc": "Success"})
		} else {
			log.Printf("error occured")
			c.JSON(http.StatusInternalServerError, gin.H{"ResultCode": 0, "ResultDesc": "Fail"})
		}
	})

	r.POST("ussd_callback", func(c *gin.Context) {
		ussdApp.UssdCallback(c)
	})
	r.POST("stk_callback", func(c *gin.Context) {
		stkapp.SafStkCallback(c)
	})

	return r
}

//	func publish() {
//		client := mqttApi.Connect("lenny", "broker.hivemq.com:1883", "lenny", "Lenny123")
//		timer := time.NewTicker(1 * time.Second)
//		for t := range timer.C {
//			client.Publish("me", 0, false, t.String())
//		}
//	}

func registerURL() {
	url := "https://sandbox.safaricom.co.ke/mpesa/c2b/v1/registerurl"
	method := "POST"
	client := &http.Client{}
	payload := []byte(`{   
		"ShortCode": "600999",
		"ResponseType":"Completed",
		"ConfirmationURL":"https://google.com/confirmation",
		"ValidationURL":"https://google.com/validation"
	 }`)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+stkapp.Token.AccessToken)

	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Error ->", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error ->", err)
		return
	}
	fmt.Println(string(body))
}

func main() {
	//First connect to database:
	dbapp.ConnectDatabase()
	//mqtt go routine
	go mqttApi.RunMqtt()
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
