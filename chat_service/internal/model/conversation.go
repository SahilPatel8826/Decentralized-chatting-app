package model

import "time"

type Conversation struct {
	ID            uint `gorm:"primaryKey"`
	Type          string
	Name          string
	LastMessageID *uint

	Participants []ConversationParticipant `gorm:"foreignKey:ConversationID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ConversationParticipant struct {
	ID             uint `gorm:"primaryKey"`
	ConversationID uint
	UserID         string
	Role           string

	JoinedAt time.Time
}
