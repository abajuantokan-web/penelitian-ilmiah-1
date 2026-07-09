package models

import "time"

// Message represents a chat message between two users.
// Messages are persisted to MySQL for history and routed in real-time via WebSocket.
type Message struct {
	ID         int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID   int32     `gorm:"not null;index" json:"sender_id"`
	ReceiverID int32     `gorm:"not null;index" json:"receiver_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	IsRead     bool      `gorm:"not null;default:false" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName overrides the default table name.
func (Message) TableName() string {
	return "messages"
}

// ChatPayload is the JSON structure exchanged over the WebSocket connection.
type ChatPayload struct {
	SenderID   int32  `json:"sender_id"`
	ReceiverID int32  `json:"receiver_id"`
	Content    string `json:"content"`
}
