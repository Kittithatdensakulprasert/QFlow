package domain

import (
	"context"
	"time"
)

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Message   string    `gorm:"not null" json:"message"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NotificationRepository interface {
	FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]Notification, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
	FindByID(ctx context.Context, id uint) (*Notification, error)
	Create(ctx context.Context, n *Notification) error
	MarkRead(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}

type NotificationService interface {
	GetNotifications(ctx context.Context, userID uint, page, limit int) ([]Notification, int64, error)
	SendNotification(ctx context.Context, userID uint, message string) (*Notification, error)
	MarkNotificationRead(ctx context.Context, id, userID uint) error
	DeleteNotification(ctx context.Context, id, userID uint) error
}
