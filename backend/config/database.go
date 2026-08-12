package config

import (
	"fmt"
	"log"
	"os" // 1. Tambahkan package os di sini
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Mengambil variabel dari Environment Vercel/Railway, atau fallback ke lokal (XAMPP)
	dbHost := os.Getenv("MYSQLHOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}

	dbUser := os.Getenv("MYSQLUSER")
	if dbUser == "" {
		dbUser = "root"
	}

	dbPassword := os.Getenv("MYSQLPASSWORD")
	if dbPassword == "" {
		dbPassword = os.Getenv("MYSQL_ROOT_PASSWORD") // Cadangan untuk lokal
	}

	dbName := os.Getenv("MYSQL_DATABASE")
	if dbName == "" {
		dbName = "db_openpeo"
	}

	dbPort := os.Getenv("MYSQLPORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	// Menyusun DSN secara dinamis
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("✅ Database connection established successfully")
}
