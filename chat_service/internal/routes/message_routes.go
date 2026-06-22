package routes

import (
	"chat_service/internal/handler"
	"chat_service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func MessageRoutes(
	router *gin.Engine,
	messageHandler *handler.MessageHandler,
) {

	messages := router.Group("/messages")

	messages.Use(
		middleware.AuthMiddleware(),
	)

	messages.POST(
		"",
		messageHandler.SendMessage,
	)

	messages.PATCH(
		"/:messageID/read",
		messageHandler.MarkRead,
	)

	messages.DELETE(
		"/:messageID",
		messageHandler.DeleteMessage,
	)

	router.GET(
		"/conversations/:conversationID/messages",
		middleware.AuthMiddleware(),
		messageHandler.GetMessages,
	)
}
