package domain

import "time"

type Queue struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	QueueNumber int       `gorm:"uniqueIndex;not null" json:"queue_number"`
	ZoneID      uint      `gorm:"not null" json:"zone_id"`
	Zone        Zone      `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status      string    `gorm:"default:waiting" json:"status"` // waiting, called, completed, skipped, cancelled
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type QueueRepository interface {
	FindZoneByID(id uint) (*Zone, error)
	GetNextQueueNumber(zoneID uint) (int, error)
	Create(queue *Queue) error
	FindByQueueNumber(queueNumber int) (*Queue, error)
	FindByID(id uint) (*Queue, error)
	FindByUserID(userID uint) ([]Queue, error)
	UpdateStatus(id uint, status string) error
}

type QueueService interface {
	BookQueue(userID, zoneID uint) (*Queue, error)
	GetQueueByNumber(queueNumber int, userID uint) (*Queue, error)
	GetQueueHistory(userID uint) ([]Queue, error)
	CancelQueue(id, userID uint) error
}
