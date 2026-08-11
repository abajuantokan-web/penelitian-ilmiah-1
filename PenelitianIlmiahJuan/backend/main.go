package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"openpeo-backend/config"
	"openpeo-backend/handlers"
	"openpeo-backend/models"
)

func main() {
	
	
	
	config.ConnectDatabase()

	
	
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.Message{},
	)
	if err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}
	fmt.Println("✅ Database schema synchronized")

	
	seedDemoData()

	
	
	
	r := gin.Default()

	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	
	
	
	api := r.Group("/api")
	{
		
		api.POST("/login", handlers.Login)
		api.POST("/register", handlers.RegisterUser)
		api.GET("/user/:id", handlers.GetCurrentUser)

		
		api.POST("/products", handlers.CreateProduct)
		api.GET("/products", handlers.GetProducts)
		api.PUT("/products/:id", handlers.UpdateProduct)
		api.DELETE("/products/:id", handlers.DeleteProduct)

		
		api.POST("/orders", handlers.CreateOrder)
		api.GET("/orders", handlers.GetOrders)
		api.DELETE("/orders/:id", handlers.DeleteOrder)

		
		api.GET("/messages", handlers.GetChatHistory)
		api.GET("/chat/contacts", handlers.GetChatContacts)

		
		admin := api.Group("/admin")
		{
			admin.GET("/sales-data", handlers.GetSalesData)
		}
	}

	
	
	
	r.GET("/ws/chat", handlers.HandleWebSocket)

	
	
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "OpenPeo Backend Engine",
			"version": "1.0.0",
		})
	})

	
	
	
	fmt.Println("🚀 OpenPeo Backend Engine starting on :8080")
	fmt.Println("📡 REST API:    http://localhost:8080/api")
	fmt.Println("🔌 WebSocket:   ws://localhost:8080/ws/chat")
	fmt.Println("❤️  Health:      http://localhost:8080/health")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}



