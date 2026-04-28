package service

import (
	"errors"
	"qflow/internal/domain"
)

var (
	ErrNotificationNotFound  = errors.New("notification not found")
	ErrNotificationForbidden = errors.New("notification does not belong to user")
)

type notificationService struct {
	repo domain.NotificationRepository
}

func NewNotificationService(repo domain.NotificationRepository) domain.NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) GetNotifications(userID uint) ([]domain.Notification, error) {
	return s.repo.FindByUserID(userID)
}

func (s *notificationService) SendNotification(userID uint, message string) (*domain.Notification, error) {
	if message == "" {
		return nil, errors.New("message is required")
	}
	n := &domain.Notification{
		UserID:  userID,
		Message: message,
		IsRead:  false,
	}
	if err := s.repo.Create(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *notificationService) MarkNotificationRead(id, userID uint) error {
	n, err := s.repo.FindByID(id)
	if err != nil {
		return ErrNotificationNotFound
	}
	if n.UserID != userID {
		return ErrNotificationForbidden
	}
	if n.IsRead {
		return nil
	}
	return s.repo.MarkRead(id)
}

func (s *notificationService) DeleteNotification(id, userID uint) error {
	n, err := s.repo.FindByID(id)
	if err != nil {
		return ErrNotificationNotFound
	}
	if n.UserID != userID {
		return ErrNotificationForbidden
	}
	return s.repo.Delete(id)
}
