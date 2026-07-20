package routes

import (
	"chat_service/internal/middleware"
	ws "chat_service/internal/websocket"

	"github.com/gin-gonic/gin"
)

func WebSocketRoutes(
	router *gin.Engine,
	wsHandler *ws.Handler,
) {

	router.GET(
		"/ws",
		middleware.AuthMiddleware(),
		wsHandler.Connect,
	)
}
