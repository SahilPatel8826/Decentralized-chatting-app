package dto

type MessagePayload struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversation_id"`
	SenderID       string `json:"sender_id"`
	Content        string `json:"content"`
}
