package restapi

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine){
	RegisterDemoApi(r)
	RegisterC2BApi(r)
	RegisterBasicAuthApi(r)
	RegisterUsersApi(r)
	RegisterOrganizationApi(r)
	RegisterStkApi(r)
	RegiserUssdApi(r)
}