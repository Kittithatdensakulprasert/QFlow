package service

import (
	"fmt"

	"qflow/internal/domain"
	"qflow/internal/jwt"
)

type authService struct {
	authRepo   domain.AuthRepository
	jwtManager *jwt.JWTManager
}

func NewAuthService(authRepo domain.AuthRepository, jwtManager *jwt.JWTManager) domain.AuthService {
	return &authService{
		authRepo:   authRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) RequestOTP(phone string) (*domain.OTP, error) {
	if phone == "" {
		return nil, domain.ErrPhoneRequired
	}

	otp, err := s.authRepo.CreateOTP(phone)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTP: %w", err)
	}

	return otp, nil
}

func (s *authService) VerifyOTP(phone, code string) (*domain.User, string, error) {
	if phone == "" {
		return nil, "", domain.ErrPhoneRequired
	}

	if code == "" {
		return nil, "", domain.ErrCodeRequired
	}

	otp, err := s.authRepo.FindValidOTP(phone, code)
	if err != nil {
		return nil, "", domain.ErrInvalidOTP
	}

	// Mark OTP as used first to prevent reuse
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

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Phone, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *authService) RegisterUser(phone, name, role, otpCode string) (*domain.User, string, error) {
	if phone == "" {
		return nil, "", domain.ErrPhoneRequired
	}

	if name == "" {
		return nil, "", domain.ErrNameRequired
	}

	// Force role to be "user" for security - role escalation should be handled by admin-only endpoints
	role = "user"

	// SECURITY: Check if user already exists
	existingUser, err := s.authRepo.FindUserByPhone(phone)
	if err == nil && existingUser != nil {
		return nil, "", domain.ErrUserExists
	}

	if otpCode == "" {
		return nil, "", domain.ErrCodeRequired
	}

	// SECURITY: Check if there's a valid OTP for this phone
	otp, err := s.authRepo.FindValidOTP(phone, otpCode)
	if err != nil || otp == nil {
		return nil, "", domain.ErrInvalidOTP
	}

	// Mark OTP as used first to prevent reuse
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

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Phone, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *authService) GetUserProfile(userID uint) (*domain.User, error) {
	if userID == 0 {
		return nil, domain.ErrUserIDRequired
	}

	user, err := s.authRepo.FindUserByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (s *authService) UpdateUserProfile(userID uint, name, role string) (*domain.User, error) {
	if userID == 0 {
		return nil, domain.ErrUserIDRequired
	}

	user, err := s.authRepo.FindUserByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	if name != "" {
		user.Name = name
	}

	// Only allow role changes if current user is admin
	// For now, prevent any role changes through this endpoint
	// Role management should be handled by admin-only endpoints
	if role != "" && role != user.Role {
		return nil, domain.ErrRoleNotAllowed
	}

	if err := s.authRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
