package service

import (
	"errors"
	"testing"
	"time"

	"qflow/internal/domain"
)

type mockAuthRepository struct {
	otp           *domain.OTP
	users         map[string]*domain.User
	usersByID     map[uint]*domain.User
	otpUsed       bool
	createError   error
	findError     error
	createUserErr error
	updateUserErr error
}

func newMockAuthRepository() *mockAuthRepository {
	return &mockAuthRepository{
		users:     make(map[string]*domain.User),
		usersByID: make(map[uint]*domain.User),
	}
}

func (m *mockAuthRepository) CreateOTP(phone string) (*domain.OTP, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.otp = &domain.OTP{
		ID:        1,
		Phone:     phone,
		Code:      "123456",
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Used:      false,
	}
	return m.otp, nil
}

func (m *mockAuthRepository) FindValidOTP(phone, code string) (*domain.OTP, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	if m.otp == nil || m.otp.Phone != phone || m.otp.Code != code || m.otpUsed || m.otp.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("otp not found")
	}
	return m.otp, nil
}

func (m *mockAuthRepository) MarkOTPAsUsed(otpID uint) error {
	if m.otp != nil && m.otp.ID == otpID {
		m.otpUsed = true
	}
	return nil
}

func (m *mockAuthRepository) FindUserByPhone(phone string) (*domain.User, error) {
	if user, exists := m.users[phone]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockAuthRepository) CreateUser(user *domain.User) error {
	if m.createUserErr != nil {
		return m.createUserErr
	}
	user.ID = uint(len(m.users) + 1)
	m.users[user.Phone] = user
	m.usersByID[user.ID] = user
	return nil
}

func (m *mockAuthRepository) UpdateUser(user *domain.User) error {
	if m.updateUserErr != nil {
		return m.updateUserErr
	}
	m.users[user.Phone] = user
	m.usersByID[user.ID] = user
	return nil
}

func (m *mockAuthRepository) FindUserByID(id uint) (*domain.User, error) {
	if user, exists := m.usersByID[id]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func TestAuthService_RequestOTP(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"valid phone", "1234567890", false},
		{"empty phone", "", true},
		{"create error", "1234567890", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockAuthRepository()
			if tt.name == "create error" {
				mock.createError = errors.New("database error")
			}

			service := NewAuthService(mock)
			otp, err := service.RequestOTP(tt.phone)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if otp == nil {
				t.Errorf("Expected OTP but got nil")
				return
			}

			if otp.Phone != tt.phone {
				t.Errorf("Expected phone %s but got %s", tt.phone, otp.Phone)
			}
		})
	}
}

func TestAuthService_VerifyOTP(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		code    string
		setup   func(*mockAuthRepository)
		wantErr bool
	}{
		{"valid OTP", "1234567890", "123456", func(m *mockAuthRepository) {
			user := &domain.User{ID: 1, Phone: "1234567890", Name: "Test User"}
			m.users["1234567890"] = user
			m.usersByID[1] = user
			m.otp = &domain.OTP{ID: 1, Phone: "1234567890", Code: "123456", Used: false, ExpiresAt: time.Now().Add(5 * time.Minute)}
		}, false},
		{"invalid phone", "", "123456", func(m *mockAuthRepository) {}, true},
		{"invalid code", "1234567890", "", func(m *mockAuthRepository) {}, true},
		{"OTP not found", "1234567890", "123456", func(m *mockAuthRepository) {}, true},
		{"user not found", "1234567890", "123456", func(m *mockAuthRepository) {
			m.otp = &domain.OTP{ID: 1, Phone: "1234567890", Code: "123456", Used: false, ExpiresAt: time.Now().Add(5 * time.Minute)}
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockAuthRepository()
			tt.setup(mock)

			service := NewAuthService(mock)
			user, token, err := service.VerifyOTP(tt.phone, tt.code)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("Expected user but got nil")
				return
			}

			if token == "" {
				t.Errorf("Expected token but got empty string")
				return
			}

			if !mock.otpUsed {
				t.Errorf("Expected OTP to be marked as used")
			}
		})
	}
}

