//go:build ignore

// fix_image_mapping.go — Exact Image Mapping Fix v2
// Resets image_url for ALL products using ORDERED keyword rules so that
// longer/more-specific keywords are always checked before shorter ones —
// avoiding the random-iteration problem of plain Go maps.
//
// Mapping (keyword → /images/ filename):
//
//	anting       → anting-penyu.png
//	cincin       → cincin-penyu.png
//	garam        → garam-gunung.png
//	gelang & kal → gelang-kalung.png  (checked before plain "gelang")
//	gelang       → gelang-kerbau.png
//	gula         → gula-air.png
//	headband     → headband-tenun.png
//	jagung       → jagung-titi.png
//	kacang       → kacang-sembunyi.png
//	kalung       → kalung-timor.jpg
//	kopi         → kopi-flores.png
//	madu         → madu-hutan.png
//	mahkota      → mahkota-tiilangga.png
//	sambal       → sambal-luat.png
//	sei          → sei-babi.png
//	sisir        → sisir-sumba.png
//	kemeja/kamej → tenun-kemeja.png
//	dress        → tenun-dress.png
//	blouse       → tenun-blouse.png
//	outer        → tenun-outer.png
//	tunik/kafta/ → product-1.png … (remaining tenun variants)
//	(fallback)   → product-1.png
//
// Usage (run from the /backend directory):
//
//	go run scripts/fix_image_mapping.go
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
// Ordered rules — evaluated top-to-bottom; first match wins.
// IMPORTANT: more-specific / multi-word patterns must come before their
//            shorter sub-strings (e.g. "gelang & kalung" before "gelang").
// ---------------------------------------------------------------------------

type rule struct {
	keyword  string
	filename string
}

var rules = []rule{
	// Accessories
	{"anting", "anting-penyu.png"},
	{"cincin", "cincin-penyu.png"},
	{"gelang & kalung", "gelang-kalung.png"}, // ← before plain "gelang"
	{"gelang", "gelang-kerbau.png"},
	{"kalung", "kalung-timor.jpg"},
	{"headband", "headband-tenun.png"},
	{"mahkota", "mahkota-tiilangga.png"},
	{"sisir", "sisir-sumba.png"},

	// Food & Beverage
	{"kopi", "kopi-flores.png"},
	{"madu", "madu-hutan.png"},
	{"sei", "sei-babi.png"},
	{"garam", "garam-gunung.png"},
	{"jagung", "jagung-titi.png"},
	{"kacang", "kacang-sembunyi.png"},
	{"sambal", "sambal-luat.png"},
	{"gula", "gula-air.png"},

	// Fashion / Tenun — specific cuts first
	{"kemeja", "tenun-kemeja.png"},
	{"kameja", "tenun-kemeja.png"}, // common typo / variant
	{"dress", "tenun-dress.png"},
	{"blouse", "tenun-blouse.png"},
	{"outer", "tenun-outer.png"},
	{"tunik", "product-1.png"},
	{"selendang", "product-2.png"},
	{"kaftan", "product-3.png"},
	{"syal", "product-4.png"},
}

// resolve returns the /images/… path for the given product name.
func resolve(name string) (string, bool) {
	n := strings.ToLower(name)
	for _, r := range rules {
		if strings.Contains(n, r.keyword) {
			return "/images/" + r.filename, true
		}
	}
	return "/images/product-1.png", false // ultimate fallback
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
	fmt.Println("🔧 Applying exact image mapping (ordered rules)...")

	var products []Product
	if err := db.Find(&products).Error; err != nil {
		log.Fatalf("❌ Failed to query products: %v\n", err)
	}
	fmt.Printf("   Found %d products\n\n", len(products))

	updated, fallback, failed := 0, 0, 0

	for i := range products {
		p := &products[i]

		imagePath, matched := resolve(p.Name)
		tag := "UPDATED "
		if !matched {
			tag = "FALLBACK"
			fallback++
		}

		if err := db.Model(p).UpdateColumn("image_url", imagePath).Error; err != nil {
			log.Printf("   [ERROR   ] #%d %s → %v\n", p.ID, p.Name, err)
			failed++
			continue
		}

		fmt.Printf("   [%s] #%-3d %-45s → %s\n",
			tag, p.ID, truncate(p.Name, 45), imagePath)
		if matched {
			updated++
		}
	}

	fmt.Println("\n────────────────────────────────────────────────────────────")
	fmt.Printf("✅ Done.  Exact match: %d | Fallback: %d | Errors: %d | Total: %d\n",
		updated, fallback, failed, len(products))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
