package websocket

type MessagePayload struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversation_id"`
	MessageID      uint   `json:"message_id"`
	SenderID       string `json:"sender_id"`
	Content        string `json:"content"`
	Status         string `json:"status"`
}

type ReadReceiptPayload struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversation_id"`
	MessageID      uint   `json:"message_id"`
	UserID         string `json:"user_id"`
}
type DeliveredPayload struct {
	Type           string `json:"type"`
	MessageID      uint   `json:"message_id"`
	ConversationID uint   `json:"conversation_id"`
	UserID         string `json:"user_id"`
}
