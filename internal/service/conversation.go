package service

import (
	"context"
	"errors"
	"fmt"

	"SDOBA/internal/model"
	"SDOBA/internal/repository"
)

var (
	ErrConversationMinMembers   = errors.New("conversation must have at least 2 members")
	ErrInvalidUserID            = errors.New("invalid user id")
	ErrDuplicateUser            = errors.New("duplicate user")
	ErrConversationNotFound     = errors.New("conversation not found")
	ErrConversationAccessDenied = errors.New("conversation access denied")
)

type ConversationService struct {
	repository *repository.ConversationRepository
}

func NewConversationService(
	repository *repository.ConversationRepository,
) *ConversationService {
	return &ConversationService{
		repository: repository,
	}
}

func (s *ConversationService) Create(ctx context.Context, userIDs []uint) (*model.Conversation, error) {

	if len(userIDs) < 2 {
		return nil, ErrConversationMinMembers
	}

	unique := make(map[uint]struct{})

	for _, userID := range userIDs {
		if userID == 0 {
			return nil, ErrInvalidUserID
		}

		if _, exists := unique[userID]; exists {
			return nil, ErrDuplicateUser
		}

		unique[userID] = struct{}{}
	}

	conversation, err := s.repository.Create(ctx, userIDs)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("create conversation: %w", err)
	}

	return conversation, nil
}

func (s *ConversationService) GetByID(ctx context.Context, conversationID uint, userID uint,
) (*model.Conversation, error) {

	isMember, err := s.repository.IsMember(
		ctx,
		conversationID,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("check conversation member: %w", err)
	}

	if !isMember {
		return nil, ErrConversationNotFound
	}

	return s.repository.GetByID(ctx, conversationID)
}

func (s *ConversationService) GetUserConversations(ctx context.Context, userID uint) ([]model.Conversation, error) {

	return s.repository.GetUserConversations(ctx, userID)
}

func (s *ConversationService) AddMember(ctx context.Context, conversationID uint, userID uint, requesterID uint) error {
	isMember, err := s.repository.IsMember(ctx, conversationID, requesterID)
	if err != nil {
		return err
	}

	if !isMember {
		return ErrConversationAccessDenied
	}

	return s.repository.AddMember(ctx, conversationID, userID)
}

func (s *ConversationService) RemoveMember(ctx context.Context, conversationID uint, userID uint, requesterID uint) error {
	isMember, err := s.repository.IsMember(ctx, conversationID, requesterID)
	if err != nil {
		return err
	}

	if !isMember {
		return ErrConversationAccessDenied
	}

	targetIsMember, err := s.repository.IsMember(ctx, conversationID, userID)
	if err != nil {
		return err
	}

	if !targetIsMember {
		return repository.ErrUserNotFound
	}

	count, err := s.repository.CountMembers(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("count conversation members: %w", err)
	}

	if count <= 2 {
		if err := s.repository.Delete(ctx, conversationID); err != nil {
			return fmt.Errorf("delete conversation: %w", err)
		}

		return nil
	}

	if err := s.repository.RemoveMember(ctx, conversationID, userID); err != nil {
		return fmt.Errorf("remove conversation member: %w", err)
	}

	return nil
}
