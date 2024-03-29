package main

import (
	"backend/database"
	"backend/handlers"
	//"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB() //database intialisation

	r := gin.Default() //router created

	r.GET("/", handlers.LandingPage) //landing page

	//admin roles

	r.POST("/libraries", handlers.CreateLibrary)                                //library creation
	r.POST("/add-book", handlers.AddBook)                                       //addition of books
	r.PATCH("/update-book", handlers.UpdateBook)                                //updation of books
	r.DELETE("/remove-book", handlers.RemoveBook)                               //deletion of books
	r.GET("/list-issue-request", handlers.ListIssueRequests)                    //list of issues request
	r.POST("/approve-issue-request/:request_id", handlers.ApprovedIssueRequest) //approve request
	r.DELETE("/reject-issue-request/:request_id", handlers.RejectIssueRequest)  //reject request

	//readers role
	r.POST("/search/book", handlers.SearchBook)          //search book
	r.POST("/issue/request", handlers.RaiseIssueRequest) //raise request

	r.Run(":8088") //running on port localhost:8088

}
