package models

import (
	"time"

	"gorm.io/gorm"
)

// Order represents a customer's purchase/pre-order transaction.
// Each order is linked to a specific product and the purchasing customer.
type Order struct {
	ID         int32          `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID int32          `gorm:"not null;index" json:"customer_id"`
	Customer   User           `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	ProductID  int32          `gorm:"not null;index" json:"product_id"`
	Product    Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity   int            `gorm:"not null" json:"quantity"`
	TotalPrice float64        `gorm:"type:decimal(12,2);not null" json:"total_price"`
	Status     string         `gorm:"size:50;not null;default:'pending'" json:"status"`
	Note       string         `gorm:"type:text" json:"note"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default table name.
func (Order) TableName() string {
	return "orders"
}

// OrderRequest is the payload shape for creating a new order.
type OrderRequest struct {
	CustomerID int32  `json:"customer_id" binding:"required"`
	ProductID  int32  `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	Note       string `json:"note"`
}
