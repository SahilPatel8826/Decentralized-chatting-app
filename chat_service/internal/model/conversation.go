package model

import "time"

type Conversation struct {
	ID            uint
	Type          string // private, group
	Name          string // group name
	LastMessageID *uint

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ConversationParticipant struct {
	ID             uint
	ConversationID uint
	UserID         string

	Role string // member, admin

	JoinedAt time.Time
}
