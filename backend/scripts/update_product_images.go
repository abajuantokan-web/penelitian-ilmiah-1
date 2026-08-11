//go:build ignore

package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Product struct {
	ID		int32	`gorm:"primaryKey"`
	Name		string	`gorm:"size:200"`
	Category	string	`gorm:"size:100"`
	ImageURL	string	`gorm:"column:image_url;size:512"`
}

func (Product) TableName() string	{ return "products" }

var fashionImages = []string{

	"https://images.unsplash.com/photo-1605814511559-009c95d90610?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1589547562479-79a405fbba76?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1528813860492-bb99459b68a7?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1572804013309-27a88b4a11b6?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1620799140408-edc6dcb6d633?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1504198266287-1659872e6590?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1602810318383-e386cc2a3ccf?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1485462537746-965f33f7f6a7?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1603400521630-9f2de124b33b?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1558769132-cb1aea458c5e?q=80&w=800&auto=format&fit=crop",
}

var foodImages = []string{

	"https://images.unsplash.com/photo-1559525839-b184a4d698c7?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1587049352847-8d4c0b4c3116?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1596040033229-a9821ebd058d?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1529193591184-b1d58069ecdd?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1558703083-61c5591e7c2f?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1447933601403-0c6688de566e?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1563636619-e9143da7973b?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1610832958506-aa56368176cf?q=80&w=800&auto=format&fit=crop",
}

var accessoryImages = []string{

	"https://images.unsplash.com/photo-1611591437281-460bfbe1220a?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1535632066927-ab7c9ab60908?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1623341214825-9f4f963727da?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1630350040792-3d3d71e4d8bf?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1605100804763-247f67b3557e?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1600950207944-0d63e8edbc3f?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1566479179817-c8b1386dbd1c?q=80&w=800&auto=format&fit=crop",

	"https://images.unsplash.com/photo-1589128777073-263566ae5e4d?q=80&w=800&auto=format&fit=crop",
}

var foodKeywords = []string{
	"kopi", "madu", "sei", "garam", "jagung", "sambal",
	"kacang", "gula", "kuliner", "rempah", "ikan", "daging",
}

var accessoryKeywords = []string{
	"anting", "gelang", "headband", "kalung", "cincin",
	"mahkota", "sisir", "bros", "aksesoris", "perhiasan",
}

var foodCategories = []string{"cita rasa", "kuliner", "makanan", "minuman"}

var accessoryCategories = []string{"aksesoris", "accessory", "perhiasan", "jewelry"}

func classify(name, category string) []string {
	nameLower := strings.ToLower(name)
	catLower := strings.ToLower(category)

	for _, fc := range foodCategories {
		if strings.Contains(catLower, fc) {
			return foodImages
		}
	}
	for _, ac := range accessoryCategories {
		if strings.Contains(catLower, ac) {
			return accessoryImages
		}
	}

	for _, kw := range foodKeywords {
		if strings.Contains(nameLower, kw) {
			return foodImages
		}
	}
	for _, kw := range accessoryKeywords {
		if strings.Contains(nameLower, kw) {
			return accessoryImages
		}
	}

	return fashionImages
}

func pickRandom(pool []string) string {
	return pool[rand.Intn(len(pool))]
}

func main() {

	rand.New(rand.NewSource(time.Now().UnixNano()))

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
	fmt.Println("🔍 Fetching ALL products (no filter — overwrite mode)...")

	var products []Product
	if err := db.Find(&products).Error; err != nil {
		log.Fatalf("❌ Failed to query products: %v\n", err)
	}
	fmt.Printf("   Found %d products\n\n", len(products))

	food, accessory, fashion, failed := 0, 0, 0, 0

	for i := range products {
		p := &products[i]

		pool := classify(p.Name, p.Category)
		newURL := pickRandom(pool)

		label := "fashion"
		switch &pool[0] {
		case &foodImages[0]:
			label = "food    "
			food++
		case &accessoryImages[0]:
			label = "accessory"
			accessory++
		default:
			label = "fashion  "
			fashion++
		}

		result := db.Model(p).UpdateColumn("image_url", newURL)
		if result.Error != nil {
			log.Printf("   [ERROR] #%d %s → %v\n", p.ID, p.Name, result.Error)
			failed++
			continue
		}

		fmt.Printf("   [%s] #%-3d %-42s → %s\n",
			label, p.ID, truncate(p.Name, 42), truncate(newURL, 55))
	}

	fmt.Printf("\n────────────────────────────────────────────────────────────\n")
	fmt.Printf("✅ Done.  Fashion: %d | Food: %d | Accessory: %d | Errors: %d | Total: %d\n",
		fashion, food, accessory, failed, len(products))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

