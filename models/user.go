package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Phone     string    `json:"phone"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type OTPRequest struct {
	Phone string `json:"phone"`
}

type OTPVerify struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

type RegisterRequest struct {
	Phone     string `json:"phone"`
	OTP       string `json:"otp"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type OTP struct {
	Phone     string    `json:"phone"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expiresAt"`
}
