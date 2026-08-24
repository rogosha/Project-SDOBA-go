package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"SDOBA/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := r.db.
		WithContext(ctx).
		Where("username = ?", email).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, name string) (*model.User, error) {
	var user model.User

	if err := r.db.
		WithContext(ctx).
		Where("username = ?", name).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User

	if err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Delete(user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
