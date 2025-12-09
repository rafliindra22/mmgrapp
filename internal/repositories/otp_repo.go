package repositories

import (
	"context"
	"mmgrapp/internal/models"
	"time"

	"gorm.io/gorm"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *models.UserOTP) error
	FindValidOTP(ctx context.Context, userID int, purpose string) (*models.UserOTP, error)
	DeleteOTP(ctx context.Context, userID int, purpose string) error
	UpdateOTP(ctx context.Context, userID int, purpose string, hashedOTP string, otpExpires_time time.Time) error
}

type otpRepo struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) OTPRepository {
	return &otpRepo{db: db}
}

func (r *otpRepo) Create(ctx context.Context, otp *models.UserOTP) error {
	return r.db.Create(otp).Error
}

func (r *otpRepo) FindValidOTP(ctx context.Context, userID int, purpose string) (*models.UserOTP, error) {
	var otp models.UserOTP

	err := r.db.Where(
		"user_id = ? AND purpose = ? AND expires_at > ?",
		userID, purpose, time.Now(),
	).First(&otp).Error

	if err != nil {
		return nil, err
	}

	return &otp, nil
}

func (r *otpRepo) DeleteOTP(ctx context.Context, userID int, purpose string) error {
	return r.db.Where("user_id = ? AND purpose = ?", userID, purpose).Delete(&models.UserOTP{}).Error
}

func (r *otpRepo) UpdateOTP(ctx context.Context, userID int, purpose string, hashedOTP string, otpExpiresTime time.Time) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND purpose = ?", userID, purpose).
		Assign(map[string]interface{}{
			"user_id":    userID,
			"purpose":    purpose,
			"otp":        hashedOTP,
			"expires_at": otpExpiresTime,
		}).
		FirstOrCreate(&models.UserOTP{}).Error
}
