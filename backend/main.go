package main

import (
	//"fmt"
	"backend/database"
	//"example.com/project1/logins"
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	r := gin.Default()
	r.GET("/", handlers.LandingPage)
	r.Run(":8088")

}
