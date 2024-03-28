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


func setupRouter(r *gin.Engine){
	// instance with logger and recovery middleware

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

	// create a user
	r.POST("/createuser", func(c *gin.Context) {
		Db := c.MustGet("db").(dbapp.Database)
	
		// Bind JSON from request body
		var newUser struct {
		  Username string  `json:"username" binding:"required"`
		  Password string  `json:"password" binding:"required"`
		  Role int		   `json:"role" binding:"required"`
		  Org string	   `json:"org" binding:"required"`
		  PhoneNumber string `json:"phonenumber" binding:"required"`
		}

		if err := c.Bind(&newUser); err != nil {
			fmt.Println(newUser)
		  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		  return
		}
	
		// Create the superuser in the database
		err := Db.CreateUser(map[string]interface{}{
		  "username": newUser.Username,
		  "password": newUser.Password,
		  "role"	: newUser.Role,
		  "org"		: newUser.Org,
		  "phonenumber" : newUser.PhoneNumber,
		})

		if err != nil {
		  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "Description":err.Error()})
		}
	
		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	  })

	  // get all users
	  r.GET("/users", func(c *gin.Context) {
		Db := c.MustGet("db").(dbapp.Database)
	
		// Get all users
		users, err := Db.GetAllUsers()

		if err != nil {
		  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users", "Description":err.Error()})
		
		}
			c.JSON(http.StatusOK, users)
		})

		r.DELETE("/deleteuser", func(c *gin.Context) {
			Db := c.MustGet("db").(dbapp.Database)

			var user struct{
				Username string `json:"username" binding:"required"`
			}
			if err := c.Bind(&user); err != nil{
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			}
			err := Db.DeleteUserByName(user.Username)
			if err != nil{
				c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete user", "Description":err.Error()})
		  return
			}
			c.JSON(http.StatusOK, gin.H{"User deleted": user.Username})
		})
	
		


	// create a user
	r.POST("/createorg", func(c *gin.Context) {
		Db := c.MustGet("db").(dbapp.Database)
	
		// Bind JSON from request body
		var newOrg struct {
		  Orgname string  `json:"orgname" binding:"required"`
		}

		if err := c.Bind(&newOrg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		err := Db.CreateOrganization(map[string]interface{}{
		  "orgname": newOrg.Orgname,
		})

		if err != nil {
		  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization", "Description":err.Error()})
		  return
		}
	
		c.JSON(http.StatusCreated, gin.H{"message": "Organization created successfully"})
	  })

	  // get all organizations
	  r.GET("/orgs", func(c *gin.Context) {
		Db := c.MustGet("db").(dbapp.Database)
	
		// Get all orgs
		orgs, err := Db.GetAllOrganizations()

		if err != nil {
		  c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organiizations", "Description":err.Error()})
		  
		}
	
		c.JSON(http.StatusOK, orgs)
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

// middleware function to set the db connection object within the context for all routes
func dbMiddleware(db dbapp.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
	  c.Set("db", db)
	  c.Next()
	}
}

func main() {
	//Connect to database as the first thing: 
	Db := &dbapp.GormDB{}

	err := Db.Connect()
	if err != nil{
		fmt.Println("Unable to initialize database")
		panic(err)
	}
	
	r := gin.Default()
	r.Use(dbMiddleware(Db))
	setupRouter(r)

	//mqtt go routine
	go mqttApi.RunMqtt()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")

}
