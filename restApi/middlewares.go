package restapi

import(
	"github.com/gin-gonic/gin"
	dbapp "cloud/dbApp"
)

// middleware function to set the db connection object within the context for all routes
func DbMiddleware(db dbapp.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
	  c.Set("db", db)
	  c.Next()
	}
}
