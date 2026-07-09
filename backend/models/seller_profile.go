package models

import (
	"time"

	"gorm.io/gorm"
)

type SellerProfile struct {
	ID                int32          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            int32          `gorm:"uniqueIndex;not null" json:"user_id"`
	User              User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	StoreName         string         `gorm:"size:255" json:"store_name"`
	Description       string         `gorm:"type:text" json:"description"`
	StoreLogo         string         `gorm:"size:255" json:"store_logo"`
	Phone             string         `gorm:"size:50" json:"phone"`
	Address           string         `gorm:"type:text" json:"address"`
	Region            string         `gorm:"size:100" json:"region"`
	BankName          string         `gorm:"size:100" json:"bank_name"`
	BankAccountNumber string         `gorm:"size:100" json:"bank_account_number"`
	BankAccountName   string         `gorm:"size:255" json:"bank_account_name"`
	ActiveBalance     float64        `gorm:"default:0" json:"active_balance"`
	PendingBalance    float64        `gorm:"default:0" json:"pending_balance"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}
