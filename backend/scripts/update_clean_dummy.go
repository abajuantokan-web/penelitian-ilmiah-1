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

const (
	colorBlue	= "3B82F6"
	colorAmber	= "F59E0B"
	colorPurple	= "8B5CF6"
	colorGrey	= "E5E7EB"
)

func categoryColor(name string) (hex, label string) {
	n := strings.ToLower(name)

	switch {
	case strings.Contains(n, "kemeja") ||
		strings.Contains(n, "dress") ||
		strings.Contains(n, "blouse") ||
		strings.Contains(n, "tunik") ||
		strings.Contains(n, "outer") ||
		strings.Contains(n, "selendang"):
		return colorBlue, "fashion  "

	case strings.Contains(n, "kopi") ||
		strings.Contains(n, "madu") ||
		strings.Contains(n, "sei") ||
		strings.Contains(n, "garam"):
		return colorAmber, "food     "

	case strings.Contains(n, "anting") ||
		strings.Contains(n, "gelang") ||
		strings.Contains(n, "headband"):
		return colorPurple, "accessory"

	default:
		return colorGrey, "fallback "
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
	fmt.Println("🎨 Applying clean solid-color placeholder images...")

	var products []Product
	if err := db.Find(&products).Error; err != nil {
		log.Fatalf("❌ Failed to query products: %v\n", err)
	}
	fmt.Printf("   Found %d products\n\n", len(products))

	counts := map[string]int{}
	failed := 0

	for i := range products {
		p := &products[i]

		hex, label := categoryColor(p.Name)
		dummyURL := fmt.Sprintf(
			"https://placehold.co/600x600/%s/%s?text=%%20",
			hex, hex,
		)

		result := db.Model(p).UpdateColumn("image_url", dummyURL)
		if result.Error != nil {
			log.Printf("   [ERROR] #%d %s → %v\n", p.ID, p.Name, result.Error)
			failed++
			continue
		}

		counts[label]++
		fmt.Printf("   [%s] #%-3d %-45s → #%s\n",
			label, p.ID, truncate(p.Name, 45), hex)
	}

	fmt.Println("\n────────────────────────────────────────────────────────────")
	fmt.Printf("✅ Done. Errors: %d | Total updated: %d\n", failed, len(products)-failed)
	fmt.Println("   Color breakdown:")
	colorMap := map[string]string{
		"fashion  ":	colorBlue,
		"food     ":	colorAmber,
		"accessory":	colorPurple,
		"fallback ":	colorGrey,
	}
	for label, count := range counts {
		fmt.Printf("     %-12s #%s : %d products\n", label, colorMap[label], count)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

