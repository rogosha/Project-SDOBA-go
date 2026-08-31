package service

import (
	"context"
	"errors"
	"fmt"

	"SDOBA/internal/model"
	"SDOBA/internal/repository"
)

var (
	ErrEmptyMessage          = errors.New("message content is empty")
	ErrNotConversationMember = errors.New("user is not a conversation member")
	ErrMessageNotFound       = errors.New("message not found")
)

type MessageService struct {
	messageRepository      *repository.MessageRepository
	conversationRepository *repository.ConversationRepository
}

func NewMessageService(
	messageRepository *repository.MessageRepository,
	conversationRepository *repository.ConversationRepository,
) *MessageService {
	return &MessageService{
		messageRepository:      messageRepository,
		conversationRepository: conversationRepository,
	}
}

func (s *MessageService) Create(ctx context.Context, conversationID uint, userID uint, content string,
) (*model.Message, error) {

	if content == "" {
		return nil, ErrEmptyMessage
	}

	isMember, err := s.conversationRepository.IsMember(
		ctx,
		conversationID,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("check conversation member: %w", err)
	}

	if !isMember {
		return nil, ErrNotConversationMember
	}

	message := &model.Message{
		ConversationID: conversationID,
		SenderID:       userID,
		Content:        content,
	}

	if err := s.messageRepository.Create(ctx, message); err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	return message, nil
}

func (s *MessageService) GetByConversationID(ctx context.Context, conversationID uint, userID uint,
) ([]model.Message, error) {

	isMember, err := s.conversationRepository.IsMember(
		ctx,
		conversationID,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("check conversation member: %w", err)
	}

	if !isMember {
		return nil, ErrNotConversationMember
	}

	return s.messageRepository.GetByConversationID(
		ctx,
		conversationID,
	)
}

func (s *MessageService) Delete(ctx context.Context, id uint, userID uint) error {
	if err := s.messageRepository.DeleteByID(ctx, id, userID); err != nil {
		if errors.Is(err, repository.ErrMessageNotFound) {
			return ErrMessageNotFound
		}

		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}
