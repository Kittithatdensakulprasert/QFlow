package repository

import (
	"context"
	"qflow/internal/domain"

	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) domain.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]domain.Notification, error) {
	var notifications []domain.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&domain.Notification{}).Where("user_id = ?", userID).Count(&total).Error
	return total, err
}

func (r *notificationRepository) FindByID(ctx context.Context, id uint) (*domain.Notification, error) {
	var n domain.Notification
	err := r.db.WithContext(ctx).First(&n, id).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *notificationRepository) MarkRead(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&domain.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *notificationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Notification{}, id).Error
}
