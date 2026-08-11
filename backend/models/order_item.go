package models

import (
	"time"
)

type OrderItem struct {
	ID		int32		`gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID		int32		`gorm:"not null;index" json:"order_id"`
	ProductID	int32		`gorm:"not null;index" json:"product_id"`
	Product		Product		`gorm:"foreignKey:ProductID" json:"product"`
	Quantity	int		`gorm:"not null;default:1" json:"quantity"`
	Price		float64		`gorm:"type:decimal(12,2);not null" json:"price"`
	CreatedAt	time.Time	`json:"created_at"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
