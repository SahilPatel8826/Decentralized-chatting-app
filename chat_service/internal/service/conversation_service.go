package service

import (
	"chat_service/internal/client"
	"chat_service/internal/dto"
	"chat_service/internal/model"
	"chat_service/internal/repository"
	"chat_service/pb"
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConversationService struct {
	repo *repository.ConversationRepository

	userClient client.UserService
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
	ErrInvalidUserID          = errors.New("invalid user id")
	ErrUserServiceUnavailable = errors.New("user service unavailable")
)

func NewConversationService(
	repo *repository.ConversationRepository,
	userClient client.UserService,
) *ConversationService {

	return &ConversationService{
		repo:       repo,
		userClient: userClient,
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

	// Must be a number. Re-format it so "05" and "5" become the same ID.
	pid, err := strconv.ParseUint(participantID, 10, 64)
	if err != nil {
		return nil, ErrInvalidUserID
	}
	participantID = strconv.FormatUint(pid, 10)

	if userID == participantID {
		return nil, ErrCannotChatWithSelf
	}
	if err := s.checkUserExists(ctx, participantID); err != nil {
		return nil, err
	}

	// ... existingConversation, err := s.repo.FindPrivateConversation(...)
	// ... and everything below stays unchanged

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
func normalizeUserID(id string) (string, error) {
	n, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return "", ErrInvalidUserID
	}
	return strconv.FormatUint(n, 10), nil
}

func (s *ConversationService) checkUserExists(ctx context.Context, id string) error {
	lookupCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	exists, err := s.userClient.UserExists(lookupCtx, id)
	if err != nil {
		switch status.Code(err) {
		case codes.Unavailable, codes.DeadlineExceeded:
			return ErrUserServiceUnavailable
		}
		return err
	}
	if !exists {
		return ErrUserNotFound
	}
	return nil
}
func (s *ConversationService) CreateGroupConversation(
	ctx context.Context,
	creatorID string,
	req dto.CreateGroupConversationRequest,
) (*model.Conversation, error) {

	if creatorID == "" {
		return nil, ErrUserIDRequired
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrGroupNameRequired
	}

	if len(req.Participants) == 0 {
		return nil, ErrParticipantsRequired
	}

	// Validate every member BEFORE writing anything.
	seen := map[string]bool{creatorID: true}
	members := make([]string, 0, len(req.Participants))

	for _, raw := range req.Participants {
		id, err := normalizeUserID(raw)
		if err != nil {
			return nil, err
		}
		if seen[id] {
			continue
		}
		seen[id] = true

		if err := s.checkUserExists(ctx, id); err != nil {
			return nil, err
		}
		members = append(members, id)
	}

	if len(members) == 0 {
		return nil, ErrParticipantsRequired
	}

	conversation := &model.Conversation{
		Type: "group",
		Name: name,
	}

	err := s.repo.CreateConversation(ctx, conversation)
	if err != nil {
		return nil, err
	}

	participants := []model.ConversationParticipant{
		{ConversationID: conversation.ID, UserID: creatorID, Role: "admin"},
	}
	for _, id := range members {
		participants = append(participants, model.ConversationParticipant{
			ConversationID: conversation.ID,
			UserID:         id,
			Role:           "member",
		})
	}

	err = s.repo.CreateParticipants(ctx, participants)
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
func (s *ConversationService) SearchUsers(
	ctx context.Context,
	query string,
) (
	[]*pb.User,
	error,
) {

	res, err :=
		s.userClient.SearchUsers(
			ctx,
			query,
		)

	if err != nil {
		return nil, err
	}

	return res.Users, nil
}
