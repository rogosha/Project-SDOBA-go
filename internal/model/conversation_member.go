package model

type ConversationMember struct {
	ID             uint `gorm:"primaryKey" json:"id"`
	ConversationID uint `json:"conversationId"`
	UserID         uint `json:"userId"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}
