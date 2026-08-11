//go:build ignore

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID		int32	`gorm:"primaryKey"`
	Role		string	`gorm:"size:20"`
	StoreName	string	`gorm:"size:100"`
}

func (User) TableName() string	{ return "users" }

type SellerProfile struct {
	ID		int32	`gorm:"primaryKey;autoIncrement"`
	UserID		int32	`gorm:"uniqueIndex;not null"`
	StoreName	string	`gorm:"size:255"`
}

func (SellerProfile) TableName() string	{ return "seller_profiles" }

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
	fmt.Println("🔄 Migrating store_name from users → seller_profiles...")

	var sellers []User
	if err := db.Where("role = ?", "seller").Find(&sellers).Error; err != nil {
		log.Fatalf("❌ Failed to query sellers: %v\n", err)
	}
	fmt.Printf("   Found %d seller(s)\n\n", len(sellers))

	created, updated, skipped := 0, 0, 0

	for _, seller := range sellers {
		if seller.StoreName == "" {
			fmt.Printf("   [SKIP   ] #%-3d ← users.store_name is empty, nothing to migrate\n", seller.ID)
			skipped++
			continue
		}

		var profile SellerProfile
		err := db.Where("user_id = ?", seller.ID).First(&profile).Error

		if err != nil {

			newProfile := SellerProfile{UserID: seller.ID, StoreName: seller.StoreName}
			if err := db.Create(&newProfile).Error; err != nil {
				log.Printf("   [ERROR  ] #%d create profile → %v\n", seller.ID, err)
				continue
			}
			fmt.Printf("   [CREATED] #%-3d → seller_profiles.store_name = %q\n", seller.ID, seller.StoreName)
			created++
		} else if profile.StoreName == "" {

			if err := db.Model(&profile).UpdateColumn("store_name", seller.StoreName).Error; err != nil {
				log.Printf("   [ERROR  ] #%d update profile → %v\n", seller.ID, err)
				continue
			}
			fmt.Printf("   [UPDATED] #%-3d → seller_profiles.store_name = %q\n", seller.ID, seller.StoreName)
			updated++
		} else {
			fmt.Printf("   [SKIP   ] #%-3d ← seller_profiles.store_name already set: %q\n", seller.ID, profile.StoreName)
			skipped++
		}
	}

	fmt.Println("\n────────────────────────────────────────────────────────────")
	fmt.Printf("✅ Done.  Created: %d | Updated: %d | Skipped: %d | Total: %d\n",
		created, updated, skipped, len(sellers))
}

