package repository

import (
	"context"
	"testing"

	"qflow/internal/domain"
)

func TestNotificationRepo_CreateAndFind(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepository(db)
	ctx := context.Background()

	n := &domain.Notification{UserID: 1, Message: "hello"}
	err := repo.Create(ctx, n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestNotificationRepo_FindByUserID(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &domain.Notification{UserID: 1, Message: "a"})
	repo.Create(ctx, &domain.Notification{UserID: 1, Message: "b"})
	repo.Create(ctx, &domain.Notification{UserID: 2, Message: "c"})

	results, err := repo.FindByUserID(ctx, 1, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2, got %d", len(results))
	}
}

func TestNotificationRepo_CountByUserID(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &domain.Notification{UserID: 1, Message: "a"})
	repo.Create(ctx, &domain.Notification{UserID: 1, Message: "b"})

	count, err := repo.CountByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func TestNotificationRepo_FindByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepository(db)
	ctx := context.Background()

	n := &domain.Notification{UserID: 1, Message: "test"}
	repo.Create(ctx, n)

	found, err := repo.FindByID(ctx, n.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Message != "test" {
		t.Errorf("expected test, got %s", found.Message)
	}
}

func TestNotificationRepo_FindByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepository(db)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, 9999)
	if err == nil {
		t.Error("expected error for non-existent notification")
	}
}

func TestNotificationRepo_MarkRead(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepository(db)
	ctx := context.Background()

	n := &domain.Notification{UserID: 1, Message: "test"}
	repo.Create(ctx, n)

	err := repo.MarkRead(ctx, n.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, _ := repo.FindByID(ctx, n.ID)
	if !found.IsRead {
		t.Error("expected IsRead=true")
	}
}

func TestNotificationRepo_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepository(db)
	ctx := context.Background()

	n := &domain.Notification{UserID: 1, Message: "test"}
	repo.Create(ctx, n)

	err := repo.Delete(ctx, n.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = repo.FindByID(ctx, n.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}
