package service

import (
	"chat_service/internal/dto"
	"chat_service/internal/model"
	"chat_service/internal/repository"
	"context"
	"errors"
	"strings"
)

type MessageService struct {
	messageRepo      *repository.MessageRepository
	conversationRepo *repository.ConversationRepository
}

func NewMessageService(
	messageRepo *repository.MessageRepository,
	conversationRepo *repository.ConversationRepository,
) *MessageService {

	return &MessageService{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
	}
}

var (
	ErrMessageIDRequired      = errors.New("message id required")
	ErrMessageContentRequired = errors.New("message content required")

	ErrMessageNotFound = errors.New("message not found")
)

func (s *MessageService) SendMessage(
	ctx context.Context,
	userID string,
	req dto.SendMessageRequest,
) (*model.Message, error) {

	// Validation

	if userID == "" {
		return nil, ErrUserIDRequired
	}

	if req.ConversationID == 0 {
		return nil, ErrConversationIDRequired
	}

	if strings.TrimSpace(req.Content) == "" {
		return nil, ErrMessageContentRequired
	}

	// Conversation Exists?

	conversation, err :=
		s.conversationRepo.GetConversationByID(
			ctx,
			req.ConversationID,
		)

	if err != nil {
		return nil, err
	}

	if conversation == nil {
		return nil, ErrConversationNotFound
	}

	// User Is Participant?

	allowed, err :=
		s.conversationRepo.IsParticipant(
			ctx,
			req.ConversationID,
			userID,
		)

	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, ErrForbidden
	}

	// Create Message

	message := &model.Message{
		ConversationID: req.ConversationID,
		SenderID:       userID,
		Content:        strings.TrimSpace(req.Content),
		Status:         "sent",
	}

	err = s.messageRepo.CreateMessage(
		ctx,
		message,
	)

	if err != nil {
		return nil, err
	}

	// Update Conversation Last Message

	err = s.conversationRepo.UpdateLastMessageID(
		ctx,
		req.ConversationID,
		message.ID,
	)

	if err != nil {
		return nil, err
	}

	return message, nil
}
func (s *MessageService) GetMessages(
	ctx context.Context,
	userID string,
	conversationID uint,
	page int,
	limit int,
) ([]model.Message, error) {

	// Validation

	if userID == "" {
		return nil, ErrUserIDRequired
	}

	if conversationID == 0 {
		return nil, ErrConversationIDRequired
	}

	// Default pagination

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Conversation Exists?

	conversation, err :=
		s.conversationRepo.GetConversationByID(
			ctx,
			conversationID,
		)

	if err != nil {
		return nil, err
	}

	if conversation == nil {
		return nil, ErrConversationNotFound
	}

	// User Is Participant?

	allowed, err :=
		s.conversationRepo.IsParticipant(
			ctx,
			conversationID,
			userID,
		)

	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, ErrForbidden
	}

	// Fetch Messages

	messages, err :=
		s.messageRepo.GetMessagesByConversation(
			ctx,
			conversationID,
			limit,
			offset,
		)

	if err != nil {
		return nil, err
	}

	return messages, nil
}
func (s *MessageService) MarkRead(
	ctx context.Context,
	userID string,
	messageID uint,
) error {

	// Validation

	if userID == "" {
		return ErrUserIDRequired
	}

	if messageID == 0 {
		return ErrMessageIDRequired
	}

	// Message Exists?

	message, err := s.messageRepo.GetMessageByID(
		ctx,
		messageID,
	)

	if err != nil {
		return err
	}

	if message == nil {
		return ErrMessageNotFound
	}

	// User belongs to conversation?

	allowed, err := s.conversationRepo.IsParticipant(
		ctx,
		message.ConversationID,
		userID,
	)

	if err != nil {
		return err
	}

	if !allowed {
		return ErrForbidden
	}

	// Don't allow sender to mark own message as read

	if message.SenderID == userID {
		return errors.New(
			"cannot mark your own message as read",
		)
	}

	// Already read?

	if message.Status == "read" {
		return nil
	}

	// Update status

	return s.messageRepo.MarkRead(
		ctx,
		messageID,
	)
}
func (s *MessageService) DeleteMessage(
	ctx context.Context,
	userID string,
	messageID uint,
) error {

	if userID == "" {
		return ErrUserIDRequired
	}

	if messageID == 0 {
		return ErrMessageIDRequired
	}

	// Message Exists?

	message, err :=
		s.messageRepo.GetMessageByID(
			ctx,
			messageID,
		)

	if err != nil {
		return err
	}

	if message == nil {
		return ErrMessageNotFound
	}

	// Authorization

	if message.SenderID != userID {
		return ErrForbidden
	}

	// Delete

	err = s.messageRepo.DeleteMessage(
		ctx,
		messageID,
	)

	if err != nil {
		return err
	}

	return nil
}
