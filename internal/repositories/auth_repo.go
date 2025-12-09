package repositories

import (
	"context"
	"errors"
	"mmgrapp/internal/models"
	"time"

	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error
	UpdatePassword(ctx context.Context, user *models.User) error
	FindRefreshTokenByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
}

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepo{db}
}

func (r *authRepo) CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	return r.db.WithContext(ctx).Create(refreshToken).Error
}

func (r *authRepo) UpdatePassword(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *authRepo) FindRefreshTokenByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token = ? AND is_revoked = false AND expires_at > ?", token, time.Now()).
		First(&rt).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("refresh token not found or invalid")
		}
		return nil, err
	}

	return &rt, nil
}

func (r *authRepo) RevokeRefreshToken(ctx context.Context, token string) error {
	result := r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("refresh token not found")
	}

	return nil
}
