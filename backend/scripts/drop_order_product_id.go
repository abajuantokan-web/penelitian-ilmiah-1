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
	fmt.Println("🔄 Dropping product_id column from orders table...")

	err = db.Exec("ALTER TABLE orders DROP FOREIGN KEY fk_orders_product;").Error
	if err != nil {
		fmt.Printf("⚠️ Warning dropping foreign key (may not exist): %v\n", err)
	}

	err = db.Exec("ALTER TABLE orders DROP COLUMN product_id;").Error
	if err != nil {
		log.Fatalf("❌ Failed to drop column product_id: %v\n", err)
	}

	fmt.Println("🎉 Migration complete: product_id removed from orders table.")
}

