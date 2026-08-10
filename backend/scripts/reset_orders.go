//go:build ignore

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/db_openpeo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v\n", err)
	}

	fmt.Println("✅ Connected to db_openpeo")
	
	// Delete all order items
	db.Exec("DELETE FROM order_items")
	// Delete all orders
	db.Exec("DELETE FROM orders")
	
	fmt.Println("🎉 All legacy orders and order items deleted. The seeder will recreate them on the next backend start.")
}
