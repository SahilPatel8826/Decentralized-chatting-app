package routes

import (
	"chat_app/internal/handler"
	"chat_app/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
) {

	users := router.Group("/users")

	users.Use(
		middleware.AuthMiddleware(),
	)

	users.GET(
		"/me",
		userHandler.GetMe,
	)

	users.GET(
		"/:id",
		userHandler.GetUserByID,
	)

	users.PATCH(
		"/me",
		userHandler.UpdateProfile,
	)

	users.DELETE(
		"/me",
		userHandler.DeleteAccount,
	)

	users.GET(
		"/search",
		userHandler.SearchUsers,
	)
}
