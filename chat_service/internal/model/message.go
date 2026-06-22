package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ConversationID uint   `json:"conversation_id"`
	SenderID       string `json:"sender_id"`
	Content        string `json:"content"`
	Status         string `json:"status"`
}
