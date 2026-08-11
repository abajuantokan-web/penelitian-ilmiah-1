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

type rule struct {
	keyword		string
	filename	string
}

var rules = []rule{

	{"anting", "anting-penyu.png"},
	{"cincin", "cincin-penyu.png"},
	{"gelang & kalung", "gelang-kalung.png"},
	{"gelang", "gelang-kerbau.png"},
	{"kalung", "kalung-timor.jpg"},
	{"headband", "headband-tenun.png"},
	{"mahkota", "mahkota-tiilangga.png"},
	{"sisir", "sisir-sumba.png"},

	{"kopi", "kopi-flores.png"},
	{"madu", "madu-hutan.png"},
	{"sei", "sei-babi.png"},
	{"garam", "garam-gunung.png"},
	{"jagung", "jagung-titi.png"},
	{"kacang", "kacang-sembunyi.png"},
	{"sambal", "sambal-luat.png"},
	{"gula", "gula-air.png"},

	{"kemeja", "tenun-kemeja.png"},
	{"kameja", "tenun-kemeja.png"},
	{"dress", "tenun-dress.png"},
	{"blouse", "tenun-blouse.png"},
	{"outer", "tenun-outer.png"},
	{"tunik", "product-1.png"},
	{"selendang", "product-2.png"},
	{"kaftan", "product-3.png"},
	{"syal", "product-4.png"},
}

func resolve(name string) (string, bool) {
	n := strings.ToLower(name)
	for _, r := range rules {
		if strings.Contains(n, r.keyword) {
			return "/images/" + r.filename, true
		}
	}
	return "/images/product-1.png", false
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

