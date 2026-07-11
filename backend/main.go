package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"openpeo-backend/config"
	"openpeo-backend/handlers"
	"openpeo-backend/middleware"
	"openpeo-backend/models"
)

func main() {
	// ---------------------------------------------------------------
	// 1. Database Initialization
	// ---------------------------------------------------------------
	config.ConnectDatabase()

	// Drop tables for clean start
	//fmt.Println("🧹 Dropping existing tables...")
	//config.DB.Migrator().DropTable(&models.OrderItem{}, &models.CartItem{}, &models.Order{}, &models.Message{}, &models.ActivityLog{}, &models.Product{}, &models.SellerProfile{}, &models.User{})

	// Auto-migrate all models to sync with the MySQL schema.
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Message{},
		&models.CartItem{},
		&models.ActivityLog{},
		&models.SellerProfile{},
		&models.WalletTransaction{},
	)
	if err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}
	fmt.Println("✅ Database schema synchronized")

	// Seed a demo vendor and customer if the users table is empty
	seedDemoData()

	// ---------------------------------------------------------------
	// 2. Gin Engine Setup
	// ---------------------------------------------------------------
	r := gin.Default()

	// CORS configuration for Vue.js dev server
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175", "http://localhost:3000", "http://127.0.0.1:5173", "http://127.0.0.1:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Serve static files from the ./images directory
	r.Static("/images", "./images")

	// ---------------------------------------------------------------
	// 3. RESTful API Routes
	// ---------------------------------------------------------------
	api := r.Group("/api")
	{
		// Authentication
		api.POST("/login", handlers.Login)
		api.POST("/register", handlers.RegisterUser)
		api.POST("/register-seller", handlers.RegisterSeller)
		api.GET("/user/:id", handlers.GetCurrentUser)

		// File Upload
		api.POST("/upload", handlers.UploadImage)

		// Payment (Midtrans Snap)
		api.POST("/checkout", handlers.HandleCheckout)

		// Product endpoints
		api.POST("/products", handlers.CreateProduct)
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProductByID)
		api.PUT("/products/:id", handlers.UpdateProduct)
		api.DELETE("/products/:id", handlers.DeleteProduct)

		// DEV tool
		api.POST("/dev/reset-data", handlers.DevResetData)

		// Protected endpoints
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Cart
			protected.POST("/cart", handlers.AddToCart)
			protected.GET("/cart", handlers.GetCart)
			protected.DELETE("/cart/:id", handlers.RemoveFromCart)

			// User Profile & Settings
			protected.GET("/user/profile", handlers.GetProfile)
			protected.PUT("/user/profile", handlers.UpdateProfile)
			protected.PUT("/user/password", handlers.ChangePassword)

			// Order endpoints
			protected.GET("/user/orders", handlers.GetMyOrders) // Legacy
			protected.GET("/orders/user", handlers.GetMyOrders) // Alias to match new requirements
			protected.POST("/orders", handlers.CreateOrder) // Cart checkout
			protected.POST("/orders/direct", handlers.CreateDirectOrder) // Buy It Now checkout
			protected.GET("/orders", handlers.GetOrders)
			protected.DELETE("/orders/:id", handlers.DeleteOrder)
			protected.PUT("/orders/:id/confirm-payment", handlers.ConfirmPayment)
			protected.PUT("/orders/:id/complete", handlers.CompleteOrder)
			
			// Seller Order Management (Moved here to explicitly match /api/orders/seller)
			protected.GET("/orders/seller", handlers.GetSellerOrders)
			protected.PUT("/orders/seller/:id/process", handlers.ProcessSellerOrder)
		}

		// Chat endpoints
		api.GET("/messages", handlers.GetChatHistory)
		api.GET("/chat/contacts", handlers.GetChatContacts)

		// Seller dashboard (protected for sellers)
		seller := api.Group("/seller")
		seller.Use(middleware.AuthMiddleware())
		// In a real app we'd have a RoleMiddleware, but we can do it in the handlers or via a simple check.
		{
			seller.GET("/dashboard/stats", handlers.GetSellerDashboardStats)
			seller.GET("/dashboard/chart", handlers.GetSellerDashboardChart)
			seller.GET("/profile", handlers.GetSellerProfile)
			seller.PUT("/profile", handlers.UpsertSellerProfile)
			seller.GET("/products", handlers.GetSellerProducts)
			seller.POST("/products", handlers.CreateSellerProduct)
			seller.PUT("/products/:id", handlers.UpdateSellerProduct)
			seller.DELETE("/products/:id", handlers.DeleteSellerProduct)

			seller.GET("/wallet", handlers.GetWalletDetails)
			seller.POST("/wallet/withdraw", handlers.WithdrawBalance)
		}

		// Admin Routes (Protected by Auth and Admin Middleware)
		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.AuthMiddleware(), middleware.RequireAdmin())
		{
			adminRoutes.GET("/stats", handlers.GetAdminStats)
			adminRoutes.GET("/users", handlers.GetAdminUsers)
			adminRoutes.GET("/products", handlers.GetAdminProducts)
			adminRoutes.DELETE("/products/:id", handlers.ForceDeleteProduct)
			adminRoutes.GET("/orders", handlers.GetAdminOrders)
			adminRoutes.PUT("/orders/:id/cancel", handlers.CancelOrder)
			adminRoutes.GET("/activity-logs", handlers.GetActivityLogs)
		}
	}

	// ---------------------------------------------------------------
	// 4. WebSocket Route
	// ---------------------------------------------------------------
	r.GET("/ws/chat", handlers.HandleWebSocket)

	// ---------------------------------------------------------------
	// 5. Health Check
	// ---------------------------------------------------------------
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "OpenPeo Backend Engine",
			"version": "1.0.0",
		})
	})

	// ---------------------------------------------------------------
	// 6. Start Server
	// ---------------------------------------------------------------
	fmt.Println("🚀 OpenPeo Backend Engine starting on :8081")
	fmt.Println("📡 REST API:    http://localhost:8081/api")
	fmt.Println("🔌 WebSocket:   ws://localhost:8081/ws/chat")
	fmt.Println("❤️  Health:      http://localhost:8081/health")

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}

