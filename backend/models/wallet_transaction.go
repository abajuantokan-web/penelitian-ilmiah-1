package models

import (
	"time"

	"gorm.io/gorm"
)

type WalletTransaction struct {
	ID		int32		`gorm:"primaryKey;autoIncrement" json:"id"`
	SellerID	int32		`gorm:"not null" json:"seller_id"`
	Type		string		`gorm:"size:50;not null" json:"type"`
	Amount		float64		`gorm:"not null" json:"amount"`
	Status		string		`gorm:"size:50;not null" json:"status"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	DeletedAt	gorm.DeletedAt	`gorm:"index" json:"-"`
}

func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}
