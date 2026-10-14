package middleware

import (
	"net/http"
	"strings"

	"chat_service/pkg/jwtutil"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader :=
			c.GetHeader("Authorization")

		// Browsers cannot set headers on a WebSocket handshake,
		// so /ws also accepts ?token=<access token>.
		if authHeader == "" && c.Request.URL.Path == "/ws" {
			if t := c.Query("token"); t != "" {
				authHeader = "Bearer " + t
			}
		}

		if authHeader == "" {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "missing token",
				},
			)

			c.Abort()
			return
		}

		tokenString :=
			strings.TrimPrefix(
				authHeader,
				"Bearer ",
			)

		claims, err :=
			jwtutil.ValidateAccessToken(
				tokenString,
			)

		if err != nil {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid token",
				},
			)

			c.Abort()
			return
		}

		c.Set(
			"userID",
			claims.UserID,
		)

		c.Set(
			"email",
			claims.Email,
		)

		c.Set(
			"role",
			claims.Role,
		)

		c.Next()
	}
}
