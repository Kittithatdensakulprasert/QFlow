package service_test

import (
	"errors"
	"qflow/internal/domain"
	"qflow/internal/service"
	"testing"
)

// mock repository
type mockNotificationRepo struct {
	notifications []domain.Notification
	nextID        uint
}

func newMockRepo() *mockNotificationRepo {
	return &mockNotificationRepo{nextID: 1}
}

func (m *mockNotificationRepo) FindByUserID(userID uint) ([]domain.Notification, error) {
	var result []domain.Notification
	for _, n := range m.notifications {
		if n.UserID == userID {
			result = append(result, n)
		}
	}
	return result, nil
}

func (m *mockNotificationRepo) FindByID(id uint) (*domain.Notification, error) {
	for i, n := range m.notifications {
		if n.ID == id {
			return &m.notifications[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockNotificationRepo) Create(n *domain.Notification) error {
	n.ID = m.nextID
	m.nextID++
	m.notifications = append(m.notifications, *n)
	return nil
}

func (m *mockNotificationRepo) MarkRead(id uint) error {
	for i, n := range m.notifications {
		if n.ID == id {
			m.notifications[i].IsRead = true
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockNotificationRepo) Delete(id uint) error {
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

	repo.Create(&domain.Notification{UserID: 1, Message: "hello"})
	repo.Create(&domain.Notification{UserID: 1, Message: "world"})
	repo.Create(&domain.Notification{UserID: 2, Message: "other"})

	result, err := svc.GetNotifications(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 notifications, got %d", len(result))
	}
}

func TestGetNotifications_Empty(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	result, err := svc.GetNotifications(99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 notifications, got %d", len(result))
	}
}

func TestSendNotification(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	n, err := svc.SendNotification(1, "test message")
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

	_, err := svc.SendNotification(1, "")
	if err == nil {
		t.Error("expected error for empty message")
	}
}

func TestMarkNotificationRead(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	n, _ := svc.SendNotification(1, "hello")

	err := svc.MarkNotificationRead(n.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, _ := repo.FindByID(n.ID)
	if !found.IsRead {
		t.Error("expected notification to be marked as read")
	}
}

func TestMarkNotificationRead_AlreadyRead(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	n, _ := svc.SendNotification(1, "hello")
	svc.MarkNotificationRead(n.ID)

	err := svc.MarkNotificationRead(n.ID)
	if err != nil {
		t.Errorf("expected no error for already-read notification, got %v", err)
	}
}

func TestMarkNotificationRead_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	err := svc.MarkNotificationRead(999)
	if err == nil {
		t.Error("expected error for non-existent notification")
	}
}

func TestDeleteNotification(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	n, _ := svc.SendNotification(1, "to delete")

	err := svc.DeleteNotification(n.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	notifications, _ := svc.GetNotifications(1)
	if len(notifications) != 0 {
		t.Error("expected notification to be deleted")
	}
}

func TestDeleteNotification_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewNotificationService(repo)

	err := svc.DeleteNotification(999)
	if err == nil {
		t.Error("expected error for non-existent notification")
	}
}
