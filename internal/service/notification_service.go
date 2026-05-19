package service

import (
	"context"
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

func (s *notificationService) GetNotifications(ctx context.Context, userID uint, page, limit int) ([]domain.Notification, int64, error) {
	offset := (page - 1) * limit
	notifications, err := s.repo.FindByUserID(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (s *notificationService) SendNotification(ctx context.Context, userID uint, message string) (*domain.Notification, error) {
	if message == "" {
		return nil, errors.New("message is required")
	}
	n := &domain.Notification{
		UserID:  userID,
		Message: message,
		IsRead:  false,
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *notificationService) MarkNotificationRead(ctx context.Context, id, userID uint) error {
	n, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return ErrNotificationNotFound
	}
	if n.UserID != userID {
		return ErrNotificationForbidden
	}
	if n.IsRead {
		return nil
	}
	return s.repo.MarkRead(ctx, id)
}

func (s *notificationService) DeleteNotification(ctx context.Context, id, userID uint) error {
	n, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return ErrNotificationNotFound
	}
	if n.UserID != userID {
		return ErrNotificationForbidden
	}
	return s.repo.Delete(ctx, id)
}
