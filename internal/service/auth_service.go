package service

import (
	"errors"
	"fmt"
	"time"

	"qflow/internal/domain"
)

type authService struct {
	authRepo domain.AuthRepository
}

func NewAuthService(authRepo domain.AuthRepository) domain.AuthService {
	return &authService{authRepo: authRepo}
}

func (s *authService) RequestOTP(phone string) (*domain.OTP, error) {
	if phone == "" {
		return nil, errors.New("phone number is required")
	}

	otp, err := s.authRepo.CreateOTP(phone)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTP: %w", err)
	}

	return otp, nil
}

func (s *authService) VerifyOTP(phone, code string) (*domain.User, string, error) {
	if phone == "" || code == "" {
		return nil, "", errors.New("phone and code are required")
	}

	otp, err := s.authRepo.FindValidOTP(phone, code)
	if err != nil {
		return nil, "", errors.New("invalid or expired OTP")
	}

	if err := s.authRepo.MarkOTPAsUsed(otp.ID); err != nil {
		return nil, "", fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	user, err := s.authRepo.FindUserByPhone(phone)
	if err != nil {
		return nil, "", errors.New("user not found")
	}

	token := fmt.Sprintf("token_%d_%d", user.ID, time.Now().Unix())

	return user, token, nil
}

func (s *authService) RegisterUser(phone, name, role string) (*domain.User, string, error) {
	if phone == "" || name == "" {
		return nil, "", errors.New("phone and name are required")
	}

	if role == "" {
		role = "user"
	}

	existingUser, err := s.authRepo.FindUserByPhone(phone)
	if err == nil && existingUser != nil {
		return nil, "", errors.New("user already exists")
	}

	user := &domain.User{
		Phone: phone,
		Name:  name,
		Role:  role,
	}

	if err := s.authRepo.CreateUser(user); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	token := fmt.Sprintf("token_%d_%d", user.ID, time.Now().Unix())

	return user, token, nil
}

func (s *authService) GetUserProfile(userID uint) (*domain.User, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	user, err := s.authRepo.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *authService) UpdateUserProfile(userID uint, name, role string) (*domain.User, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	user, err := s.authRepo.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if name != "" {
		user.Name = name
	}

	if role != "" {
		user.Role = role
	}

	if err := s.authRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
