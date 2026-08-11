package models

import "time"

type ActivityLog struct {
	ID		int32		`gorm:"primaryKey" json:"id"`
	UserID		int32		`json:"user_id"`
	UserRole	string		`json:"user_role"`
	ActionType	string		`json:"action_type"`
	Description	string		`json:"description"`
	CreatedAt	time.Time	`json:"created_at"`
}
