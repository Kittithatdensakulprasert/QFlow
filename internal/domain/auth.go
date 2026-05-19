package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidOTP     = errors.New("invalid or expired OTP")
	ErrOTPAlreadyUsed = errors.New("OTP has already been used or does not exist")
	ErrUserNotFound   = errors.New("user not found")
	ErrUserExists     = errors.New("user with this phone number already exists")
	ErrPhoneRequired  = errors.New("phone number is required")
	ErrCodeRequired   = errors.New("code is required")
	ErrNameRequired   = errors.New("name is required")
	ErrUserIDRequired = errors.New("user ID is required")
	ErrRoleNotAllowed = errors.New("role changes are not allowed through this endpoint")
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
	CreateOTP(phone string) (*OTP, error)
	FindValidOTP(phone, code string) (*OTP, error)
	MarkOTPAsUsed(otpID uint) error
	DeleteExpiredOTPs(now time.Time) (int64, error)
	FindUserByPhone(phone string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	FindUserByID(id uint) (*User, error)
}

type AuthService interface {
	RequestOTP(phone string) (*OTP, error)
	VerifyOTP(phone, code string) (*User, string, error)
	RegisterUser(phone, name, role, otpCode string) (*User, string, error)
	GetUserProfile(userID uint) (*User, error)
	UpdateUserProfile(userID uint, name, role string) (*User, error)
}
