//go:build ignore

package main

import (
	"fmt"
	"log"
	"strings"

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

var kemejaImages = []string{
	"https://images.unsplash.com/photo-1598033129183-c4f50c736f10?q=80&w=600&auto=format&fit=crop",
	"https://images.unsplash.com/photo-1603252109303-2751441dd157?q=80&w=600&auto=format&fit=crop",
	"https://images.unsplash.com/photo-1589310261154-f58564f02758?q=80&w=600&auto=format&fit=crop",
}

var dressImages = []string{
	"https://images.unsplash.com/photo-1612336307429-8a898d10e223?q=80&w=600&auto=format&fit=crop",
	"https://images.unsplash.com/photo-1515372039744-b8f02a3ae446?q=80&w=600&auto=format&fit=crop",
	"https://images.unsplash.com/photo-1572804013309-27a88b4a11b6?q=80&w=600&auto=format&fit=crop",
}

var kopiImages = []string{
	"https://images.unsplash.com/photo-1559525839-b184a4d698c7?q=80&w=600&auto=format&fit=crop",
	"https://images.unsplash.com/photo-1511920170033-f8396924c348?q=80&w=600&auto=format&fit=crop",
	"https://images.unsplash.com/photo-1587734195503-904fca47e0e9?q=80&w=600&auto=format&fit=crop",
}

var makananImages = []string{
	"https://images.unsplash.com/photo-1587049352847-8d4c0b4c3116?q=80&w=600&auto=format&fit=crop",
	"https://images.unsplash.com/photo-1529193591184-b1d58069ecdd?q=80&w=600&auto=format&fit=crop",
	"https://images.unsplash.com/photo-1596040033229-a9821ebd058d?q=80&w=600&auto=format&fit=crop",
}

var aksesorisImages = []string{
	"https://images.unsplash.com/photo-1611591437281-460bfbe1220a?q=80&w=600&auto=format&fit=crop",
	"https://images.unsplash.com/photo-1535632066927-ab7c9ab60908?q=80&w=600&auto=format&fit=crop",
	"https://images.unsplash.com/photo-1623341214825-9f4f963727da?q=80&w=600&auto=format&fit=crop",
}

const fallbackImage = "https://images.unsplash.com/photo-1605814511559-009c95d90610?q=80&w=600&auto=format&fit=crop"

type bankResult struct {
	pool	[]string
	label	string
}

func classify(name string) bankResult {
	n := strings.ToLower(name)

	switch {
	case strings.Contains(n, "kemeja"):
		return bankResult{kemejaImages, "kemeja  "}

	case strings.Contains(n, "dress") ||
		strings.Contains(n, "blouse") ||
		strings.Contains(n, "outer") ||
		strings.Contains(n, "tunik") ||
		strings.Contains(n, "selendang"):
		return bankResult{dressImages, "dress   "}

	case strings.Contains(n, "kopi"):
		return bankResult{kopiImages, "kopi    "}

	case strings.Contains(n, "madu") ||
		strings.Contains(n, "sei") ||
		strings.Contains(n, "garam"):
		return bankResult{makananImages, "makanan "}

	case strings.Contains(n, "anting") ||
		strings.Contains(n, "gelang") ||
		strings.Contains(n, "headband"):
		return bankResult{aksesorisImages, "aksesoris"}

	default:
		return bankResult{nil, "fallback "}
	}
}

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
	fmt.Println("🔍 Fetching ALL products (diverse-image mode)...")

	var products []Product
	if err := db.Find(&products).Error; err != nil {
		log.Fatalf("❌ Failed to query products: %v\n", err)
	}
	fmt.Printf("   Found %d products\n\n", len(products))

	counts := map[string]int{}
	failed := 0

	for i := range products {
		p := &products[i]

		var imageURL string
		res := classify(p.Name)

		if res.pool != nil {
			idInt := int(p.ID)
			imageURL = res.pool[idInt%len(res.pool)]
		} else {
			imageURL = fallbackImage
		}

		result := db.Model(p).UpdateColumn("image_url", imageURL)
		if result.Error != nil {
			log.Printf("   [ERROR] #%d %s → %v\n", p.ID, p.Name, result.Error)
			failed++
			continue
		}

		counts[res.label]++
		fmt.Printf("   [%s] #%-3d %-42s → %s\n",
			res.label, p.ID, truncate(p.Name, 42), truncate(imageURL, 60))
	}

	fmt.Println("\n────────────────────────────────────────────────────────────")
	fmt.Printf("✅ Done. Errors: %d | Total updated: %d\n", failed, len(products)-failed)
	fmt.Println("   Breakdown by bank:")
	for label, count := range counts {
		fmt.Printf("     %-12s : %d\n", label, count)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

