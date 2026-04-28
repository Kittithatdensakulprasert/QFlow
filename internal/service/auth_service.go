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
	if phone == "" {
		return nil, "", errors.New("phone number is required")
	}

	if code == "" {
		return nil, "", errors.New("code is required")
	}

	otp, err := s.authRepo.FindValidOTP(phone, code)
	if err != nil {
		return nil, "", errors.New("invalid or expired OTP")
	}

	// Mark OTP as used
	if err := s.authRepo.MarkOTPAsUsed(otp.ID); err != nil {
		return nil, "", err
	}

	// Find user - if not found, create new user
	user, err := s.authRepo.FindUserByPhone(phone)
	if err != nil {
		// Auto-create user after OTP verification
		user = &domain.User{
			Phone: phone,
			Name:  phone,  // Default name to phone number
			Role:  "user", // Default role
		}

		if err := s.authRepo.CreateUser(user); err != nil {
			return nil, "", err
		}
	}

	token := fmt.Sprintf("token_%d_%d", user.ID, time.Now().Unix())

	return user, token, nil
}

func (s *authService) RegisterUser(phone, name, role string) (*domain.User, string, error) {
	if phone == "" {
		return nil, "", errors.New("phone number is required")
	}

	if name == "" {
		return nil, "", errors.New("name is required")
	}

	if role == "" {
		role = "user"
	}

	// SECURITY: Check if user already exists
	existingUser, err := s.authRepo.FindUserByPhone(phone)
	if err == nil && existingUser != nil {
		return nil, "", errors.New("user with this phone number already exists")
	}

	// SECURITY: Check if there's a valid OTP for this phone
	otp, err := s.authRepo.FindValidOTP(phone, "123456") // In production, this should be the actual OTP code
	if err != nil || otp == nil {
		return nil, "", errors.New("phone number not verified. Please request OTP first")
	}

	// Mark OTP as used to prevent reuse
	if err := s.authRepo.MarkOTPAsUsed(otp.ID); err != nil {
		return nil, "", err
	}

	// Create new user
	user := &domain.User{
		Phone: phone,
		Name:  name,
		Role:  role,
	}

	if err := s.authRepo.CreateUser(user); err != nil {
		return nil, "", err
	}

	// Generate token
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
