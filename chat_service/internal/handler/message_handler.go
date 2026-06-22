package handler

import (
	"chat_service/internal/dto"
	"chat_service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(
	messageService *service.MessageService,
) *MessageHandler {

	return &MessageHandler{
		messageService: messageService,
	}
}
func (h *MessageHandler) SendMessage(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	var req dto.SendMessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	message, err :=
		h.messageService.SendMessage(
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
		http.StatusCreated,
		gin.H{
			"message": message,
		},
	)
}
func (h *MessageHandler) GetMessages(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	conversationIDStr :=
		c.Param("conversationID")

	conversationID, err :=
		strconv.ParseUint(
			conversationIDStr,
			10,
			64,
		)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid conversation id",
			},
		)

		return
	}

	page, _ :=
		strconv.Atoi(
			c.DefaultQuery("page", "1"),
		)

	limit, _ :=
		strconv.Atoi(
			c.DefaultQuery("limit", "20"),
		)

	messages, err :=
		h.messageService.GetMessages(
			c.Request.Context(),
			userID,
			uint(conversationID),
			page,
			limit,
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
			"messages": messages,
		},
	)
}
func (h *MessageHandler) MarkRead(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	messageIDStr :=
		c.Param("messageID")

	messageID, err :=
		strconv.ParseUint(
			messageIDStr,
			10,
			64,
		)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid message id",
			},
		)

		return
	}

	err = h.messageService.MarkRead(
		c.Request.Context(),
		userID,
		uint(messageID),
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
			"message": "message marked as read",
		},
	)
}
func (h *MessageHandler) DeleteMessage(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	messageIDStr :=
		c.Param("messageID")

	messageID, err :=
		strconv.ParseUint(
			messageIDStr,
			10,
			64,
		)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid message id",
			},
		)

		return
	}

	err = h.messageService.DeleteMessage(
		c.Request.Context(),
		userID,
		uint(messageID),
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
			"message": "message deleted successfully",
		},
	)
}
