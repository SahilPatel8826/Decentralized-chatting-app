package handler

import (
	"net/http"
	"strconv"

	"chat_service/internal/dto"
	"chat_service/internal/service"

	"github.com/gin-gonic/gin"
)

type ConversationHandler struct {
	conversationService *service.ConversationService
}

func NewConversationHandler(
	conversationService *service.ConversationService,
) *ConversationHandler {

	return &ConversationHandler{
		conversationService: conversationService,
	}
}
func (h *ConversationHandler) CreatePrivateConversation(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	var req dto.CreatePrivateConversationRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	conversation, err :=
		h.conversationService.CreatePrivateConversation(
			c.Request.Context(),
			userID,
			req.ParticipantID,
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
			"conversation": conversation,
		},
	)
}
func (h *ConversationHandler) CreateGroupConversation(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	var req dto.CreateGroupConversationRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	conversation, err :=
		h.conversationService.CreateGroupConversation(
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
			"conversation": conversation,
		},
	)
}
func (h *ConversationHandler) GetConversationByID(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	idParam := c.Param("conversationID")

	id, err := strconv.ParseUint(
		idParam,
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

	conversation, err :=
		h.conversationService.GetConversationByID(
			c.Request.Context(),
			userID,
			uint(id),
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
			"conversation": conversation,
		},
	)
}
func (h *ConversationHandler) GetUserConversations(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	conversations, err :=
		h.conversationService.GetUserConversations(
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
			"conversations": conversations,
		},
	)
}
func (h *ConversationHandler) AddParticipant(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	conversationIDParam :=
		c.Param("conversationID")

	conversationID, err :=
		strconv.ParseUint(
			conversationIDParam,
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

	var req dto.AddParticipantRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err = h.conversationService.AddParticipant(
		c.Request.Context(),
		userID,
		uint(conversationID),
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
			"message": "participant added",
		},
	)
}
func (h *ConversationHandler) RemoveParticipant(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	conversationIDParam :=
		c.Param("conversationID")

	targetUserID :=
		c.Param("userID")

	conversationID, err :=
		strconv.ParseUint(
			conversationIDParam,
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

	err = h.conversationService.RemoveParticipant(
		c.Request.Context(),
		userID,
		uint(conversationID),
		targetUserID,
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
			"message": "participant removed",
		},
	)
}
func (h *ConversationHandler) LeaveGroup(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	conversationIDParam :=
		c.Param("conversationID")

	conversationID, err :=
		strconv.ParseUint(
			conversationIDParam,
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

	err = h.conversationService.LeaveGroup(
		c.Request.Context(),
		userID,
		uint(conversationID),
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
			"message": "left group successfully",
		},
	)
}
