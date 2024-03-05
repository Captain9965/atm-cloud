package stkapp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"time"
	"github.com/google/uuid"
	"cloud/mqttApp"
	// "github.com/gin-gonic/gin"
	"strconv"
)

type AuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

type StkResponseData_t struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

type QueryResponseData_t struct {
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResultCode          string `json:"ResultCode"`
	ResultDesc          string `json:"ResultDesc"`
}

var Token *AuthToken
var amountToSend int

func GetAuthToken() {
	consumerKey := "944ffP8lbnwq9uKoqzgL3SAKeAt3OhJR"
	consumerSecret := "e6yZUQbeXoKzQOud"
	password := consumerKey + ":" + consumerSecret
	b64Password := base64.StdEncoding.EncodeToString([]byte(password))

	url := "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+b64Password)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(body))
	Token = &AuthToken{}
	err = json.Unmarshal(body, Token)
	if err != nil {
		fmt.Println("Error ->", err)
		return
	}
	fmt.Println(Token.AccessToken)
}

func SafStkCallback(c *gin.Context) {
	fmt.Println("Response received..")
}

func SendStkPush(number string, amount int) {
	amountToSend = amount
	GetAuthToken()
	fmt.Println("Sending stk push..")
	// get time:
	loc, err := time.LoadLocation("Africa/Nairobi")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	timeNow := time.Now().In(loc)

	// Format the time manually
	formattedTime := timeNow.Format("20060102150405")
	fmt.Println("Formatted time:", formattedTime)

	//
	mpesa_pass := "bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919"
	short_code := "174379"
	int_short_code, _ := strconv.Atoi(short_code)

	// get password:
	password := short_code + mpesa_pass + formattedTime
	b64Password := base64.StdEncoding.EncodeToString([]byte(password))

	// Check if the phoneNumber starts with '+254'
	if strings.HasPrefix(number, "+254") {
		number = number[1:]
	}

	int_number, err := strconv.Atoi(number)

	if err != nil {
		fmt.Printf("Error in phone number -> %s , %s", err, number)
		return
	}

	url := "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"
	method := "POST"
	client := &http.Client{}
	payload_str := fmt.Sprintf(`{   
		"BusinessShortCode": %d,
		"Password": "%s",
		"Timestamp": "%s",
		"TransactionType": "CustomerPayBillOnline",
		"Amount": %d,
		"PartyA": %s,
		"PartyB": 174379,
		"PhoneNumber": %d,
		"CallBackURL": "https://c0a2-41-90-178-243.ngrok-free.app/stk_callback",
		"AccountReference": "EagoLTD",
		"TransactionDesc": "Payment for beverage"
	 }`, int_short_code, b64Password, formattedTime, amount, number, int_number)

	payload := []byte(payload_str)
	fmt.Println(string(payload))
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err)
		return
	}

	if Token == nil {
		fmt.Println("No token provided")
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return
	}
	Response := &StkResponseData_t{}
	err = json.Unmarshal(body, &Response)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	go checkTransactionInLoop(Response)
}

func checkTransactionInLoop(stkResponse *StkResponseData_t){
	// Wait for 2 seconds..
	time.Sleep(4 * time.Second)
	// Initialize a counter
	count := 0
	// Create a ticker that ticks once every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // Stop the ticker when the loop ends
	// Loop until either count reaches 10 or function returns true
	for range ticker.C {
		count++
		if CheckTransactionStatus(stkResponse) || count >= 20 {
			fmt.Println("Exiting loop.")
			break
		}
	}
}


func CheckTransactionStatus(stkResponse *StkResponseData_t)bool{
	fmt.Println("Checking transaction status...")
	if stkResponse == nil{
		fmt.Println("No new transaction to check..")
		return false
	}
	GetAuthToken()
	// get time:
	loc, err := time.LoadLocation("Africa/Nairobi")

	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	timeNow := time.Now().In(loc)

	// Format the time manually
	formattedTime := timeNow.Format("20060102150405")
	fmt.Println("Formatted time:", formattedTime)

	//
	mpesa_pass := "bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919"
	short_code := "174379"
	int_short_code, _ := strconv.Atoi(short_code)

	// get password:
	password := short_code + mpesa_pass + formattedTime
	b64Password := base64.StdEncoding.EncodeToString([]byte(password))

	url := "https://sandbox.safaricom.co.ke/mpesa/stkpushquery/v1/query"
	method := "POST"
	client := &http.Client{}

	payload_str := fmt.Sprintf(`{
		"BusinessShortCode": %d,
    	"Password": "%s",
    	"Timestamp": "%s",
    	"CheckoutRequestID": "%s"
	}`, int_short_code, b64Password, formattedTime, stkResponse.CheckoutRequestID)

	payload := []byte(payload_str)
	fmt.Println(string(payload))
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err)
		return false
	}

	if Token == nil {
		fmt.Println("No token provided")
		return false
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(string(body))
	Response := &QueryResponseData_t{}
	err = json.Unmarshal(body, &Response)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	switch Response.ResultCode{
		case "0":
			fmt.Println("Transaction was successful..forwading data to mqttApp")
			uid := uuid.New()
			amount := strconv.Itoa(amountToSend)
			payload := fmt.Sprintf(`{"ev":"pay","ms":"45667","a":"%s","u":"%s"}`, amount, uid.String())
			mqttApi.Publish("v/p/30004a-31385105-38353933", payload)
			return true
		case "1032":
			fmt.Println("Request cancelled by the user..")
			return true
		case "2001":
			fmt.Println("The user entered the wrong pin..")
			return true
		default:
			fmt.Println("Let's try this again..")
			return false
	}
}