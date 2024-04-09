package restapi

import(
	"github.com/gin-gonic/gin"
	"net/http"
)

// ping pong test..
func RegisterDemoApi(r *gin.Engine){
	r.GET("/ping", func(c *gin.Context) {
		PingPong(c)
	})

}

func PingPong(c *gin.Context){
	c.String(http.StatusOK, "pong")
	return
}