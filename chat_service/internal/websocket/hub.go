package websocket

import (
	"encoding/json"
	"sync"
)

// Hub tracks connected clients (one connection per user).
// All access to clients goes through mu.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client

	// Set from main. Hub stays decoupled from the database.
	// Participants returns the user IDs in a conversation.
	Participants func(conversationID uint) []string
	// Contacts returns the user IDs that share a conversation with userID.
	Contacts func(userID string) []string
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

// Register adds a client. If the user is already connected,
// the old connection is closed and replaced.
func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if old, ok := h.clients[c.UserID]; ok {
		delete(h.clients, c.UserID)
		close(old.Send)
	}

	h.clients[c.UserID] = c
}

// Unregister removes c only if it is still the active
// connection for that user (a newer connection stays).
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if cur, ok := h.clients[c.UserID]; ok && cur == c {
		delete(h.clients, c.UserID)
		close(c.Send)
	}
}

// SendToUser does a non-blocking send. A client whose buffer
// is full is dropped.
func (h *Hub) SendToUser(userID string, message []byte) {
	h.mu.RLock()
	c, ok := h.clients[userID]
	sent := false
	if ok {
		select {
		case c.Send <- message:
			sent = true
		default:
		}
	}
	h.mu.RUnlock()

	if ok && !sent {
		h.Unregister(c)
	}
}

func (h *Hub) IsOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.clients[userID]
	return ok
}

func (h *Hub) OnlineUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]string, 0, len(h.clients))
	for id := range h.clients {
		users = append(users, id)
	}
	return users
}

type PresencePayload struct {
	Type   string `json:"type"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

// BroadcastPresence notifies only the user's contacts.
func (h *Hub) BroadcastPresence(userID string, status string) {
	if h.Contacts == nil {
		return
	}

	data, err := json.Marshal(PresencePayload{
		Type:   "presence",
		UserID: userID,
		Status: status,
	})
	if err != nil {
		return
	}

	for _, id := range h.Contacts(userID) {
		h.SendToUser(id, data)
	}
}

// BroadcastTyping notifies the other participants of the
// conversation. Ignored if the sender is not a participant.
func (h *Hub) BroadcastTyping(
	conversationID uint,
	userID string,
	isTyping bool,
) {
	if h.Participants == nil {
		return
	}

	participants := h.Participants(conversationID)

	member := false
	for _, id := range participants {
		if id == userID {
			member = true
			break
		}
	}
	if !member {
		return
	}

	eventType := "typing_stop"
	if isTyping {
		eventType = "typing_start"
	}

	data, err := json.Marshal(TypingPayload{
		Type:           eventType,
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		return
	}

	for _, id := range participants {
		if id != userID {
			h.SendToUser(id, data)
		}
	}
}
