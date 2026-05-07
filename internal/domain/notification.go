package domain

import "time"

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
	FindByUserID(userID uint) ([]Notification, error)
	FindByID(id uint) (*Notification, error)
	Create(n *Notification) error
	MarkRead(id uint) error
	Delete(id uint) error
}

type NotificationService interface {
	GetNotifications(userID uint) ([]Notification, error)
	SendNotification(userID uint, message string) (*Notification, error)
	MarkNotificationRead(id, userID uint) error
	DeleteNotification(id, userID uint) error
}
