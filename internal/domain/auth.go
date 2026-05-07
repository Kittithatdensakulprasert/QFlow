package domain

import (
	"time"
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
	Phone     string    `gorm:"not null" json:"phone"`
	Code      string    `gorm:"not null" json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
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