func TestAuthService_RegisterUser(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		userName string
		role     string
		setup    func(*mockAuthRepository)
		wantErr  bool
	}{
		{"valid registration", "1234567890", "Test User", "user", func(m *mockAuthRepository) {}, false},
		{"empty phone", "", "Test User", "user", func(m *mockAuthRepository) {}, true},
		{"empty name", "1234567890", "", "user", func(m *mockAuthRepository) {}, true},
		{"user exists", "1234567890", "Test User", "user", func(m *mockAuthRepository) {
			user := &domain.User{ID: 1, Phone: "1234567890", Name: "Existing User"}
			m.users["1234567890"] = user
		}, true},
		{"create error", "1234567890", "Test User", "user", func(m *mockAuthRepository) {
			m.createUserErr = errors.New("database error")
		}, true},
		{"default role", "1234567890", "Test User", "", func(m *mockAuthRepository) {}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockAuthRepository()
			tt.setup(mock)

			service := NewAuthService(mock)
			user, token, err := service.RegisterUser(tt.phone, tt.userName, tt.role)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("Expected user but got nil")
				return
			}

			if token == "" {
				t.Errorf("Expected token but got empty string")
				return
			}

			if user.Phone != tt.phone {
				t.Errorf("Expected phone %s but got %s", tt.phone, user.Phone)
			}

			if user.Name != tt.userName {
				t.Errorf("Expected name %s but got %s", tt.userName, user.Name)
			}

			expectedRole := tt.role
			if expectedRole == "" {
				expectedRole = "user"
			}
			if user.Role != expectedRole {
				t.Errorf("Expected role %s but got %s", expectedRole, user.Role)
			}
		})
	}
}

func TestAuthService_GetUserProfile(t *testing.T) {
	tests := []struct {
		name    string
		userID  uint
		setup   func(*mockAuthRepository)
		wantErr bool
	}{
		{"valid user", 1, func(m *mockAuthRepository) {
			user := &domain.User{ID: 1, Phone: "1234567890", Name: "Test User"}
			m.usersByID[1] = user
		}, false},
		{"invalid user ID", 0, func(m *mockAuthRepository) {}, true},
		{"user not found", 999, func(m *mockAuthRepository) {}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockAuthRepository()
			tt.setup(mock)

			service := NewAuthService(mock)
			user, err := service.GetUserProfile(tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("Expected user but got nil")
				return
			}

			if user.ID != tt.userID {
				t.Errorf("Expected user ID %d but got %d", tt.userID, user.ID)
			}
		})
	}
}

func TestAuthService_UpdateUserProfile(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		userName string
		role     string
		setup    func(*mockAuthRepository)
		wantErr  bool
	}{
		{"valid update", 1, "Updated Name", "admin", func(m *mockAuthRepository) {
			user := &domain.User{ID: 1, Phone: "1234567890", Name: "Test User", Role: "user"}
			m.usersByID[1] = user
		}, false},
		{"invalid user ID", 0, "Updated Name", "admin", func(m *mockAuthRepository) {}, true},
		{"user not found", 999, "Updated Name", "admin", func(m *mockAuthRepository) {}, true},
		{"update error", 1, "Updated Name", "admin", func(m *mockAuthRepository) {
			user := &domain.User{ID: 1, Phone: "1234567890", Name: "Test User", Role: "user"}
			m.usersByID[1] = user
			m.updateUserErr = errors.New("database error")
		}, true},
		{"partial update name", 1, "Updated Name", "", func(m *mockAuthRepository) {
			user := &domain.User{ID: 1, Phone: "1234567890", Name: "Test User", Role: "user"}
			m.usersByID[1] = user
		}, false},
		{"partial update role", 1, "", "admin", func(m *mockAuthRepository) {
			user := &domain.User{ID: 1, Phone: "1234567890", Name: "Test User", Role: "user"}
			m.usersByID[1] = user
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockAuthRepository()
			tt.setup(mock)

			service := NewAuthService(mock)
			user, err := service.UpdateUserProfile(tt.userID, tt.userName, tt.role)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("Expected user but got nil")
				return
			}

			if user.ID != tt.userID {
				t.Errorf("Expected user ID %d but got %d", tt.userID, user.ID)
			}

			if tt.userName != "" && user.Name != tt.userName {
				t.Errorf("Expected name %s but got %s", tt.userName, user.Name)
			}

			if tt.role != "" && user.Role != tt.role {
				t.Errorf("Expected role %s but got %s", tt.role, user.Role)
			}
		})
	}
}
