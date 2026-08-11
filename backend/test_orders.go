//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"openpeo-backend/config"
	"openpeo-backend/models"
)

func main() {
	config.ConnectDatabase()

	var orders []models.Order
	config.DB.
		Where("id = ?", 31).
		Preload("Product").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Find(&orders)

	b, _ := json.MarshalIndent(orders, "", "  ")
	fmt.Println(string(b))
}

