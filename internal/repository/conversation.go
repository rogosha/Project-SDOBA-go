package repository

import (
	"context"
	"errors"

	"SDOBA/internal/model"

	"gorm.io/gorm"
)

var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrUserAlreadyMember    = errors.New("user already a member")
)

type ConversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(
	ctx context.Context,
	userIDs []uint,
) (*model.Conversation, error) {

	existing, err := r.FindByMembers(ctx, userIDs)

	if err == nil {
		return existing, nil
	}

	if !errors.Is(err, ErrConversationNotFound) {
		return nil, err
	}

	var conversation model.Conversation

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var count int64

		if err := tx.Model(&model.User{}).
			Where("id IN ?", userIDs).
			Count(&count).Error; err != nil {
			return err
		}

		if count != int64(len(userIDs)) {
			return ErrUserNotFound
		}

		if err := tx.Create(&conversation).Error; err != nil {
			return err
		}

		for _, userID := range userIDs {
			member := model.ConversationMember{
				ConversationID: conversation.ID,
				UserID:         userID,
			}

			if err := tx.Create(&member).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *ConversationRepository) GetUserConversations(ctx context.Context, userID uint) ([]model.Conversation, error) {
	var conversations []model.Conversation

	if err := r.db.WithContext(ctx).
		Joins("JOIN conversation_members ON conversation_members.conversation_id = conversations.id").
		Where("conversation_members.user_id = ?", userID).
		Order("conversations.updated_at DESC").
		Find(&conversations).Error; err != nil {
		return nil, err
	}

	return conversations, nil
}

func (r *ConversationRepository) IsMember(ctx context.Context, conversationID uint, userID uint) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.ConversationMember{}).
		Where(
			"conversation_id = ? AND user_id = ?",
			conversationID,
			userID,
		).
		Count(&count).Error

	return count > 0, err
}

func (r *ConversationRepository) FindByMembers(ctx context.Context, userIDs []uint) (*model.Conversation, error) {

	var conversation model.Conversation

	err := r.db.WithContext(ctx).
		Model(&model.Conversation{}).
		Joins("JOIN conversation_members ON conversation_members.conversation_id = conversations.id").
		Where("conversation_members.user_id IN ?", userIDs).
		Group("conversations.id").
		Having("COUNT(DISTINCT conversation_members.user_id) = ?", len(userIDs)).
		Having(`
			(SELECT COUNT(*)
			 FROM conversation_members cm
			 WHERE cm.conversation_id = conversations.id) = ?
		`, len(userIDs)).
		First(&conversation).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrConversationNotFound
	}

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *ConversationRepository) AddMember(ctx context.Context, conversationID uint, userID uint) error {
	var user model.User

	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	var member model.ConversationMember

	err := r.db.WithContext(ctx).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		First(&member).Error

	if err == nil {
		return ErrUserAlreadyMember
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	member = model.ConversationMember{
		ConversationID: conversationID,
		UserID:         userID,
	}

	return r.db.WithContext(ctx).Create(&member).Error
}

func (r *ConversationRepository) GetByID(ctx context.Context, id uint) (*model.Conversation, error) {
	var conversation model.Conversation

	err := r.db.WithContext(ctx).
		Preload("Members").
		First(&conversation, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrConversationNotFound
	}

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *ConversationRepository) RemoveMember(ctx context.Context, conversationID uint, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Delete(&model.ConversationMember{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *ConversationRepository) CountMembers(ctx context.Context, conversationID uint) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.ConversationMember{}).
		Where("conversation_id = ?", conversationID).
		Count(&count).Error

	return count, err
}

func (r *ConversationRepository) Delete(ctx context.Context, conversationID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("conversation_id = ?", conversationID).
			Delete(&model.Message{}).Error; err != nil {
			return err
		}

		if err := tx.Where("conversation_id = ?", conversationID).
			Delete(&model.ConversationMember{}).Error; err != nil {
			return err
		}

		result := tx.Delete(&model.Conversation{}, conversationID)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return ErrConversationNotFound
		}

		return nil
	})
}
