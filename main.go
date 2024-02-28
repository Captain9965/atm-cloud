package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	// "time"
	"cloud/mqttApp"
	"strconv"
)

var db = make(map[string]string)

type AuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

var authToken *AuthToken

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

	return r
}

//	func publish() {
//		client := mqttApi.Connect("lenny", "broker.hivemq.com:1883", "lenny", "Lenny123")
//		timer := time.NewTicker(1 * time.Second)
//		for t := range timer.C {
//			client.Publish("me", 0, false, t.String())
//		}
//	}
func getToken() {

	consumerKey := "ITHYll8TtraEAEgQzmh9xlHPSBB3fjIsGMEiAreNBAzvnnkg"
	consumerSecret := "Pc7nWwXoCO8qBzQAjhBGPzpm0vGW8FXO0u0GusIArTiCv5QhrfLbC4rJeTxb3IhU"
	password := consumerKey + ":" + consumerSecret
	b64Password := base64.StdEncoding.EncodeToString([]byte(password))

	url := "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+b64Password)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(body))
	authToken = &AuthToken{}
	err = json.Unmarshal(body, authToken)
	if err != nil {
		fmt.Println("Error ->", err)
		return
	}
	fmt.Println(authToken.AccessToken)

}

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
	req.Header.Add("Authorization", "Bearer "+authToken.AccessToken)

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
	//loadtoken
	getToken()
	//mqtt go routine
	go mqttApi.RunMqtt()
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
