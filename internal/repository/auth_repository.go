package repository

import (
	"context"
	"crypto/rand"
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

func (r *authRepository) CreateOTP(ctx context.Context, phone string) (*domain.OTP, error) {
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

	if err := r.db.WithContext(ctx).Create(otp).Error; err != nil {
		return nil, err
	}

	return otp, nil
}

func (r *authRepository) FindValidOTP(ctx context.Context, phone, code string) (*domain.OTP, error) {
	var otp domain.OTP
	err := r.db.WithContext(ctx).Where("phone = ? AND code = ? AND used = ? AND expires_at > ?",
		phone, code, false, time.Now()).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *authRepository) MarkOTPAsUsed(ctx context.Context, otpID uint) error {
	result := r.db.WithContext(ctx).Model(&domain.OTP{}).Where("id = ? AND used = ?", otpID, false).Update("used", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrOTPAlreadyUsed
	}
	return nil
}

func (r *authRepository) DeleteExpiredOTPs(ctx context.Context, now time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Where("expires_at <= ?", now).Delete(&domain.OTP{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (r *authRepository) FindUserByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *authRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *authRepository) FindUserByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
