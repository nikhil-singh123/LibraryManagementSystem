package handlers

import (
	"backend/database"
	"backend/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func LandingPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello the website is running."})
}

func AddBook(c *gin.Context){
	var request struct{
		Book		models.BookInventory	`json:"book"`
		AdminEmail	string 					`json:"email`		
	}

	
}