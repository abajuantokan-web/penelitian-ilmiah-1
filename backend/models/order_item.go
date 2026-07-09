package models

import (
	"time"
)

// OrderItem represents a single product line within an order.
// This allows an order to contain multiple different products.
type OrderItem struct {
	ID        int32   `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   int32   `gorm:"not null;index" json:"order_id"`
	ProductID int32   `gorm:"not null;index" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"not null;default:1" json:"quantity"`
	Price     float64 `gorm:"type:decimal(12,2);not null" json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName overrides the default table name.
func (OrderItem) TableName() string {
	return "order_items"
}
