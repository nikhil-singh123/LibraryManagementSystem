package handlers

import (
	"backend/database"
	"backend/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	newLibrary := models.Library{
		Name: request.LibraryName,
	}

	if request.OwnerRole != "Reader" {
		if err := database.DB.Create(&newLibrary).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create library"})
			return
		}
	}

	if request.OwnerRole == "Reader" {
		if result.RowsAffected > 0 {
			num = existingLibrary.ID
		}
	} else {
		num = newLibrary.ID
	}

	//Create Owner User
	var mod models.User
	newOwner := models.User{
		Name:          request.OwnerName,
		Email:         request.OwnerEmail,
		ContactNumber: request.OwnerPhone,
		Role:          request.OwnerRole,
		LibID:         num,
	}

	if err := database.DB.Create(&newOwner).Error; err != nil {
		mod.ID--
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Creation of User Failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Library Created Succesfully", "library_id": num})

}

//Add Book feature

func AddBook(c *gin.Context) {
	var request struct {
		Book       models.BookInventory `json:"book"`
		AdminEmail string               `json:"admin_email"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Validating email of admin
	var existingEmail models.User
	result1 := database.DB.Where("email=?", request.AdminEmail).First(&existingEmail)
	if result1.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, "Admin email is incorrect")
		return
	}

	//check for availabiltiy of book
	var exisitingBook models.BookInventory
	result := database.DB.Where("isbn=?", request.Book.ISBN).First(&exisitingBook)
	if result.RowsAffected > 0 {
		exisitingBook.TotalCopies++
		if err := database.DB.Save(&exisitingBook).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to updating book"})
			return
		}

		c.JSON(http.StatusOK, exisitingBook)
		return
	}

	//Book is not present
	request.Book.TotalCopies = 1
	if err := database.DB.Create(&request.Book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the  book"})
		return
	}

}

// Removing book
func RemoveBook(c *gin.Context) {
	var request struct {
		ISBN string `json:"isbn"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	//Validating if thee book is present or not

	var book models.BookInventory
	result := database.DB.Where("isbn=?", request.ISBN).First(&book)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Validating if book is present but doesn't have copies left in library
	if book.TotalCopies == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No copies available in Library"})
		return
	}

	// Validating if copy of book is issued

	if book.AvailableCopies < book.TotalCopies {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Copy of book can't be removed"})
		return
	}

	//Function for decreasing the copies
	book.TotalCopies--

	if err := database.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update the book"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book removed!"})
}

// Function for updation of Books
func UpdateBook(c *gin.Context) {
	var request struct {
		ISBN           string               `json:"isbn"`
		UpdatedDetails models.BookInventory `json:"updated-details"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//condition to check if the book exists
	var book models.BookInventory
	result := database.DB.Where("isbn=?", request.ISBN).First(&book)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No books found by the provided ISBN"})
		return
	}

	//updation of details
	if err := database.DB.Model(&book).Updates(request.UpdatedDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

//Listing issue request

func ListIssueRequests(c *gin.Context) {
	var issueRequests []models.RequestEvent

	//fetching request from database

	if err := database.DB.Find(&issueRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch issue requests"})
		return
	}
	c.JSON(http.StatusOK, issueRequests)
}

//request approval

func ApprovedIssueRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}
	var request models.RequestEvent
	if err := database.DB.First(&request, requestID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	//update request details

	request.ApprovalDate = time.Now()
	request.ApproverID = 1
	if err := database.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve requests "})
		return
	}

	var issueRegistry models.IssueRegistry
	//putting value in issue registry
	result := database.DB.Where("reader_id = ?", requestID).First(&issueRegistry)
	fmt.Println(requestID)
	if result.RowsAffected == 0 {
		issueRegistry := models.IssueRegistry{
			ISBN:               request.BookID,
			ReaderID:           request.ReaderID,
			IssueApproverID:    request.ApproverID,
			IssueStatus:        "APPROVED",
			IssueDate:          time.Now(),
			ExpectedReturnDate: time.Now(),
			ReturnDate:         time.Now(),
			ReturnApproverID:   request.ApproverID,
		}
		if err := database.DB.Create(&issueRegistry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create issue request"})
			return
		}

	} else {
		if err := database.DB.Where("reader_id = ?", request.ReaderID).First(&issueRegistry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find issue registory"})
			return
		}

		issueRegistry.IssueStatus = "Approved"
		issueRegistry.IssueDate = time.Now()

		if err := database.DB.Save(&issueRegistry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update issue registry"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Issue request approved successfully"})

}

func RejectIssueRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var request models.RequestEvent
	if err := database.DB.First(&request, requestID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	//Deleting request

	if err := database.DB.First(&request, requestID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "issue request rejected successfully"})
}

//Search for Books matching

func SearchBook(c *gin.Context) {
	var request struct {
		Title     string `json:"title"`
		Author    string `json:"author"`
		Publisher string `json:"publisher"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var books models.BookInventory
	result := database.DB.Where("title LIKE ?", "%"+request.Title+"%").Where("authors LIKE ?", "%"+request.Author+"%").Where("publisher LIKE ?", "%"+request.Publisher+"%").Find(&books)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

//Raise an issue request

func RaiseIssueRequest(c *gin.Context) {
	var request struct {
		BookID string `json:"book_id"`
		Email  string `json:"email"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check if book exists

	var book models.BookInventory
	result := database.DB.Where("isbn=?", request.BookID).First((&book))
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book is not found"})
		return
	}

	//Providing approver ID
	var col models.User
	result1 := database.DB.Where("email = ?", request.Email).First(&col)
	if result1.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	//Check if book is available
	if book.AvailableCopies == 0 {
		//Create a new issue request
		issueRequest := models.RequestEvent{
			BookID:      request.BookID,
			ReaderID:    col.ID,
			RequestDate: time.Now(),
			RequestType: "Issue",
		}

		if err := database.DB.Create(&issueRequest).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create issue request"})
			return
		}
		//printing issue request in json format
		c.JSON(http.StatusCreated, issueRequest)
		return
	}
	book.AvailableCopies--
	if err := database.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book availability"})
		return
	}

}
