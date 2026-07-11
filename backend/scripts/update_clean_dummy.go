//go:build ignore

// update_clean_dummy.go — Clean Solid-Color Dummy Seeder v1
// Assigns a perfectly blank, category-colored 600×600 placeholder image
// to every product using placehold.co with `?text=%20` to suppress all text.
//
// Color mapping:
//   - Kemeja / Dress / Blouse / Tunik / Outer / Selendang → #3B82F6 (Blue)
//   - Kopi / Madu / Sei / Garam / Makanan                 → #F59E0B (Amber)
//   - Anting / Gelang / Headband / Aksesoris              → #8B5CF6 (Purple)
//   - Everything else                                      → #E5E7EB (Light Grey)
//
// Usage (run from the /backend directory):
//
//	go run scripts/update_clean_dummy.go
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
// Minimal Product model — only the columns this script needs.
// ---------------------------------------------------------------------------

type Product struct {
	ID       int32  `gorm:"primaryKey"`
	Name     string `gorm:"size:200"`
	Category string `gorm:"size:100"`
	ImageURL string `gorm:"column:image_url;size:512"`
}

func (Product) TableName() string { return "products" }

// ---------------------------------------------------------------------------
// Color palette — solid HEX codes (no leading #, as placehold.co expects)
// ---------------------------------------------------------------------------

const (
	colorBlue   = "3B82F6" // Fashion: kemeja, dress, blouse, tunik, outer, selendang
	colorAmber  = "F59E0B" // Food & Beverage: kopi, madu, sei, garam
	colorPurple = "8B5CF6" // Accessories: anting, gelang, headband
	colorGrey   = "E5E7EB" // Fallback / unclassified
)

// categoryColor returns a HEX color code and a human-readable label
// based on keywords found in the product name.
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
	fmt.Println("🎨 Applying clean solid-color placeholder images...")

	var products []Product
	if err := db.Find(&products).Error; err != nil {
		log.Fatalf("❌ Failed to query products: %v\n", err)
	}
	fmt.Printf("   Found %d products\n\n", len(products))

	// Counters per category label
	counts := map[string]int{}
	failed := 0

	for i := range products {
		p := &products[i]

		// ── Determine color & build placehold.co URL ───────────────────────
		// ?text=%20 forces placehold.co to render a blank space instead of
		// the default "600x600" dimension text → perfectly clean solid square.
		hex, label := categoryColor(p.Name)
		dummyURL := fmt.Sprintf(
			"https://placehold.co/600x600/%s/%s?text=%%20",
			hex, hex,
		)

		// ── Persist ────────────────────────────────────────────────────────
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
		"fashion  ": colorBlue,
		"food     ": colorAmber,
		"accessory": colorPurple,
		"fallback ": colorGrey,
	}
	for label, count := range counts {
		fmt.Printf("     %-12s #%s : %d products\n", label, colorMap[label], count)
	}
}

// truncate shortens s to at most n chars, appending "…" if needed.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
