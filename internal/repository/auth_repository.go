package repository

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"qflow/internal/domain"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateOTP(phone string) (*domain.OTP, error) {
	// Generate cryptographically secure 6-digit OTP
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return nil, fmt.Errorf("failed to generate secure OTP: %w", err)
	}
	code := fmt.Sprintf("%06d", n.Int64())

	otp := &domain.OTP{
		Phone:     phone,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Used:      false,
	}

	if err := r.db.Create(otp).Error; err != nil {
		return nil, err
	}

	return otp, nil
}

func (r *authRepository) FindValidOTP(phone, code string) (*domain.OTP, error) {
	var otp domain.OTP
	err := r.db.Where("phone = ? AND code = ? AND used = ? AND expires_at > ?",
		phone, code, false, time.Now()).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *authRepository) MarkOTPAsUsed(otpID uint) error {
	result := r.db.Model(&domain.OTP{}).Where("id = ? AND used = ?", otpID, false).Update("used", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("OTP has already been used or does not exist")
	}
	return nil
}

func (r *authRepository) FindUserByPhone(phone string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) UpdateUser(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *authRepository) FindUserByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
