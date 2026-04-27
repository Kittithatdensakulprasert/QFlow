package domain

import "time"

type Queue struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	QueueNumber int       `gorm:"not null" json:"queue_number"`
	ZoneID      uint      `gorm:"not null" json:"zone_id"`
	Zone        Zone      `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status      string    `gorm:"default:waiting" json:"status"` // waiting, called, completed, skipped, cancelled
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type QueueRepository interface {
	// TODO: define methods
}

type QueueService interface {
	// TODO: define methods
}
