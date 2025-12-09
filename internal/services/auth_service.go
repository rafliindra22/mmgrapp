package services

import (
	"context"
	"errors"
	"fmt"
	"mmgrapp/internal/dto"
	"mmgrapp/internal/models"
	"mmgrapp/internal/repositories"
	"mmgrapp/pkg/utils"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (interface{}, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email, otp, newPassword string) error
	RefreshToken(ctx context.Context, oldRefreshToken string) (map[string]interface{}, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	authRepo repositories.AuthRepository
	userRepo repositories.UserRepository
	otpRepo  repositories.OTPRepository
}

func NewAuthService(authRepo repositories.AuthRepository, userRepo repositories.UserRepository, otpRepo repositories.OTPRepository) AuthService {
	return &authService{
		authRepo: authRepo,
		userRepo: userRepo,
		otpRepo:  otpRepo,
	}
}

func (s *authService) Login(ctx context.Context, username, password string) (interface{}, error) {
	// 1. Coba login pakai email
	user, err := s.userRepo.FindByEmail(ctx, username)

	// 2. Jika tidak ditemukan sebagai email → coba sebagai username
	if err != nil {
		user, err = s.userRepo.FindByUsername(ctx, username)
		if err != nil {
			return nil, fmt.Errorf("email/username atau password salah")
		}
	}

	// 3. Cek Password
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("email/username atau password salah")
	}

	// 4. Cek Verified status
	if !user.IsVerified {
		return nil, errors.New("email belum diverifikasi")
	}

	// generate JWT
	accessToken, err := utils.GenerateAccessToken(user.ID, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	// generate refresh token
	refreshToken, err := utils.GenerateRefreshTokenJWT(user.ID)
	if err != nil {
		return nil, err
	}

	// save refresh token
	rt := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.authRepo.CreateRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	// set response
	userResponse := dto.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}

	return map[string]interface{}{
		"token_type":    "bearer",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          userResponse,
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	// cek user email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("email tidak ditemukan")
	}

	otp, hashedOTP, err := utils.GenerateOTP()
	if err != nil {
		return errors.New("gagal generate OTP")
	}

	if err := utils.SendOTP(user.Email, otp); err != nil {
		return err
	}

	return s.otpRepo.Create(ctx, &models.UserOTP{
		UserID:    user.ID,
		OTP:       hashedOTP,
		Purpose:   "password_reset",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})
}

func (s *authService) ResetPassword(ctx context.Context, email, otp, newPassword string) error {
	// cek user email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("email tidak ditemukan")
	}

	// cek valid OTP
	_, err = s.otpRepo.FindValidOTP(ctx, user.ID, "password_reset")
	if err != nil {
		return errors.New("OTP tidak valid atau kadaluarsa")
	}

	hashedPassword, _ := utils.HashPassword(newPassword)

	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	s.otpRepo.DeleteOTP(ctx, user.ID, "password_reset")

	return s.authRepo.UpdatePassword(ctx, user)
}

func (s *authService) RefreshToken(ctx context.Context, oldRefreshToken string) (map[string]interface{}, error) {
	// 1. Cek refresh token di DB
	rt, err := s.authRepo.FindRefreshTokenByToken(ctx, oldRefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, rt.UserID)

	// 2. Jika tidak ditemukan sebagai email → coba sebagai username
	if err != nil {
		user, err = s.userRepo.FindByUsername(ctx, user.Username)
		if err != nil {
			return nil, fmt.Errorf("Token tidak valid")
		}
	}

	// 2. Generate access token baru
	accessToken, err := utils.GenerateAccessToken(rt.UserID, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	// 3. Rotate refresh token (optional tapi lebih aman)
	newRefreshToken, err := utils.GenerateRefreshTokenJWT(rt.UserID)
	if err != nil {
		return nil, err
	}

	// 4. Revoke token lama
	if err := s.authRepo.RevokeRefreshToken(ctx, oldRefreshToken); err != nil {
		return nil, err
	}

	// 5. Simpan token baru
	newRT := &models.RefreshToken{
		UserID:    rt.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.authRepo.CreateRefreshToken(ctx, newRT); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token_type":    "bearer",
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	// Revoke token di DB
	err := s.authRepo.RevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	return nil
}
