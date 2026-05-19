package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"qflow/internal/domain"
	"qflow/internal/jwt"
)

var phonePattern = regexp.MustCompile(`^0[0-9]{9}$`)

func isValidPhone(phone string) bool {
	return phonePattern.MatchString(phone)
}

type authService struct {
	authRepo   domain.AuthRepository
	jwtManager *jwt.JWTManager
	tokenGen   func(userID uint, phone, role string) (string, error)
}

func NewAuthService(authRepo domain.AuthRepository, jwtManager *jwt.JWTManager) domain.AuthService {
	return &authService{
		authRepo:   authRepo,
		jwtManager: jwtManager,
		tokenGen:   jwtManager.GenerateToken,
	}
}

func (s *authService) RequestOTP(ctx context.Context, phone string) (*domain.OTP, error) {
	if phone == "" {
		return nil, domain.ErrPhoneRequired
	}
	if !isValidPhone(phone) {
		return nil, domain.ErrPhoneInvalid
	}

	otp, err := s.authRepo.CreateOTP(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTP: %w", err)
	}

	return otp, nil
}

func (s *authService) VerifyOTP(ctx context.Context, phone, code string) (*domain.User, string, error) {
	if phone == "" {
		return nil, "", domain.ErrPhoneRequired
	}
	if !isValidPhone(phone) {
		return nil, "", domain.ErrPhoneInvalid
	}

	if code == "" {
		return nil, "", domain.ErrCodeRequired
	}

	otp, err := s.authRepo.FindValidOTP(ctx, phone, code)
	if err != nil {
		return nil, "", errors.New("invalid or expired OTP")
	}

	// Mark OTP as used
	if err := s.authRepo.MarkOTPAsUsed(ctx, otp.ID); err != nil {
		return nil, "", err
	}

	// Find user - if not found, create new user
	user, err := s.authRepo.FindUserByPhone(ctx, phone)
	if err != nil {
		// Auto-create user after OTP verification
		user = &domain.User{
			Phone: phone,
			Name:  phone,  // Default name to phone number
			Role:  "user", // Default role
		}

		if err := s.authRepo.CreateUser(ctx, user); err != nil {
			return nil, "", err
		}
	}

	// Generate JWT token
	token, err := s.tokenGen(user.ID, user.Phone, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *authService) RegisterUser(ctx context.Context, phone, name, role, otpCode string) (*domain.User, string, error) {
	if phone == "" {
		return nil, "", domain.ErrPhoneRequired
	}
	if !isValidPhone(phone) {
		return nil, "", domain.ErrPhoneInvalid
	}

	if name == "" {
		return nil, "", domain.ErrNameRequired
	}

	// Force role to be "user" for security - role escalation should be handled by admin-only endpoints
	role = "user"

	// SECURITY: Check if user already exists
	existingUser, err := s.authRepo.FindUserByPhone(ctx, phone)
	if err == nil && existingUser != nil {
		return nil, "", errors.New("user with this phone number already exists")
	}

	if otpCode == "" {
		return nil, "", domain.ErrOTPCodeRequired
	}

	// SECURITY: Check if there's a valid OTP for this phone
	otp, err := s.authRepo.FindValidOTP(ctx, phone, otpCode)
	if err != nil || otp == nil {
		return nil, "", errors.New("phone number not verified. Please request OTP first")
	}

	// Mark OTP as used to prevent reuse
	if err := s.authRepo.MarkOTPAsUsed(ctx, otp.ID); err != nil {
		return nil, "", err
	}

	// Create new user
	user := &domain.User{
		Phone: phone,
		Name:  name,
		Role:  role,
	}

	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		return nil, "", err
	}

	// Generate JWT token
	token, err := s.tokenGen(user.ID, user.Phone, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *authService) GetUserProfile(ctx context.Context, userID uint) (*domain.User, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *authService) UpdateUserProfile(ctx context.Context, userID uint, name, role string) (*domain.User, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if name != "" {
		user.Name = name
	}

	// Only allow role changes if current user is admin
	// For now, prevent any role changes through this endpoint
	// Role management should be handled by admin-only endpoints
	if role != "" && role != user.Role {
		return nil, errors.New("role changes are not allowed through this endpoint")
	}

	if err := s.authRepo.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
