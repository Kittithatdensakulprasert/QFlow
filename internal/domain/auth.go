package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidOTP      = errors.New("invalid or expired OTP")
	ErrOTPAlreadyUsed  = errors.New("OTP has already been used or does not exist")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user with this phone number already exists")
	ErrPhoneRequired   = errors.New("phone number is required")
	ErrCodeRequired    = errors.New("code is required")
	ErrNameRequired    = errors.New("name is required")
	ErrUserIDRequired  = errors.New("user ID is required")
	ErrRoleNotAllowed  = errors.New("role changes are not allowed through this endpoint")
	ErrPhoneInvalid    = errors.New("phone number format is invalid")
	ErrOTPCodeRequired = errors.New("OTP code is required")
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Phone     string    `gorm:"uniqueIndex;not null" json:"phone"`
	Name      string    `json:"name"`
	Role      string    `gorm:"default:user" json:"role"` // user, provider, admin
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OTP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Phone     string    `gorm:"index;not null" json:"phone"`
	Code      string    `gorm:"not null" json:"code"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthRepository interface {
	CreateOTP(ctx context.Context, phone string) (*OTP, error)
	FindValidOTP(ctx context.Context, phone, code string) (*OTP, error)
	MarkOTPAsUsed(ctx context.Context, otpID uint) error
	DeleteExpiredOTPs(ctx context.Context, now time.Time) (int64, error)
	FindUserByPhone(ctx context.Context, phone string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	FindUserByID(ctx context.Context, id uint) (*User, error)
}

type AuthService interface {
	RequestOTP(ctx context.Context, phone string) (*OTP, error)
	VerifyOTP(ctx context.Context, phone, code string) (*User, string, error)
	RegisterUser(ctx context.Context, phone, name, role, otpCode string) (*User, string, error)
	GetUserProfile(ctx context.Context, userID uint) (*User, error)
	UpdateUserProfile(ctx context.Context, userID uint, name, role string) (*User, error)
}
