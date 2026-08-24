package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"SDOBA/internal/model"
)

type AuthService struct {
	userService  *UserService
	tokenService *TokenService
}

func NewAuthService(
	userService *UserService,
	tokenService *TokenService,
) *AuthService {
	return &AuthService{
		userService:  userService,
		tokenService: tokenService,
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

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {

	user, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", errors.New("invalid password or email")
	}

	token, err := s.tokenService.Generate(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func (s *AuthService) GetCurrentUser(ctx context.Context, userID uint) (*model.User, error) {
	return s.userService.GetByID(ctx, userID)
}
