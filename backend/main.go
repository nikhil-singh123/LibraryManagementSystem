package main

import (
	//"fmt"
	"backend/database"
	//"backend/handlers"
	//"example.com/project1/logins"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()
	r.Run(":8088")

}
