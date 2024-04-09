package restapi

import (
	mqttApi "cloud/mqttApp"
	safc2bapp "cloud/safc2bApp"
	"fmt"
	"log"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterC2BApi(r *gin.Engine){
	r.GET("/register_urls", func(c *gin.Context) {
		safc2bapp.RegisterURL()
	})
	r.POST("confirmation", func(c *gin.Context) {
		confirmPaymentCallback(c)
	})
}

func confirmPaymentCallback(c *gin.Context){
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
}