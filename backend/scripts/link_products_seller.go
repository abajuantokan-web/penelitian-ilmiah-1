//go:build ignore

package main

import (
	"fmt"
	"log"

	"openpeo-backend/models"

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
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	fmt.Println("✅ Connected to db_openpeo")
	fmt.Println("🔄 Fetching the primary SellerProfile...")
	var primarySeller models.SellerProfile
	if err := db.First(&primarySeller).Error; err != nil {
		log.Fatalf("❌ No SellerProfile found in database! Please ensure at least one seller profile exists.")
	}

	fmt.Printf("✅ Found primary seller profile (ID: %d, StoreName: %s)\n", primarySeller.ID, primarySeller.StoreName)
	fmt.Println("🔄 Updating all products to link to this seller profile...")

	// Disable FK checks temporarily for the update if needed, though simple update should work if we just run raw SQL.
	db.Exec("SET FOREIGN_KEY_CHECKS=0;")
	
	// Update all existing products to point to the new seller_profile ID
	result := db.Exec("UPDATE products SET seller_id = ?", primarySeller.ID)
	if result.Error != nil {
		log.Fatalf("❌ Failed to update products: %v\n", result.Error)
	}
	
	db.Exec("SET FOREIGN_KEY_CHECKS=1;")
	fmt.Printf("✅ Successfully updated %d products to seller_id = %d\n", result.RowsAffected, primarySeller.ID)

	fmt.Println("🔄 AutoMigrating Product struct to add FK constraint...")
	if err := db.AutoMigrate(&models.Product{}); err != nil {
		log.Fatalf("❌ Failed to AutoMigrate: %v\n", err)
	}

	fmt.Printf("✅ Successfully updated %d products to seller_id = %d\n", result.RowsAffected, primarySeller.ID)
	fmt.Println("🎉 Database normalization complete.")
}
