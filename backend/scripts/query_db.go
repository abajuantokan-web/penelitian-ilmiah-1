//go:build ignore

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID		int32
	Email		string
	Password	string
	Role		string
}

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/db_openpeo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	var user User
	if err := db.Where("email = ?", "admin@openpeo.com").First(&user).Error; err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("User ID: %d, Email: %s, Role: %s\n", user.ID, user.Email, user.Role)
	fmt.Printf("Hash in DB: %s\n", user.Password)
}