func seedDemoData() {
	var userCount int64
	config.DB.Model(&models.User{}).Count(&userCount)

	
	if userCount == 0 {
		fmt.Println("🌱 Seeding demo users...")
		users := []models.User{
			{
				Name:     "Admin Flores",
				Username: "admin_flores",
				Email:    "admin@openpeo.id",
				Password: "password123",
				Role:     "admin",
				Phone:    "081234567890",
				Address:  "Kupang, NTT",
			},
			{
				Name:     "Admin Sumba",
				Username: "admin_sumba",
				Email:    "sumba@openpeo.id",
				Password: "password123",
				Role:     "admin",
				Phone:    "081234567891",
				Address:  "Waingapu, Sumba Timur, NTT",
			},
			{
				Name:     "Pembeli Flores",
				Username: "pembeli_flores",
				Email:    "flores@email.com",
				Password: "password123",
				Role:     "customer",
				Phone:    "081234567892",
				Address:  "Ende, Flores, NTT",
			},
			{
				Name:     "Pembeli Sumba",
				Username: "pembeli_sumba",
				Email:    "budi@email.com",
				Password: "password123",
				Role:     "customer",
				Phone:    "081234567893",
				Address:  "Jakarta, Indonesia",
			},
		}
		for i := range users {
			config.DB.Create(&users[i])
		}
		fmt.Printf("✅ Seeded %d users\n", len(users))
	}

	
	var productCount int64
	config.DB.Model(&models.Product{}).Count(&productCount)
	if productCount < 8 {
		fmt.Println("🌱 Seeding demo products...")

		
		var adminUser models.User
		if err := config.DB.Where("role = ?", "admin").First(&adminUser).Error; err != nil {
			if err := config.DB.First(&adminUser).Error; err != nil {
				fmt.Println("⚠️ Cannot seed products: no users found in database")
				return
			}
		}

		products := []models.Product{
			{
				VendorID:    adminUser.ID,
				Name:        "Kain Tenun Ikat Sumba Pahikung",
				Description: "Kain tenun ikat tradisional Sumba Timur bermotif Pahikung yang melambangkan kebangsawanan dan kekuatan. Dibuat dengan teknik tenun ikat pakan warisan turun-temurun menggunakan pewarna alami dari akar dan daun.",
				Price:       1500000,
				ImageURL:    "https://images.unsplash.com/photo-1558171813-4c088753af8f?w=400",
				Region:      "Sumba",
				Category:    "Tenun",
				MinOrder:    2,
				PODuration:  30,
				Stock:       15,
				IsActive:    true,
			},
			{
				VendorID:    adminUser.ID,
				Name:        "Selimut Tenun Sumba Hinggi Kombu",
				Description: "Selimut tenun Hinggi Kombu khas Sumba dengan motif hewan dan geometris. Simbol kehormatan dan status sosial dalam budaya Marapu.",
				Price:       2500000,
				ImageURL:    "https://images.unsplash.com/photo-1606722590583-6951b5ea92ad?w=400",
				Region:      "Sumba",
				Category:    "Tenun",
				MinOrder:    1,
				PODuration:  45,
				Stock:       8,
				IsActive:    true,
			},
			{
				VendorID:    adminUser.ID,
				Name:        "Songke Manggarai — Kain Adat",
				Description: "Kain tenun songke khas Manggarai dengan motif wela mpuu (bunga penuh). Ditenun secara tradisional oleh para perempuan Manggarai untuk upacara adat.",
				Price:       850000,
				ImageURL:    "https://images.unsplash.com/photo-1621600411688-4be93cd68504?w=400",
				Region:      "Manggarai",
				Category:    "Tenun",
				MinOrder:    3,
				PODuration:  21,
				Stock:       12,
				IsActive:    true,
			},
			{
				VendorID:    adminUser.ID,
				Name:        "Kopi Flores Bajawa Single Origin",
				Description: "Kopi arabika single origin dari dataran tinggi Bajawa, Flores. Ditanam di ketinggian 1200-1600 mdpl dengan cita rasa coklat, karamel, dan sentuhan rempah khas.",
				Price:       120000,
				ImageURL:    "https://images.unsplash.com/photo-1447933601403-0c6688de566e?w=400",
				Region:      "Flores",
				Category:    "Kuliner",
				MinOrder:    5,
				PODuration:  14,
				Stock:       50,
				IsActive:    true,
			},
			{
				VendorID:    adminUser.ID,
				Name:        "Manik-Manik Sumba Handmade Necklace",
				Description: "Kalung manik-manik handmade khas Sumba dengan warna-warna tanah. Setiap butir dibentuk dan diwarnai secara manual oleh pengrajin lokal.",
				Price:       350000,
				ImageURL:    "https://images.unsplash.com/photo-1611085583191-a3b181a88401?w=400",
				Region:      "Sumba",
				Category:    "Aksesoris",
				MinOrder:    3,
				PODuration:  14,
				Stock:       20,
				IsActive:    true,
			},
			{
				VendorID:    adminUser.ID,
				Name:        "Madu Hutan Timor Asli",
				Description: "Madu hutan asli dari pegunungan Timor, NTT. Dipanen langsung dari sarang lebah liar di hutan lindung dengan rasa yang kaya dan kental.",
				Price:       175000,
				ImageURL:    "https://images.unsplash.com/photo-1587049352846-4a222e784d38?w=400",
				Region:      "Timor",
				Category:    "Kuliner",
				MinOrder:    4,
				PODuration:  10,
				Stock:       35,
				IsActive:    true,
			},
			{
				VendorID:    adminUser.ID,
				Name:        "Patung Ukir Kayu Sandalwood Kupang",
				Description: "Patung ukiran kayu cendana (sandalwood) buatan tangan pengrajin Kupang. Motif tradisional NTT dengan aroma khas kayu cendana yang tahan lama.",
				Price:       950000,
				ImageURL:    "https://images.unsplash.com/photo-1513519245088-0e12902e35ca?w=400",
				Region:      "Kupang",
				Category:    "Kerajinan",
				MinOrder:    1,
				PODuration:  30,
				Stock:       5,
				IsActive:    true,
			},
			{
				VendorID:    adminUser.ID,
				Name:        "Syal Tenun Ende Lio",
				Description: "Syal tenun ikat Ende Lio dengan motif bunga dan geometris dalam warna-warna alam. Pewarna alami dari tumbuhan lokal daerah Ende.",
				Price:       450000,
				ImageURL:    "https://images.unsplash.com/photo-1601924921557-45e16393d8e1?w=400",
				Region:      "Flores",
				Category:    "Tenun",
				MinOrder:    2,
				PODuration:  21,
				Stock:       18,
				IsActive:    true,
			},
		}

		for i := range products {
			config.DB.Create(&products[i])
		}
		fmt.Printf("✅ Seeded %d products\n", len(products))
	}
}
