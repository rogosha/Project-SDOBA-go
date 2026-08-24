package service

import (
	"SDOBA/internal/model"
	"SDOBA/internal/repository"
	"context"
	"errors"
	"fmt"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUsernameExists = errors.New("username already exists")
	ErrEmailExists    = errors.New("email already exists")
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	exists, err := s.userRepository.FindByUsername(ctx, user.Username)
	if err != nil {
		//return fmt.Errorf("check username: %w", err)
	}

	if exists != nil {
		return fmt.Errorf("user already exists")
	}

	exists, err = s.userRepository.FindByEmail(ctx, user.Email)
	if err != nil {
		//return fmt.Errorf("check email: %w", err)
	}

	if exists != nil {
		return fmt.Errorf("user already exists")
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.userRepository.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by username: %w", err)
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, user *model.User) error {
	if err := s.userRepository.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	if err := s.userRepository.Delete(ctx, user); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
