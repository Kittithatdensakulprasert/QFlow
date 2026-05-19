package repository

import (
	"context"
	"testing"
	"time"

	"qflow/internal/domain"
)

func TestAuthRepo_CreateOTP(t *testing.T) {
	db := newTestDB(t)
	repo := NewAuthRepository(db)
	ctx := context.Background()

	otp, err := repo.CreateOTP(ctx, "0812345678")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if otp.ID == 0 {
		t.Error("expected OTP to have ID")
	}
	if len(otp.Code) != 6 {
		t.Errorf("expected 6-digit code, got %s", otp.Code)
	}
	if otp.Used {
		t.Error("new OTP should not be used")
	}
}

func TestAuthRepo_FindValidOTP(t *testing.T) {
	db := newTestDB(t)
	repo := NewAuthRepository(db)
	ctx := context.Background()

	otp, _ := repo.CreateOTP(ctx, "0812345678")

	found, err := repo.FindValidOTP(ctx, "0812345678", otp.Code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.ID != otp.ID {
		t.Errorf("expected ID %d, got %d", otp.ID, found.ID)
	}
}

func TestAuthRepo_FindValidOTP_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewAuthRepository(db)
	ctx := context.Background()

	_, err := repo.FindValidOTP(ctx, "0812345678", "000000")
	if err == nil {
		t.Error("expected error for non-existent OTP")
	}
}

func TestAuthRepo_MarkOTPAsUsed(t *testing.T) {
	db := newTestDB(t)
	repo := NewAuthRepository(db)
	ctx := context.Background()

	otp, _ := repo.CreateOTP(ctx, "0812345678")

	err := repo.MarkOTPAsUsed(ctx, otp.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Mark again should fail (already used)
	err = repo.MarkOTPAsUsed(ctx, otp.ID)
	if err != domain.ErrOTPAlreadyUsed {
		t.Errorf("expected ErrOTPAlreadyUsed, got %v", err)
	}
}

func TestAuthRepo_DeleteExpiredOTPs(t *testing.T) {
	db := newTestDB(t)
	repo := NewAuthRepository(db)
	ctx := context.Background()

	// Create an expired OTP directly via db
	db.Create(&domain.OTP{
		Phone:     "0812345678",
		Code:      "123456",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		Used:      false,
	})

	n, err := repo.DeleteExpiredOTPs(ctx, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 deleted, got %d", n)
	}
}

func TestAuthRepo_CreateAndFindUser(t *testing.T) {
	db := newTestDB(t)
	repo := NewAuthRepository(db)
	ctx := context.Background()

	user := &domain.User{Phone: "0812345678", Name: "Test", Role: "user"}
	err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID == 0 {
		t.Error("expected user to have ID")
	}

	found, err := repo.FindUserByPhone(ctx, "0812345678")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Name != "Test" {
		t.Errorf("expected name Test, got %s", found.Name)
	}
}

func TestAuthRepo_FindUserByPhone_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewAuthRepository(db)
	ctx := context.Background()

	_, err := repo.FindUserByPhone(ctx, "0000000000")
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestAuthRepo_FindUserByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewAuthRepository(db)
	ctx := context.Background()

	user := &domain.User{Phone: "0812345678", Name: "Test", Role: "user"}
	repo.CreateUser(ctx, user)

	found, err := repo.FindUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.ID != user.ID {
		t.Errorf("expected ID %d, got %d", user.ID, found.ID)
	}
}

func TestAuthRepo_FindUserByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewAuthRepository(db)
	ctx := context.Background()

	_, err := repo.FindUserByID(ctx, 9999)
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestAuthRepo_UpdateUser(t *testing.T) {
	db := newTestDB(t)
	repo := NewAuthRepository(db)
	ctx := context.Background()

	user := &domain.User{Phone: "0812345678", Name: "Old", Role: "user"}
	repo.CreateUser(ctx, user)

	user.Name = "New"
	err := repo.UpdateUser(ctx, user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, _ := repo.FindUserByID(ctx, user.ID)
	if found.Name != "New" {
		t.Errorf("expected name New, got %s", found.Name)
	}
}
