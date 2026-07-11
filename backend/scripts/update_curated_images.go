//go:build ignore

// update_curated_images.go — Contextual Curated Image Seeder v4
// Assigns highly relevant, ethnic/rustic Unsplash images to every product
// using keyword matching on product name + ID-based modulo cycling for variety.
//
// Photo Banks:
//   tenunImages    → woven fabrics, artisan weaving, ethnic textiles
//   kopiImages     → rustic coffee beans / brewing scenes
//   makananImages  → honey jars, smoked meat, artisan spices
//   aksesorisImages → ethnic bracelets, artisan jewelry
//
// Usage (run from the /backend directory):
//
//	go run scripts/update_curated_images.go
package main

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ---------------------------------------------------------------------------
// Minimal Product model
// ---------------------------------------------------------------------------

type Product struct {
	ID       int32  `gorm:"primaryKey"`
	Name     string `gorm:"size:200"`
	Category string `gorm:"size:100"`
	ImageURL string `gorm:"column:image_url;size:512"`
}

func (Product) TableName() string { return "products" }

// ---------------------------------------------------------------------------
// Photo Banks — curated for ethnic / rustic Indonesian vibe
// ---------------------------------------------------------------------------

var tenunImages = []string{
	"https://images.unsplash.com/photo-1605814511559-009c95d90610?q=80&w=600&auto=format&fit=crop", // Woven fabric texture
	"https://images.unsplash.com/photo-1528813860492-bb99459b68a7?q=80&w=600&auto=format&fit=crop", // Artisan weaving on loom
	"https://images.unsplash.com/photo-1589547562479-79a405fbba76?q=80&w=600&auto=format&fit=crop", // Ethnic draped textile
}

var kopiImages = []string{
	"https://images.unsplash.com/photo-1559525839-b184a4d698c7?q=80&w=600&auto=format&fit=crop", // Coffee beans in wooden scoop
	"https://images.unsplash.com/photo-1511920170033-f8396924c348?q=80&w=600&auto=format&fit=crop", // Rustic coffee brewing
}

var makananImages = []string{
	"https://images.unsplash.com/photo-1587049352847-8d4c0b4c3116?q=80&w=600&auto=format&fit=crop", // Rustic honey jar
	"https://images.unsplash.com/photo-1529193591184-b1d58069ecdd?q=80&w=600&auto=format&fit=crop", // Smoked meat / spices
}

var aksesorisImages = []string{
	"https://images.unsplash.com/photo-1611591437281-460bfbe1220a?q=80&w=600&auto=format&fit=crop", // Ethnic beaded bracelet
	"https://images.unsplash.com/photo-1535632066927-ab7c9ab60908?q=80&w=600&auto=format&fit=crop", // Artisan jewelry on velvet
}

// ---------------------------------------------------------------------------
// classify returns the right photo bank and label for a product name
// ---------------------------------------------------------------------------

type result struct {
	pool  []string
	label string
}

func classify(name string) result {
	n := strings.ToLower(name)

	switch {
	case strings.Contains(n, "kemeja") ||
		strings.Contains(n, "dress") ||
		strings.Contains(n, "blouse") ||
		strings.Contains(n, "outer") ||
		strings.Contains(n, "tunik") ||
		strings.Contains(n, "selendang"):
		return result{tenunImages, "tenun    "}

	case strings.Contains(n, "kopi"):
		return result{kopiImages, "kopi     "}

	case strings.Contains(n, "madu") ||
		strings.Contains(n, "sei") ||
		strings.Contains(n, "garam"):
		return result{makananImages, "makanan  "}

	case strings.Contains(n, "anting") ||
		strings.Contains(n, "gelang") ||
		strings.Contains(n, "headband"):
		return result{aksesorisImages, "aksesoris"}

	default:
		// Fallback: woven fabric — fits the general NTT artisan theme
		return result{tenunImages, "fallback "}
	}
}

// ---------------------------------------------------------------------------
// Entry point
// ---------------------------------------------------------------------------

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
	fmt.Println("🖼️  Applying contextual curated images...")

	var products []Product
	if err := db.Find(&products).Error; err != nil {
		log.Fatalf("❌ Failed to query products: %v\n", err)
	}
	fmt.Printf("   Found %d products\n\n", len(products))

	counts := map[string]int{}
	failed := 0

	for i := range products {
		p := &products[i]

		res := classify(p.Name)
		idInt := int(p.ID)
		imageURL := res.pool[idInt%len(res.pool)]

		result := db.Model(p).UpdateColumn("image_url", imageURL)
		if result.Error != nil {
			log.Printf("   [ERROR] #%d %s → %v\n", p.ID, p.Name, result.Error)
			failed++
			continue
		}

		counts[res.label]++
		fmt.Printf("   [%s] #%-3d %-45s → %s\n",
			res.label, p.ID, truncate(p.Name, 45), truncate(imageURL, 55))
	}

	fmt.Println("\n────────────────────────────────────────────────────────────")
	fmt.Printf("✅ Done. Errors: %d | Total updated: %d\n", failed, len(products)-failed)
	fmt.Println("   Breakdown by bank:")
	for label, count := range counts {
		fmt.Printf("     %-12s : %d products\n", label, count)
	}
}

// truncate shortens s to at most n chars, appending "…" if needed.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
