package repository

import (
	"qflow/internal/domain"

	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) domain.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) FindByUserID(userID uint) ([]domain.Notification, error) {
	var notifications []domain.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) FindByID(id uint) (*domain.Notification, error) {
	var n domain.Notification
	err := r.db.First(&n, id).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepository) Create(n *domain.Notification) error {
	return r.db.Create(n).Error
}

func (r *notificationRepository) MarkRead(id uint) error {
	return r.db.Model(&domain.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *notificationRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Notification{}, id).Error
}
