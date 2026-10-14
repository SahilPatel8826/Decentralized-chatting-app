package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	Hub *Hub
}

func NewHandler(
	hub *Hub,
) *Handler {

	return &Handler{
		Hub: hub,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(
		r *http.Request,
	) bool {

		return true
	},
}

func (h *Handler) Connect(
	c *gin.Context,
) {

	userID :=
		c.GetString(
			"userID",
		)

	conn, err :=
		upgrader.Upgrade(
			c.Writer,
			c.Request,
			nil,
		)

	if err != nil {
		return
	}

	client := &Client{
		UserID: userID,

		Conn: conn,

		Send: make(
			chan []byte,
			256,
		),
	}

	h.Hub.Register(client)

	go h.Hub.BroadcastPresence(
		client.UserID,
		"online",
	)

	go client.WritePump()

	go client.ReadPump(
		h.Hub,
	)
}
