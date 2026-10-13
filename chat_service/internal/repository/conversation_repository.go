package repository

import (
	"chat_service/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type ConversationRepository struct {
	DB *gorm.DB
}

func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{
		DB: db,
	}
}
func (r *ConversationRepository) CreateConversation(
	ctx context.Context,
	conversation *model.Conversation,
) error {

	return r.DB.
		WithContext(ctx).
		Create(conversation).
		Error
}

func (r *ConversationRepository) GetConversationByID(
	ctx context.Context,
	id uint,
) (*model.Conversation, error) {

	var conversation model.Conversation

	err := r.DB.
		WithContext(ctx).
		Preload("Participants").
		First(&conversation, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}
func (r *ConversationRepository) FindPrivateConversation(
	ctx context.Context,
	user1 string,
	user2 string,
) (*model.Conversation, error) {

	var conversation model.Conversation

	err := r.DB.
		WithContext(ctx).
		Joins(
			"JOIN conversation_participants cp1 ON cp1.conversation_id = conversations.id",
		).
		Joins(
			"JOIN conversation_participants cp2 ON cp2.conversation_id = conversations.id",
		).
		Where("cp1.user_id = ?", user1).
		Where("cp2.user_id = ?", user2).
		Where("conversations.type = ?", "private").
		First(&conversation).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}
func (r *ConversationRepository) GetUserConversations(
	ctx context.Context,
	userID string,
) ([]model.Conversation, error) {

	var conversations []model.Conversation

	err := r.DB.
		WithContext(ctx).
		Joins(
			"JOIN conversation_participants cp ON cp.conversation_id = conversations.id",
		).
		Where(
			"cp.user_id = ?",
			userID,
		).
		Preload("Participants").
		Find(&conversations).
		Error

	return conversations, err
}
func (r *ConversationRepository) UpdateLastMessageID(
	ctx context.Context,
	conversationID uint,
	messageID uint,
) error {

	return r.DB.
		WithContext(ctx).
		Model(&model.Conversation{}).
		Where("id = ?", conversationID).
		Update("last_message_id", messageID).
		Error
}
func (r *ConversationRepository) DeleteConversation(
	ctx context.Context,
	id uint,
) error {

	return r.DB.
		WithContext(ctx).
		Delete(
			&model.Conversation{},
			id,
		).
		Error
}
func (r *ConversationRepository) CreateParticipants(
	ctx context.Context,
	participants []model.ConversationParticipant,
) error {

	return r.DB.
		WithContext(ctx).
		Create(&participants).
		Error
}
func (r *ConversationRepository) GetParticipants(
	ctx context.Context,
	conversationID uint,
) ([]model.ConversationParticipant, error) {

	var participants []model.ConversationParticipant

	err := r.DB.
		WithContext(ctx).
		Where(
			"conversation_id = ?",
			conversationID,
		).
		Find(&participants).
		Error

	return participants, err
}
func (r *ConversationRepository) IsParticipant(
	ctx context.Context,
	conversationID uint,
	userID string,
) (bool, error) {

	var count int64

	err := r.DB.
		WithContext(ctx).
		Model(
			&model.ConversationParticipant{},
		).
		Where(
			"conversation_id = ? AND user_id = ?",
			conversationID,
			userID,
		).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
func (r *ConversationRepository) AddParticipant(
	ctx context.Context,
	participant *model.ConversationParticipant,
) error {

	return r.DB.
		WithContext(ctx).
		Create(participant).
		Error
}
func (r *ConversationRepository) RemoveParticipant(
	ctx context.Context,
	conversationID uint,
	userID string,
) error {

	return r.DB.
		WithContext(ctx).
		Delete(
			&model.ConversationParticipant{},
			"conversation_id = ? AND user_id = ?",
			conversationID,
			userID,
		).
		Error
}
