package handler

import (
	"net/http"

	"chat_app/internal/dto"
	"chat_app/internal/middleware"
	"chat_app/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(
	userService *service.UserService,
) *UserHandler {

	return &UserHandler{
		userService: userService,
	}
}
func (h *UserHandler) GetMe(
	c *gin.Context,
) {

	userID := c.GetString(
		middleware.ContextUserID,
	)

	user, err :=
		h.userService.GetByID(
			c.Request.Context(),
			userID,
		)

	if err != nil {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "user not found",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		user,
	)
}
func (h *UserHandler) GetUserByID(
	c *gin.Context,
) {

	id := c.Param("id")

	user, err :=
		h.userService.GetByID(
			c.Request.Context(),
			id,
		)

	if err != nil {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "user not found",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		user,
	)
}
func (h *UserHandler) UpdateProfile(
	c *gin.Context,
) {

	userID := c.GetString(
		middleware.ContextUserID,
	)

	var req dto.UpdateProfileRequest

	if err :=
		c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err := h.userService.UpdateProfile(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "profile updated",
		},
	)
}
func (h *UserHandler) DeleteAccount(
	c *gin.Context,
) {

	userID := c.GetString(
		middleware.ContextUserID,
	)

	err := h.userService.Delete(
		c.Request.Context(),
		userID,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "account deleted",
		},
	)
}
func (h *UserHandler) SearchUsers(
	c *gin.Context,
) {

	query := c.Query("q")

	users, err :=
		h.userService.Search(
			c.Request.Context(),
			query,
		)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		users,
	)
}
