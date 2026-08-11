package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID			int32		`gorm:"primaryKey;autoIncrement" json:"id"`
	SellerID		int32		`gorm:"not null;index" json:"seller_id"`
	SellerProfile		SellerProfile	`gorm:"foreignKey:SellerID" json:"seller_profile,omitempty"`
	Name			string		`gorm:"size:200;not null" json:"name"`
	Description		string		`gorm:"type:text" json:"description"`
	Price			float64		`gorm:"type:decimal(12,2);not null" json:"price"`
	ImageURL		string		`gorm:"size:255" json:"image_url"`
	Region			string		`gorm:"size:100;not null;index" json:"region"`
	Category		string		`gorm:"size:100;index" json:"category"`
	MinOrder		int		`gorm:"not null;default:1" json:"min_order"`
	PreOrderDuration	int		`gorm:"not null;default:7" json:"pre_order_duration"`
	Stock			int		`gorm:"not null;default:0" json:"stock"`
	IsActive		bool		`gorm:"not null;default:true" json:"is_active"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	DeletedAt		gorm.DeletedAt	`gorm:"index" json:"-"`
}

func (Product) TableName() string {
	return "products"
}