// seedDemoData populates the database with sample users and products
// if the tables are empty. This allows the frontend to work immediately.
func seedDemoData() {
	var userCount int64
	config.DB.Model(&models.User{}).Count(&userCount)

	// 1. Seed users if users table is empty
	if userCount == 0 {
		fmt.Println("🌱 Seeding demo users...")
		users := []models.User{
			{
				Name:     "Admin Flores",
				Email:    "admin@openpeo.com",
				Password: "$2a$10$.4gRKHk/3kFjryRKYb9rq.li1V1pE08g0KtplV2/yzjpd7kAB3bte", // password123
				Role:     "admin",
				Phone:    "081234567890",
				Address:  "Kupang, NTT",
			},
			{
				Name:     "Pembeli Sumba",
				Email:    "budi@email.com",
				Password: "$2a$10$.4gRKHk/3kFjryRKYb9rq.li1V1pE08g0KtplV2/yzjpd7kAB3bte", // password123
				Role:     "customer",
				Phone:    "081234567893",
				Address:  "Jakarta, Indonesia",
			},
			{
				Name:      "Toko Tenun NTT",
				Email:     "tenun@openpeo.com",
				Password:  "$2a$10$.4gRKHk/3kFjryRKYb9rq.li1V1pE08g0KtplV2/yzjpd7kAB3bte",
				Role:      "seller",
				Phone:     "081234567894",
				Address:   "Sumba, NTT",
			},
			{
				Name:      "Toko Kuliner Lokal",
				Email:     "kuliner@openpeo.com",
				Password:  "$2a$10$.4gRKHk/3kFjryRKYb9rq.li1V1pE08g0KtplV2/yzjpd7kAB3bte",
				Role:      "seller",
				Phone:     "081234567895",
				Address:   "Kupang, NTT",
			},
			{
				Name:      "Toko Aksesoris Flores",
				Email:     "aksesoris@openpeo.com",
				Password:  "$2a$10$.4gRKHk/3kFjryRKYb9rq.li1V1pE08g0KtplV2/yzjpd7kAB3bte",
				Role:      "seller",
				Phone:     "081234567896",
				Address:   "Flores, NTT",
			},
		}
		for i := range users {
			config.DB.Create(&users[i])
		}
		fmt.Printf("✅ Seeded %d users\n", len(users))
	}

	// 2. Seed products if products table is empty or missing items
	var productCount int64
	config.DB.Model(&models.Product{}).Count(&productCount)
	if productCount == 0 {
		fmt.Println("🌱 Seeding demo products...")

		// Find sellers to assign products to
		var sellerTenun, sellerKuliner, sellerAksesoris models.User
		config.DB.Where("email = ?", "tenun@openpeo.com").First(&sellerTenun)
		config.DB.Where("email = ?", "kuliner@openpeo.com").First(&sellerKuliner)
		config.DB.Where("email = ?", "aksesoris@openpeo.com").First(&sellerAksesoris)

		products := []models.Product{
			// SECTION 1: Koleksi Tenun NTT (Seller 1)
			{SellerID: sellerTenun.ID, Name: "Kaftan Indigo Sumba", Price: 2850000, Category: "Koleksi Tenun NTT", ImageURL: "images/product-1.png", Region: "Sumba", Description: "Kaftan tenun ikat berwarna indigo biru tua dengan motif geometris putih dari Sumba. Ditenun secara tradisional memakan waktu berbulan-bulan, menghasilkan karya seni yang elegan dan otentik.", PreOrderDuration: 14, IsActive: true},
			{SellerID: sellerTenun.ID, Name: "Selendang Merah Sabu", Price: 1950000, Category: "Koleksi Tenun NTT", ImageURL: "images/product-2.png", Region: "Sabu", Description: "Selendang tenun berwarna merah coral dengan garis-garis halus dan motif tradisional Sabu. Sangat cocok dipadukan dengan pakaian modern maupun tradisional.", PreOrderDuration: 14, IsActive: true},
			{SellerID: sellerTenun.ID, Name: "Syal Tenun Amarasi", Price: 1250000, Category: "Koleksi Tenun NTT", ImageURL: "images/product-3.png", Region: "Amarasi", Description: "Syal tenun ikat berwarna coklat dan krem dengan motif tradisional Amarasi. Lembut dan hangat, memberikan sentuhan etnik pada setiap penampilan.", PreOrderDuration: 7, IsActive: true},
			{SellerID: sellerTenun.ID, Name: "Tunik Hitam Rote", Price: 3150000, Category: "Koleksi Tenun NTT", ImageURL: "images/product-4.png", Region: "Rote", Description: "Tunik tenun ikat berwarna hitam dengan motif emas dan krem dari Pulau Rote. Desain timeless dengan kualitas benang premium terbaik.", PreOrderDuration: 21, IsActive: true},
			{SellerID: sellerTenun.ID, Name: "Outer Tenun Ikat Sumba", Price: 1450000, Category: "Koleksi Tenun NTT", ImageURL: "images/tenun-outer.png", Region: "Sumba", Description: "Mannequin/model wearing a chic, modern-cut outer featuring bold Sumba horse motifs. Elegan untuk acara formal maupun semi-formal.", PreOrderDuration: 14, IsActive: true},
			{SellerID: sellerTenun.ID, Name: "Kemeja Tenun Ende Pria", Price: 950000, Category: "Koleksi Tenun NTT", ImageURL: "images/tenun-kemeja.png", Region: "Ende", Description: "Professional man's shirt with classic, geometric Ende patterns. Memadukan kesan profesional dan tradisi, pas untuk gaya kantor modern.", PreOrderDuration: 10, IsActive: true},
			{SellerID: sellerTenun.ID, Name: "Dress Tenun Manggarai", Price: 2100000, Category: "Koleksi Tenun NTT", ImageURL: "images/tenun-dress.png", Region: "Manggarai", Description: "Elegant women's dress with premium, detailed Manggarai weave patterns. Siluet yang menawan untuk acara spesial.", PreOrderDuration: 21, IsActive: true},
			{SellerID: sellerTenun.ID, Name: "Blouse Tenun Alor", Price: 1150000, Category: "Koleksi Tenun NTT", ImageURL: "images/tenun-blouse.png", Region: "Alor", Description: "Stylish blouse featuring the distinct, colorful Alor textile designs. Nyaman dikenakan sepanjang hari dengan potongan yang modern.", PreOrderDuration: 14, IsActive: true},

			// SECTION 2: Cita Rasa Lokal (Seller 2)
			{SellerID: sellerKuliner.ID, Name: "Sei Babi Asap Kupang", Price: 185000, Category: "Cita Rasa Lokal", ImageURL: "images/sei-babi.png", Region: "Kupang", Description: "Sei babi asap tradisional dari Kupang dengan irisan daging berwarna amber keemasan. Diproses dengan kayu kosambi yang memberikan aroma khas.", PreOrderDuration: 3, IsActive: true},
			{SellerID: sellerKuliner.ID, Name: "Madu Hutan Timor Murni", Price: 145000, Category: "Cita Rasa Lokal", ImageURL: "images/madu-hutan.png", Region: "Timor", Description: "Toples madu hutan murni dari Timor berwarna emas tua. Diperoleh langsung dari hutan belantara, kaya akan nutrisi dan antioksidan.", PreOrderDuration: 3, IsActive: true},
			{SellerID: sellerKuliner.ID, Name: "Kopi Arabika Flores Bajawa", Price: 95000, Category: "Cita Rasa Lokal", ImageURL: "images/kopi-flores.png", Region: "Flores", Description: "Kemasan kopi arabika Flores Bajawa dengan biji kopi roasting gelap. Memiliki hints cokelat dan rempah yang kuat dengan acidity medium.", PreOrderDuration: 5, IsActive: true},
			{SellerID: sellerKuliner.ID, Name: "Garam Gunung Organik Khas NTT", Price: 65000, Category: "Cita Rasa Lokal", ImageURL: "images/garam-gunung.png", Region: "NTT", Description: "Garam gunung organik NTT berwarna merah muda keabu-abuan. Tinggi mineral alami dan sangat cocok untuk fine dining.", PreOrderDuration: 5, IsActive: true},
			{SellerID: sellerKuliner.ID, Name: "Jagung Titi Khas NTT", Price: 35000, Category: "Cita Rasa Lokal", ImageURL: "images/jagung-titi.png", Region: "NTT", Description: "Close-up of crispy, flat traditional corn snack in a premium rustic wooden bowl. Camilan ringan dan sehat untuk teman minum kopi.", PreOrderDuration: 2, IsActive: true},
			{SellerID: sellerKuliner.ID, Name: "Sambal Luat Kupang", Price: 45000, Category: "Cita Rasa Lokal", ImageURL: "images/sambal-luat.png", Region: "Kupang", Description: "Macro shot of authentic, spicy red Sambal Luat in a small glass jar with fresh chili garnishes. Pedas yang menyegarkan dengan sentuhan jeruk nipis.", PreOrderDuration: 2, IsActive: true},
			{SellerID: sellerKuliner.ID, Name: "Kacang Sembunyi NTT", Price: 30000, Category: "Cita Rasa Lokal", ImageURL: "images/kacang-sembunyi.png", Region: "NTT", Description: "Pile of crunchy, golden-brown caramelized nuts, studio lighting, clean background. Gurih manis yang adiktif.", PreOrderDuration: 2, IsActive: true},
			{SellerID: sellerKuliner.ID, Name: "Gula Air Sabu", Price: 55000, Category: "Cita Rasa Lokal", ImageURL: "images/gula-air.png", Region: "Sabu", Description: "Close-up of authentic, thick palm nectar in a clear bottle, honey-colored glow. Pemanis alami rendah indeks glikemik.", PreOrderDuration: 3, IsActive: true},

			// SECTION 3: Koleksi Aksesoris (Seller 3)
			{SellerID: sellerAksesoris.ID, Name: "Headband Kain Tenun Sumba", Price: 275000, Category: "Koleksi Aksesoris", ImageURL: "images/headband-tenun.png", Region: "Sumba", Description: "Headband kain tenun ikat Sumba berwarna indigo biru dengan motif geometris krem. Aksesori chic yang mudah dipadupadankan.", PreOrderDuration: 7, IsActive: true},
			{SellerID: sellerAksesoris.ID, Name: "Gelang & Kalung Serbuk Gading Maumere", Price: 450000, Category: "Koleksi Aksesoris", ImageURL: "images/gelang-kalung.png", Region: "Maumere", Description: "Set gelang dan kalung artisan dari Maumere dengan manik-manik serbuk gading dan kuningan. Statement piece untuk berbagai gaya.", PreOrderDuration: 14, IsActive: true},
			{SellerID: sellerAksesoris.ID, Name: "Mahkota Ti'i Langga Rote", Price: 850000, Category: "Koleksi Aksesoris", ImageURL: "images/mahkota-tiilangga.png", Region: "Rote", Description: "Mahkota Ti'i Langga tradisional dari Rote yang terbuat dari anyaman daun lontar. Sebuah karya seni ikonis yang bersejarah.", PreOrderDuration: 21, IsActive: true},
			{SellerID: sellerAksesoris.ID, Name: "Cincin & Aksesoris Ornamen Penyu Tradisional", Price: 320000, Category: "Koleksi Aksesoris", ImageURL: "images/cincin-penyu.png", Region: "NTT", Description: "Koleksi cincin dan aksesoris ornamen penyu tradisional NTT dari bahan alami yang lestari.", PreOrderDuration: 7, IsActive: true},
			{SellerID: sellerAksesoris.ID, Name: "Kalung Khas Timor", Price: 150000, Category: "Koleksi Aksesoris", ImageURL: "images/kalung-timor.jpg", Region: "Timor", Description: "Kalung manik-manik khas Timor dengan untaian manik oranye cerah dan ornamen hitam-putih bermotif etnik. Karya tangan perajin lokal yang teliti.", PreOrderDuration: 7, IsActive: true},
			{SellerID: sellerAksesoris.ID, Name: "Sisir Adat Sumba", Price: 225000, Category: "Koleksi Aksesoris", ImageURL: "images/sisir-sumba.png", Region: "Sumba", Description: "Intricate, hand-carved traditional wooden hair comb with traditional motifs. Aksesori rambut eksotis khas pulau Sumba.", PreOrderDuration: 10, IsActive: true},
			{SellerID: sellerAksesoris.ID, Name: "Gelang Kulit Kerbau Rote", Price: 85000, Category: "Koleksi Aksesoris", ImageURL: "images/gelang-kerbau.png", Region: "Rote", Description: "Detailed shot of handcrafted leather and bead bracelet. Tahan lama dan bertekstur klasik.", PreOrderDuration: 5, IsActive: true},
			{SellerID: sellerAksesoris.ID, Name: "Anting Motif Penyu Flores", Price: 120000, Category: "Koleksi Aksesoris", ImageURL: "images/anting-penyu.png", Region: "Flores", Description: "Elegant metalwork earrings with traditional sea-turtle-inspired motifs. Menambah kesan ayu pada penampilan elegan Anda.", PreOrderDuration: 7, IsActive: true},
		}

		for i := range products {
			config.DB.Create(&products[i])
		}
		fmt.Printf("✅ Seeded %d products\n", len(products))
	}

	// 3. Seed demo orders so the profile page has order history to show
	var orderCount int64
	config.DB.Model(&models.Order{}).Count(&orderCount)
	if orderCount == 0 {
		fmt.Println("🌱 Seeding demo orders...")

		var customer models.User
		if err := config.DB.Where("role = ?", "customer").First(&customer).Error; err == nil {
			demoOrders := []models.Order{
				{
					CustomerID:  customer.ID,
					SellerID:    1,
					OrderItems: []models.OrderItem{
						{ProductID: 1, Quantity: 1, Price: 2850000},
					},
					Quantity:    1,
					TotalPrice:  2850000,
					Status:      "Diproses Perajin",
					Note:        "Kaftan Indigo Sumba",
					CustomNotes: "Ukuran custom: Panjang 120cm, Lebar 60cm",
				},
				{
					CustomerID: customer.ID,
					SellerID:   1,
					OrderItems: []models.OrderItem{
						{ProductID: 11, Quantity: 3, Price: 95000},
					},
					Quantity:   3,
					TotalPrice: 285000,
					Status:     "Dikirim",
					Note:       "Kopi Arabika Flores Bajawa",
				},
				{
					CustomerID: customer.ID,
					SellerID:   1,
					OrderItems: []models.OrderItem{
						{ProductID: 19, Quantity: 1, Price: 850000},
					},
					Quantity:   1,
					TotalPrice: 850000,
					Status:     "Selesai",
					Note:       "Mahkota Ti'i Langga Rote",
				},
				{
					CustomerID:  customer.ID,
					SellerID:    1,
					OrderItems: []models.OrderItem{
						{ProductID: 4, Quantity: 1, Price: 3150000},
					},
					Quantity:    1,
					TotalPrice:  3150000,
					Status:      "Menunggu Pembayaran",
					Note:        "Tunik Hitam Rote",
					CustomNotes: "Request ukuran XL, lengan 3/4",
				},
			}

			for i := range demoOrders {
				config.DB.Create(&demoOrders[i])
			}
			fmt.Printf("✅ Seeded %d demo orders\n", len(demoOrders))
		}
	}
}
