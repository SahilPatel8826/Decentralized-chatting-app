package routes

import (
	"github.com/gin-gonic/gin"

	"chat_app/internal/handler"
)

func AuthRoutes(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
) {

	auth := router.Group("/auth")

	auth.POST(
		"/register",
		authHandler.Register,
	)

	auth.POST(
		"/login",
		authHandler.Login,
	)

	auth.POST(
		"/refresh",
		authHandler.RefreshToken,
	)

	auth.POST(
		"/logout",
		authHandler.Logout,
	)
}
