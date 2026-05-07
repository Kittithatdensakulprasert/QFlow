package service

import (
	"testing"
	"time"

	"qflow/internal/domain"
)

func TestOTPCleanupJob_CleanupExpiredOTPs(t *testing.T) {
	mock := newMockAuthRepository()
	now := time.Now()
	mock.otp = &domain.OTP{
		ID:        1,
		Phone:     "1234567890",
		Code:      "123456",
		ExpiresAt: now.Add(-time.Minute),
		Used:      false,
	}

	job := NewOTPCleanupJob(mock, time.Hour, nil)
	deleted, err := job.CleanupExpiredOTPs(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if deleted != 1 {
		t.Fatalf("expected 1 deleted OTP, got %d", deleted)
	}
	if mock.otp != nil {
		t.Fatal("expected expired OTP to be deleted")
	}
}
