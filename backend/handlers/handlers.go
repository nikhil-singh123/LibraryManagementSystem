package handlers

import (
	//"backend/database"
	//"backend/models"
	"backend/database"
	"backend/models"
	"net/http"

	//"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func LandingPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello the website is running."})
}

//Creating the Library

var num uint

func CreateLibrary(c *gin.Context) {
	var request struct {
		LibraryName string `json:"library_name"`
		OwnerName   string `json:"owner_name"`
		OwnerEmail  string `json:"owner_email"`
		OwnerPhone  string `json:"owner_phone"`
		OwnerRole   string `json:"owne_role"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Condition for check if library already exists

	var existingLibrary models.Library

	result := database.DB.Where("name=?", request.LibraryName).First(&existingLibrary)

	if request.OwnerRole == "Admin" {
		if result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Library is already existed with this User Name"})
			return
		}
	}

	//If library doesn't exist

	newLibrary :=models.Library{
		Name: request.LibraryName,
	}

	if request.OwnerRole!="Reader"{
		if err:=database.DB.Create(&newLibrary).Error;err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create library"})
			return
		}
	}

	if request.OwnerRole=="Reader"{
		if result.RowsAffected>0{
			num=existingLibrary.ID
		}
	}else{
		num=newLibrary.ID
	}


	//Create Owner User
	var mod models.User
	newOwner :=models.User{
		Name:	request.OwnerName,
		Email:  request.OwnerEmail,
		ContactNumber: request.OwnerPhone,
		Role:    request.OwnerRole,
		LibID:   num,
	}

	if err := database.DB.Create(&newOwner).Error; err!=nil{
		mod.ID--
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Creation of User Failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message":"Library Created Succesfully", "library_id":num})

}
