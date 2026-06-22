package routes

import (
	"chat_service/internal/handler"
	"chat_service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ConversationRoutes(
	router *gin.Engine,
	conversationHandler *handler.ConversationHandler,
) {

	conversations := router.Group("/conversations")

	conversations.Use(
		middleware.AuthMiddleware(),
	)

	conversations.POST(
		"/private",
		conversationHandler.CreatePrivateConversation,
	)

	conversations.POST(
		"/group",
		conversationHandler.CreateGroupConversation,
	)

	conversations.GET(
		"",
		conversationHandler.GetUserConversations,
	)

	conversations.GET(
		"/:conversationID",
		conversationHandler.GetConversationByID,
	)

	conversations.POST(
		"/:conversationID/participants",
		conversationHandler.AddParticipant,
	)

	conversations.DELETE(
		"/:conversationID/participants/:userID",
		conversationHandler.RemoveParticipant,
	)

	conversations.DELETE(
		"/:conversationID/leave",
		conversationHandler.LeaveGroup,
	)
}
