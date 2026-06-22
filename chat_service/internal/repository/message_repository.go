package repository

import (
	"chat_service/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{DB: db}

}
func (r *MessageRepository) CreateMessage(
	ctx context.Context,
	msg *model.Message,
) error {

	return r.DB.
		WithContext(ctx).
		Create(msg).
		Error
}
func (r *MessageRepository) GetMessageByID(
	ctx context.Context,
	id uint,
) (*model.Message, error) {

	var message model.Message

	err := r.DB.
		WithContext(ctx).
		First(&message, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &message, nil
}
func (r *MessageRepository) GetMessagesByConversation(
	ctx context.Context,
	conversationID uint,
	limit int,
	offset int,
) ([]model.Message, error) {

	var messages []model.Message

	err := r.DB.
		WithContext(ctx).
		Where(
			"conversation_id = ?",
			conversationID,
		).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&messages).
		Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}
func (r *MessageRepository) MarkRead(
	ctx context.Context,
	messageID uint,
) error {

	return r.DB.
		WithContext(ctx).
		Model(&model.Message{}).
		Where("id = ?", messageID).
		Update("status", "read").
		Error
}
func (r *MessageRepository) MarkDelivered(
	ctx context.Context,
	messageID uint,
) error {

	return r.DB.
		WithContext(ctx).
		Model(&model.Message{}).
		Where("id = ?", messageID).
		Update("status", "delivered").
		Error
}
func (r *MessageRepository) DeleteMessage(
	ctx context.Context,
	messageID uint,
) error {

	return r.DB.
		WithContext(ctx).
		Delete(
			&model.Message{},
			messageID,
		).
		Error
}
func (r *MessageRepository) GetUnreadMessages(
	ctx context.Context,
	conversationID uint,
) ([]model.Message, error) {

	var messages []model.Message

	err := r.DB.
		WithContext(ctx).
		Where(
			"conversation_id = ? AND status != ?",
			conversationID,
			"read",
		).
		Find(&messages).
		Error

	return messages, err
}
func (r *MessageRepository) CountMessages(
	ctx context.Context,
	conversationID uint,
) (int64, error) {

	var count int64

	err := r.DB.
		WithContext(ctx).
		Model(&model.Message{}).
		Where(
			"conversation_id = ?",
			conversationID,
		).
		Count(&count).
		Error

	return count, err
}
func (r *MessageRepository) UpdateMessage(
	ctx context.Context,
	message *model.Message,
) error {

	return r.DB.
		WithContext(ctx).
		Save(message).
		Error
}
