//go:build ignore

package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash := "$2a$10$wK2Vv.6vM6o8O/FkE3tPoe2i.ZfO6jH8VzW5C/cO1k9G6nU.M0FvW"
	password := "password123"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Hash mismatch!", err)
	} else {
		fmt.Println("Hash matches!")
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	fmt.Println("New Hash:", string(newHash))
}

