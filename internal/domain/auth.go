package domain

import "time"

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
	// TODO: define methods
}

type AuthService interface {
	// TODO: define methods
}
