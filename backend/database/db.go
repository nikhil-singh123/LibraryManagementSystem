package database

import (
	"fmt"

	"backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := "Nikhilsingh:password@tcp(127.0.0.1:3306)/NIKHIL?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect with Database")
	}

	fmt.Println("Database Connected Sucessfully")
	DB.AutoMigrate(&models.Library{})
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.BookInventory{})
	DB.AutoMigrate(&models.RequestEvent{})
	DB.AutoMigrate(&models.IssueRegistry{})

}
