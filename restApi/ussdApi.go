package restapi

import (
	ussdapp "cloud/ussdApp"
	"github.com/gin-gonic/gin"
)

func RegiserUssdApi(r *gin.Engine){
	r.POST("ussd_callback", func(c *gin.Context) {
		ussdapp.UssdCallback(c)
	})
}