package main

import (
	
	"backend/database"
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	r := gin.Default()
	r.GET("/", handlers.LandingPage)
	r.POST("/libraries",handlers.CreateLibrary)
	r.POST("/add-book", handlers.AddBook)
	r.Run(":8088")

}
