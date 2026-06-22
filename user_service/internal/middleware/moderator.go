package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ModeratorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		role := c.GetString(ContextRole)

		if role != "moderator" &&
			role != "admin" {

			c.JSON(http.StatusForbidden, gin.H{
				"error": "moderator access required",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
