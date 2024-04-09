package restapi

import(
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterBasicAuthApi(r * gin.Engine) *gin.RouterGroup{
	
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		//demo user password pairs
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
			c.JSON(http.StatusOK, gin.H{"status": "ok", "user": user})
			return
		}
	})

	return authorized
}