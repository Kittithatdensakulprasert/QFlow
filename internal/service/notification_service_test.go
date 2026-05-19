package service_test

import (
	"context"
	"errors"
	"qflow/internal/domain"
	"qflow/internal/service"
	"testing"
)

// mock repository
type mockNotificationRepo struct {
	notifications []domain.Notification
	nextID        uint
	createErr     error
}

func newMockRepo() *mockNotificationRepo {
	return &mockNotificationRepo{nextID: 1}
}

func (m *mockNotificationRepo) FindByUserID(_ context.Context, userID uint, offset, limit int) ([]domain.Notification, error) {
	var result []domain.Notification
	for _, n := range m.notifications {
		if n.UserID == userID {
			result = append(result, n)
		}
	}
	if offset >= len(result) {
		return []domain.Notification{}, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], nil
}

func (m *mockNotificationRepo) CountByUserID(_ context.Context, userID uint) (int64, error) {
	var count int64
	for _, n := range m.notifications {
		if n.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (m *mockNotificationRepo) FindByID(_ context.Context, id uint) (*domain.Notification, error) {
	for i, n := range m.notifications {
		if n.ID == id {
			return &m.notifications[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockNotificationRepo) Create(_ context.Context, n *domain.Notification) error {
	if m.createErr != nil {
		return m.createErr
	}
	n.ID = m.nextID
	m.nextID++
	m.notifications = append(m.notifications, *n)
	return nil
}

func (m *mockNotificationRepo) MarkRead(_ context.Context, id uint) error {
	for i, n := range m.notifications {
		if n.ID == id {
			m.notifications[i].IsRead = true
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockNotificationRepo) Delete(_ context.Context, id uint) error {
	for i, n := range m.notifications {
		if n.ID == id {
			m.notifications = append(m.notifications[:i], m.notifications[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func TestGetNotifications(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	repo.Create(context.Background(), &domain.Notification{UserID: 1, Message: "hello"})
	repo.Create(context.Background(), &domain.Notification{UserID: 1, Message: "world"})
	repo.Create(context.Background(), &domain.Notification{UserID: 2, Message: "other"})

	result, total, err := svc.GetNotifications(context.Background(), 1, 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 notifications, got %d", len(result))
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

func TestGetNotifications_Empty(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	result, total, err := svc.GetNotifications(context.Background(), 99, 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 notifications, got %d", len(result))
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
}

func TestSendNotification(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	n, err := svc.SendNotification(context.Background(), 1, "test message")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.ID == 0 {
		t.Error("expected notification to have an ID")
	}
	if n.Message != "test message" {
		t.Errorf("expected message 'test message', got '%s'", n.Message)
	}
	if n.IsRead {
		t.Error("new notification should not be read")
	}
}

func TestSendNotification_EmptyMessage(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	_, err := svc.SendNotification(context.Background(), 1, "")
	if err == nil {
		t.Error("expected error for empty message")
	}
}

func TestSendNotification_CreateError(t *testing.T) {
	repo := newMockRepo()
	repo.createErr = errors.New("insert failed")
	svc := service.NewNotificationService(repo)

	n, err := svc.SendNotification(context.Background(), 1, "hello")
	if !errors.Is(err, repo.createErr) {
		t.Fatalf("expected create error, got %v", err)
	}
	if n != nil {
		t.Fatalf("expected nil notification when create fails, got %+v", n)
	}
}

func TestMarkNotificationRead(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	n, _ := svc.SendNotification(context.Background(), 1, "hello")

	err := svc.MarkNotificationRead(context.Background(), n.ID, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, _ := repo.FindByID(context.Background(), n.ID)
	if !found.IsRead {
		t.Error("expected notification to be marked as read")
	}
}

func TestMarkNotificationRead_AlreadyRead(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	n, _ := svc.SendNotification(context.Background(), 1, "hello")
	svc.MarkNotificationRead(context.Background(), n.ID, 1)

	err := svc.MarkNotificationRead(context.Background(), n.ID, 1)
	if err != nil {
		t.Errorf("expected no error for already-read notification, got %v", err)
	}
}

func TestMarkNotificationRead_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	err := svc.MarkNotificationRead(context.Background(), 999, 1)
	if err == nil {
		t.Error("expected error for non-existent notification")
	}
}

func TestMarkNotificationRead_Forbidden(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	n, _ := svc.SendNotification(context.Background(), 1, "hello")
	err := svc.MarkNotificationRead(context.Background(), n.ID, 2)
	if !errors.Is(err, service.ErrNotificationForbidden) {
		t.Errorf("expected forbidden error, got %v", err)
	}
}

func TestDeleteNotification(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	n, _ := svc.SendNotification(context.Background(), 1, "to delete")

	err := svc.DeleteNotification(context.Background(), n.ID, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	notifications, _, _ := svc.GetNotifications(context.Background(), 1, 1, 20)
	if len(notifications) != 0 {
		t.Error("expected notification to be deleted")
	}
}

func TestDeleteNotification_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	err := svc.DeleteNotification(context.Background(), 999, 1)
	if err == nil {
		t.Error("expected error for non-existent notification")
	}
}

func TestDeleteNotification_Forbidden(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	n, _ := svc.SendNotification(context.Background(), 1, "to delete")
	err := svc.DeleteNotification(context.Background(), n.ID, 2)
	if !errors.Is(err, service.ErrNotificationForbidden) {
		t.Errorf("expected forbidden error, got %v", err)
	}
}
