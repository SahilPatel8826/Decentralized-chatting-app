package websocket

import "encoding/json"

type Hub struct {
	clients map[string]*Client

	register chan *Client

	unregister chan *Client

	broadcast chan []byte
}

func NewHub() *Hub {

	return &Hub{
		clients: make(map[string]*Client),

		register: make(chan *Client),

		unregister: make(chan *Client),

		broadcast: make(chan []byte),
	}
}

func (h *Hub) Run() {

	for {

		select {

		case client := <-h.register:

			h.clients[client.UserID] = client

		case client := <-h.unregister:

			if _, ok := h.clients[client.UserID]; ok {

				delete(
					h.clients,
					client.UserID,
				)

				close(client.Send)
			}

		case message := <-h.broadcast:

			for _, client := range h.clients {

				select {

				case client.Send <- message:

				default:

					close(client.Send)

					delete(
						h.clients,
						client.UserID,
					)
				}
			}
		}
	}
}

func (h *Hub) SendToUser(
	userID string,
	message []byte,
) {

	client, ok := h.clients[userID]

	if !ok {
		return
	}

	select {

	case client.Send <- message:

	default:
		close(client.Send)
		delete(h.clients, userID)
	}
}
func (h *Hub) IsOnline(
	userID string,
) bool {

	_, ok := h.clients[userID]

	return ok
}
func (h *Hub) OnlineUsers() []string {

	users := make(
		[]string,
		0,
	)

	for userID := range h.clients {

		users = append(
			users,
			userID,
		)
	}

	return users
}
func (h *Hub) Broadcast(
	message []byte,
) {

	for _, client := range h.clients {

		select {

		case client.Send <- message:

		default:

			close(client.Send)

			delete(
				h.clients,
				client.UserID,
			)
		}
	}

}

type PresencePayload struct {
	Type   string `json:"type"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

func (h *Hub) BroadcastPresence(
	userID string,
	status string,
) {

	payload := PresencePayload{
		Type:   "presence",
		UserID: userID,
		Status: status,
	}

	data, err := json.Marshal(
		payload,
	)

	if err != nil {
		return
	}

	h.Broadcast(data)

}
func (h *Hub) BroadcastTyping(
	conversationID uint,
	userID string,
	isTyping bool,
) {

	eventType := "typing_stop"

	if isTyping {
		eventType = "typing_start"
	}

	payload := TypingPayload{
		Type:           eventType,
		ConversationID: conversationID,
		UserID:         userID,
	}

	data, err := json.Marshal(
		payload,
	)

	if err != nil {
		return
	}

	h.Broadcast(data)

}
