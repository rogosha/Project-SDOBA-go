package model

import "time"

type Conversation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Members  []ConversationMember `gorm:"foreignKey:ConversationID" json:"members"`
	Messages []Message            `json:"messages,omitempty"`
}
