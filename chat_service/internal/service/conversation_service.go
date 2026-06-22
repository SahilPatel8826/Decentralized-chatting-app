package service

import (
	"chat_service/internal/dto"
	"chat_service/internal/model"
	"chat_service/internal/repository"
	"context"
	"errors"
)

type ConversationService struct {
	repo *repository.ConversationRepository
}

var (
	ErrUserIDRequired         = errors.New("user id is required")
	ErrParticipantIDRequired  = errors.New("participant id is required")
	ErrConversationIDRequired = errors.New("conversation id is required")
	ErrCannotChatWithSelf     = errors.New("cannot create conversation with self")
	ErrUserNotFound           = errors.New("user not found")
	ErrConversationNotFound   = errors.New("conversation not found")
	ErrForbidden              = errors.New("forbidden")
	ErrGroupNameRequired      = errors.New("group name is required")
	ErrParticipantsRequired   = errors.New("participants are required")
)

func NewConversationService(repo *repository.ConversationRepository) *ConversationService {
	return &ConversationService{
		repo: repo,
	}
}

func (s *ConversationService) CreatePrivateConversation(
	ctx context.Context,
	userID string,
	participantID string,
) (*model.Conversation, error) {

	if userID == "" {
		return nil, ErrUserIDRequired
	}

	if participantID == "" {
		return nil, ErrParticipantIDRequired
	}

	if userID == participantID {
		return nil, ErrCannotChatWithSelf
	}

	// exists, err := s.userClient.UserExists(
	// 	ctx,
	// 	participantID,
	// )

	// if err != nil {
	// 	return nil, err
	// }

	// if !exists {
	// 	return nil, ErrUserNotFound
	// }

	existingConversation, err :=
		s.repo.FindPrivateConversation(
			ctx,
			userID,
			participantID,
		)

	if err != nil {
		return nil, err
	}

	if existingConversation != nil {
		return existingConversation, nil
	}

	conversation := &model.Conversation{
		Type: "private",
	}

	err = s.repo.CreateConversation(
		ctx,
		conversation,
	)

	if err != nil {
		return nil, err
	}

	participants := []model.ConversationParticipant{
		{
			ConversationID: conversation.ID,
			UserID:         userID,
			Role:           "member",
		},
		{
			ConversationID: conversation.ID,
			UserID:         participantID,
			Role:           "member",
		},
	}

	err = s.repo.CreateParticipants(
		ctx,
		participants,
	)

	if err != nil {
		return nil, err
	}

	return conversation, nil
}
func (s *ConversationService) GetConversationByID(
	ctx context.Context,
	userID string,
	conversationID uint,
) (*model.Conversation, error) {

	if userID == "" {
		return nil, ErrUserIDRequired
	}

	if conversationID == 0 {
		return nil, ErrConversationNotFound
	}

	conversation, err :=
		s.repo.GetConversationByID(
			ctx,
			conversationID,
		)

	if err != nil {
		return nil, err
	}

	if conversation == nil {
		return nil, ErrConversationNotFound
	}

	allowed, err :=
		s.repo.IsParticipant(
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

	return conversation, nil
}
func (s *ConversationService) GetUserConversations(
	ctx context.Context,
	userID string,
) ([]model.Conversation, error) {

	if userID == "" {
		return nil, ErrUserIDRequired
	}

	conversations, err :=
		s.repo.GetUserConversations(
			ctx,
			userID,
		)

	if err != nil {
		return nil, err
	}

	if len(conversations) == 0 {
		return []model.Conversation{}, nil
	}

	return conversations, nil
}
func (s *ConversationService) CreateGroupConversation(
	ctx context.Context,
	creatorID string,
	req dto.CreateGroupConversationRequest,
) (*model.Conversation, error) {

	if creatorID == "" {
		return nil, ErrUserIDRequired
	}

	if req.Name == "" {
		return nil, ErrGroupNameRequired
	}

	if len(req.Participants) == 0 {
		return nil, ErrParticipantsRequired
	}

	conversation := &model.Conversation{
		Type: "group",
		Name: req.Name,
	}

	err := s.repo.CreateConversation(
		ctx,
		conversation,
	)

	if err != nil {
		return nil, err
	}

	participants := []model.ConversationParticipant{
		{
			ConversationID: conversation.ID,
			UserID:         creatorID,
			Role:           "admin",
		},
	}

	for _, participantID := range req.Participants {

		if participantID == creatorID {
			continue
		}

		participants = append(
			participants,
			model.ConversationParticipant{
				ConversationID: conversation.ID,
				UserID:         participantID,
				Role:           "member",
			},
		)
	}

	err = s.repo.CreateParticipants(
		ctx,
		participants,
	)

	if err != nil {
		return nil, err
	}

	return conversation, nil
}
func (s *ConversationService) AddParticipant(
	ctx context.Context,
	requesterID string,
	conversationID uint,
	req dto.AddParticipantRequest,
) error {

	if requesterID == "" {
		return ErrUserIDRequired
	}

	if req.UserID == "" {
		return ErrUserIDRequired
	}

	conversation, err :=
		s.repo.GetConversationByID(
			ctx,
			conversationID,
		)

	if err != nil {
		return err
	}

	if conversation == nil {
		return ErrConversationNotFound
	}

	if conversation.Type != "group" {
		return errors.New(
			"cannot add participants to private conversation",
		)
	}

	// requester belongs to group?

	allowed, err :=
		s.repo.IsParticipant(
			ctx,
			conversationID,
			requesterID,
		)

	if err != nil {
		return err
	}

	if !allowed {
		return ErrForbidden
	}

	// requester admin?

	participants, err :=
		s.repo.GetParticipants(
			ctx,
			conversationID,
		)

	if err != nil {
		return err
	}
	isAdmin := false

	for _, participant := range participants {

		if participant.UserID == requesterID &&
			participant.Role == "admin" {

			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return ErrForbidden
	}
	// already exists?

	exists, err :=
		s.repo.IsParticipant(
			ctx,
			conversationID,
			req.UserID,
		)

	if err != nil {
		return err
	}

	if exists {
		return errors.New(
			"user already participant",
		)
	}

	err = s.repo.AddParticipant(
		ctx,
		&model.ConversationParticipant{
			ConversationID: conversationID,
			UserID:         req.UserID,
			Role:           "member",
		},
	)

	if err != nil {
		return err
	}

	return nil
}
func (s *ConversationService) RemoveParticipant(
	ctx context.Context,
	requesterID string,
	conversationID uint,
	targetUserID string,
) error {

	// Validation

	if requesterID == "" {
		return ErrUserIDRequired
	}

	if targetUserID == "" {
		return ErrUserIDRequired
	}

	if conversationID == 0 {
		return ErrConversationIDRequired
	}

	// Conversation Exists?

	conversation, err := s.repo.GetConversationByID(
		ctx,
		conversationID,
	)

	if err != nil {
		return err
	}

	if conversation == nil {
		return ErrConversationNotFound
	}

	// Only group chats

	if conversation.Type != "group" {
		return errors.New(
			"participants can only be removed from groups",
		)
	}

	// Requester must belong to group

	allowed, err := s.repo.IsParticipant(
		ctx,
		conversationID,
		requesterID,
	)

	if err != nil {
		return err
	}

	if !allowed {
		return ErrForbidden
	}

	// Load participants

	participants, err := s.repo.GetParticipants(
		ctx,
		conversationID,
	)

	if err != nil {
		return err
	}

	// Check requester is admin

	isAdmin := false

	for _, participant := range participants {

		if participant.UserID == requesterID &&
			participant.Role == "admin" {

			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return ErrForbidden
	}

	// Check target exists

	targetFound := false

	for _, participant := range participants {

		if participant.UserID == targetUserID {

			targetFound = true
			break
		}
	}

	if !targetFound {
		return errors.New(
			"user is not a participant",
		)
	}

	// Optional safety

	if requesterID == targetUserID {
		return errors.New(
			"use leave group instead",
		)
	}

	// Remove participant

	err = s.repo.RemoveParticipant(
		ctx,
		conversationID,
		targetUserID,
	)

	if err != nil {
		return err
	}

	return nil
}
func (s *ConversationService) LeaveGroup(
	ctx context.Context,
	userID string,
	conversationID uint,
) error {

	if userID == "" {
		return ErrUserIDRequired
	}

	if conversationID == 0 {
		return ErrConversationIDRequired
	}

	// Conversation exists?

	conversation, err := s.repo.GetConversationByID(
		ctx,
		conversationID,
	)

	if err != nil {
		return err
	}

	if conversation == nil {
		return ErrConversationNotFound
	}

	// Only groups can be left

	if conversation.Type != "group" {
		return errors.New(
			"cannot leave private conversation",
		)
	}

	// User belongs to group?

	allowed, err := s.repo.IsParticipant(
		ctx,
		conversationID,
		userID,
	)

	if err != nil {
		return err
	}

	if !allowed {
		return ErrForbidden
	}

	// Remove user

	err = s.repo.RemoveParticipant(
		ctx,
		conversationID,
		userID,
	)

	if err != nil {
		return err
	}

	// Check remaining members

	participants, err := s.repo.GetParticipants(
		ctx,
		conversationID,
	)

	if err != nil {
		return err
	}

	// No participants left

	if len(participants) == 0 {

		err = s.repo.DeleteConversation(
			ctx,
			conversationID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
