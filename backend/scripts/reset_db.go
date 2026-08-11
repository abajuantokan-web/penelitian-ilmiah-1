//go:build ignore

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/db_openpeo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.Exec("DROP TABLE IF EXISTS order_items, orders, cart_items, products, users").Error
	if err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}

	fmt.Println("Tables dropped successfully.")
}

