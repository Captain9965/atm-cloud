package restapi

import (
	stkapp "cloud/stkApp"
	"github.com/gin-gonic/gin"
)

func RegisterStkApi(r *gin.Engine){
	r.POST("stk_callback", func(c *gin.Context) {
		stkapp.SafStkCallback(c)
	})
}