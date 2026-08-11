package models

import "time"



type Message struct {
	ID         int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID   int32     `gorm:"not null;index" json:"sender_id"`
	ReceiverID int32     `gorm:"not null;index" json:"receiver_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	IsRead     bool      `gorm:"not null;default:false" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}


func (Message) TableName() string {
	return "messages"
}


type ChatPayload struct {
	SenderID   int32  `json:"sender_id"`
	ReceiverID int32  `json:"receiver_id"`
	Content    string `json:"content"`
}
