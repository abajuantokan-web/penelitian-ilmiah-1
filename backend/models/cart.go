package models

import (
	"time"
)

// CartItem represents an item in a user's shopping cart.
type CartItem struct {
	ID        int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int32     `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	ProductID int32     `gorm:"not null;index" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the default table name.
func (CartItem) TableName() string {
	return "cart_items"
}
