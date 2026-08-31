package repository

import (
	"context"
	"errors"

	"SDOBA/internal/model"

	"gorm.io/gorm"
)

var ErrMessageNotFound = errors.New("message not found")

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(
	ctx context.Context,
	message *model.Message,
) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *MessageRepository) GetByConversationID(
	ctx context.Context,
	conversationID uint,
) ([]model.Message, error) {

	var messages []model.Message

	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepository) GetByID(
	ctx context.Context,
	id uint,
) (*model.Message, error) {

	var message model.Message

	err := r.db.WithContext(ctx).
		First(&message, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrMessageNotFound
	}

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *MessageRepository) Delete(
	ctx context.Context,
	id uint,
) error {

	result := r.db.WithContext(ctx).
		Delete(&model.Message{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrMessageNotFound
	}

	return nil
}

func (r *MessageRepository) DeleteByID(ctx context.Context, id uint, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND sender_id = ?", id, userID).
		Delete(&model.Message{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrMessageNotFound
	}

	return nil
}
