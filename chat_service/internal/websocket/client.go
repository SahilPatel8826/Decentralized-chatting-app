package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string

	Conn *websocket.Conn

	Send chan []byte
}
type TypingPayload struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversation_id"`
	UserID         string `json:"user_id"`
}

func (c *Client) ReadPump(
	hub *Hub,
) {

	defer func() {

		hub.unregister <- c

		hub.BroadcastPresence(
			c.UserID,
			"offline",
		)

		c.Conn.Close()
	}()

	for {

		_, data, err :=
			c.Conn.ReadMessage()

		if err != nil {
			break
		}

		var payload TypingPayload

		err = json.Unmarshal(
			data,
			&payload,
		)

		if err != nil {
			continue
		}

		switch payload.Type {

		case "typing_start":

			hub.BroadcastTyping(
				payload.ConversationID,
				c.UserID,
				true,
			)

		case "typing_stop":

			hub.BroadcastTyping(
				payload.ConversationID,
				c.UserID,
				false,
			)
		}
	}

}

func (c *Client) WritePump() {

	defer c.Conn.Close()

	for {

		message, ok := <-c.Send

		if !ok {

			c.Conn.WriteMessage(
				websocket.CloseMessage,
				[]byte{},
			)

			return
		}

		err := c.Conn.WriteMessage(
			websocket.TextMessage,
			message,
		)

		if err != nil {

			log.Println(err)

			return
		}
	}
}
