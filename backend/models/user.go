package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents users in the OpenPeo marketplace.
// Supports three roles: admin, vendor, and customer.
type User struct {
	ID        int32          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"size:100" json:"name"`
	Email     string         `gorm:"size:100;uniqueIndex" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Role      string         `gorm:"size:20;not null;default:'customer'" json:"role"`

	Phone     string         `gorm:"size:20" json:"phone"`
	Address   string         `gorm:"type:text" json:"address"`
	Avatar        string         `gorm:"size:255" json:"avatar"`
	SellerProfile *SellerProfile `gorm:"foreignKey:UserID" json:"seller_profile,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default table name.
func (User) TableName() string {
	return "users"
}
