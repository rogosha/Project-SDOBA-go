package model

import "time"

type Message struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ConversationID uint      `json:"conversationId"`
	SenderID       uint      `json:"senderId"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	Sender User `gorm:"foreignKey:SenderID" json:"sender"`
}
