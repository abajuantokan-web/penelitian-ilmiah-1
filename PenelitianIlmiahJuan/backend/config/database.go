package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance used throughout the application.
var DB *gorm.DB

// ConnectDatabase initializes the GORM connection to the XAMPP MySQL database.
// It configures connection pooling, logging, and returns any connection errors.
func ConnectDatabase() {
	// DSN for XAMPP MySQL: root user, no password, port 3306, database db_openpeo
	dsn := "root:@tcp(127.0.0.1:3306)/db_openpeo?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Configure connection pool for optimal performance
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("✅ Database connection established successfully")
}
