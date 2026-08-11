package models

import (
	"time"
)

type CartItem struct {
	ID		int32		`gorm:"primaryKey;autoIncrement" json:"id"`
	UserID		int32		`gorm:"not null;index" json:"user_id"`
	User		User		`gorm:"foreignKey:UserID" json:"-"`
	ProductID	int32		`gorm:"not null;index" json:"product_id"`
	Product		Product		`gorm:"foreignKey:ProductID" json:"product"`
	Quantity	int		`gorm:"not null;default:1" json:"quantity"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

func (CartItem) TableName() string {
	return "cart_items"
}
