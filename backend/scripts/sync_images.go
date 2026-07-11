//go:build ignore

// sync_images.go — Local Image Sync v1
// Maps every product in the `products` table to its corresponding physical
// file in frontend/public/images/ using keyword matching on the product name.
//
// Image path format stored in DB: /images/{filename}
// (served statically by the Vue dev server from the /public directory)
//
// Matching strategy:
//  1. A prioritised keyword map is checked against the lowercased product name.
//  2. The FIRST matching rule wins (longest/most-specific keywords are listed first).
//  3. If no keyword matches, the product is skipped and reported — no overwrite.
//
// Usage (run from the /backend directory):
//
//	go run scripts/sync_images.go
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
// Keyword → local image filename mapping
// Rules are evaluated top-to-bottom; the first match wins.
// Specific / multi-word patterns are listed before generic single-word ones.
// ---------------------------------------------------------------------------

type rule struct {
	keyword  string // checked with strings.Contains on lowercased product name
	filename string // relative filename only; script prepends "/images/"
}

var rules = []rule{
	// ── Accessories ──────────────────────────────────────────────────────────
	{keyword: "anting", filename: "anting-penyu.png"},
	{keyword: "cincin", filename: "cincin-penyu.png"},
	{keyword: "gelang & kalung", filename: "gelang-kalung.png"}, // before plain "gelang"
	{keyword: "gelang", filename: "gelang-kerbau.png"},
	{keyword: "kalung", filename: "kalung-timor.jpg"},
	{keyword: "headband", filename: "headband-tenun.png"},
	{keyword: "mahkota", filename: "mahkota-tiilangga.png"},
	{keyword: "sisir", filename: "sisir-sumba.png"},

	// ── Food & Beverage ──────────────────────────────────────────────────────
	{keyword: "kopi", filename: "kopi-flores.png"},
	{keyword: "madu", filename: "madu-hutan.png"},
	{keyword: "sei", filename: "sei-babi.png"},
	{keyword: "garam", filename: "garam-gunung.png"},
	{keyword: "jagung", filename: "jagung-titi.png"},
	{keyword: "kacang", filename: "kacang-sembunyi.png"},
	{keyword: "sambal", filename: "sambal-luat.png"},
	{keyword: "gula", filename: "gula-air.png"},

	// ── Fashion / Tenun ──────────────────────────────────────────────────────
	{keyword: "kemeja", filename: "tenun-kemeja.png"},
	{keyword: "dress", filename: "tenun-dress.png"},
	{keyword: "blouse", filename: "tenun-blouse.png"},
	{keyword: "outer", filename: "tenun-outer.png"},
	// tunik, selendang, kaftan, syal, kameja → fallback to numbered product images
	{keyword: "tunik", filename: "product-1.png"},
	{keyword: "selendang", filename: "product-2.png"},
	{keyword: "kaftan", filename: "product-3.png"},
	{keyword: "syal", filename: "product-4.png"},
	{keyword: "kameja", filename: "tenun-kemeja.png"}, // typo variant of "kemeja"
}

// resolve returns the /images/... path for the given product name,
// and ok=false if no rule matched.
func resolve(name string) (path string, ok bool) {
	n := strings.ToLower(name)
	for _, r := range rules {
		if strings.Contains(n, r.keyword) {
			return "/images/" + r.filename, true
		}
	}
	return "", false
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
	fmt.Println("🔗 Syncing local image files to products table...")

	var products []Product
	if err := db.Find(&products).Error; err != nil {
		log.Fatalf("❌ Failed to query products: %v\n", err)
	}
	fmt.Printf("   Found %d products\n\n", len(products))

	updated, skipped, failed := 0, 0, 0

	for i := range products {
		p := &products[i]

		imagePath, ok := resolve(p.Name)
		if !ok {
			fmt.Printf("   [SKIP   ] #%-3d %s  ← no matching image file\n", p.ID, p.Name)
			skipped++
			continue
		}

		result := db.Model(p).UpdateColumn("image_url", imagePath)
		if result.Error != nil {
			log.Printf("   [ERROR  ] #%d %s → %v\n", p.ID, p.Name, result.Error)
			failed++
			continue
		}

		fmt.Printf("   [UPDATED] #%-3d %-45s → %s\n",
			p.ID, truncate(p.Name, 45), imagePath)
		updated++
	}

	fmt.Println("\n────────────────────────────────────────────────────────────")
	fmt.Printf("✅ Done.  Updated: %d | Skipped: %d | Errors: %d | Total: %d\n",
		updated, skipped, failed, len(products))
}

// truncate shortens s to at most n chars, appending "…" if needed.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
