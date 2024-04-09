package main

import (
	dbapp "cloud/dbApp"
	"cloud/mqttApp"
	restapi "cloud/restApi"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	//Connect to database as the first thing: 
	Db := &dbapp.GormDB{}

	err := Db.Connect()
	if err != nil{
		fmt.Println("Unable to initialize database")
		panic(err)
	}
	r := gin.Default()
	//  Attach middlewares
	r.Use(restapi.DbMiddleware(Db))
	restapi.SetupRouter(r)

	//mqtt go routine
	go mqttApi.RunMqtt()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")

}
