package service

import (
	"context"
	"errors"
	"fmt"

	"SDOBA/internal/model"
	"SDOBA/internal/repository"
)

var ErrMessageNotFound = errors.New("message not found")

type MessageService struct {
	repository *repository.MessageRepository
}

func NewMessageService(
	repository *repository.MessageRepository,
) *MessageService {
	return &MessageService{
		repository: repository,
	}
}

func (s *MessageService) Create(
	ctx context.Context,
	conversationID uint,
	userID uint,
	content string,
) (*model.Message, error) {

	if content == "" {
		return nil, errors.New("message content is empty")
	}

	message := &model.Message{
		ConversationID: conversationID,
		SenderID:       userID,
		Content:        content,
	}

	if err := s.repository.Create(ctx, message); err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	return message, nil
}

func (s *MessageService) GetByConversationID(
	ctx context.Context,
	conversationID uint,
) ([]model.Message, error) {

	return s.repository.GetByConversationID(ctx, conversationID)
}

func (s *MessageService) Delete(
	ctx context.Context,
	id uint,
) error {

	if err := s.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrMessageNotFound) {
			return ErrMessageNotFound
		}

		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}
