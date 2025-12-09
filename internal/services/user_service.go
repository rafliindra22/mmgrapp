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

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, username, email, password string) (*models.User, error)
	VerifyEmailOTP(ctx context.Context, email, otp string) error
	ResendOTP(ctx context.Context, email string, purpose string) error
	GetUserByID(ctx context.Context, id int) (interface{}, error)
}

type userService struct {
	repo    repositories.UserRepository
	otpRepo repositories.OTPRepository
}

func NewUserService(userRepo repositories.UserRepository, otpRepo repositories.OTPRepository) UserService {
	return &userService{
		repo:    userRepo,
		otpRepo: otpRepo,
	}
}

func (s *userService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	// check duplikasi
	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		return nil, errors.New("email sudah digunakan")
	}

	if _, err := s.repo.FindByUsername(ctx, username); err == nil {
		return nil, errors.New("username sudah digunakan")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal meng-hash password")
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	// create user
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// generate OTP
	otp, hashedOTP, err := utils.GenerateOTP()
	if err != nil {
		return nil, err
	}

	// delete OTP lama (optional)
	s.otpRepo.DeleteOTP(ctx, user.ID, "email_verification")

	// simpan OTP
	err = s.otpRepo.Create(ctx, &models.UserOTP{
		UserID:    user.ID,
		OTP:       hashedOTP,
		Purpose:   "email_verification",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})
	if err != nil {
		return nil, err
	}

	// kirim email
	if err := utils.SendOTP(user.Email, otp); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) VerifyEmailOTP(ctx context.Context, email, otp string) error {
	// 1. Get user
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("email tidak ditemukan")
	}

	// 2. Get latest OTP for verification
	storedOTP, err := s.otpRepo.FindValidOTP(ctx, user.ID, "email_verification")
	if err != nil {
		return errors.New("OTP tidak ditemukan atau sudah kadaluarsa")
	}

	// 3. Check expiry
	if time.Now().After(storedOTP.ExpiresAt) {
		return errors.New("OTP sudah kadaluarsa")
	}

	// 4. Compare OTP
	if err := bcrypt.CompareHashAndPassword([]byte(storedOTP.OTP), []byte(otp)); err != nil {
		return errors.New("OTP salah")
	}

	// 5. Update user as verified
	_, err = s.repo.VerifyUser(ctx, email)
	if err != nil {
		return err
	}

	// 6. Hapus OTP setelah dipakai
	s.otpRepo.DeleteOTP(ctx, user.ID, "email_verification")

	return nil
}

func (s *userService) ResendOTP(ctx context.Context, email string, purpose string) error {
	// cek purpose
	const (
		PurposeEmailVerification = "email_verification"
		PurposePasswordReset     = "password_reset"
	)

	if purpose != PurposeEmailVerification && purpose != PurposePasswordReset {
		return fmt.Errorf("invalid purpose")
	}

	// cek user
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan")
	}

	if user.IsVerified {
		return fmt.Errorf("email sudah diverifikasi, silakan login")
	}

	// Generate OTP baru
	otp, hashedOTP, err := utils.GenerateOTP()
	if err != nil {
		return err
	}
	otpExpiresTime := time.Now().Add(5 * time.Minute)

	// Update OTP di DB
	if err := s.otpRepo.UpdateOTP(ctx, user.ID, purpose, hashedOTP, otpExpiresTime); err != nil {
		return err
	}

	// Kirim OTP via email
	if err := utils.SendOTP(user.Email, otp); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (interface{}, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userResponse := dto.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}

	return userResponse, nil
}
