package dto

import "time"

type CreatePrivateConversationRequest struct {
	ParticipantID string `json:"participant_id" binding:"required"`
}

type CreateGroupConversationRequest struct {
	Name         string   `json:"name" binding:"required"`
	Participants []string `json:"participants" binding:"required"`
}

type AddParticipantRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type RemoveParticipantRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type SendMessageRequest struct {
	ConversationID uint   `json:"conversation_id" binding:"required"`
	Content        string `json:"content" binding:"required"`
}

type GetMessagesRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type MarkReadRequest struct {
	MessageID uint `json:"message_id" binding:"required"`
}

type ConversationResponse struct {
	ID            uint      `json:"id"`
	Type          string    `json:"type"`
	Name          string    `json:"name"`
	LastMessageID *uint     `json:"last_message_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type MessageResponse struct {
	ID             uint      `json:"id"`
	ConversationID uint      `json:"conversation_id"`
	SenderID       string    `json:"sender_id"`
	Content        string    `json:"content"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
