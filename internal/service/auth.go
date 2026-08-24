package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"SDOBA/internal/model"
)

type AuthService struct {
	userService *UserService
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (s *AuthService) Register(ctx context.Context, username string, email string, password string,
) (*model.User, error) {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	if err := s.userService.Create(ctx, user); err != nil {
		if errors.Is(err, ErrUsernameExists) {
			return nil, ErrUsernameExists
		}

		if errors.Is(err, ErrEmailExists) {
			return nil, ErrEmailExists
		}

		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*model.User, error) {

	user, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
