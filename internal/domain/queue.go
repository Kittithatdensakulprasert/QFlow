package domain

import (
	"errors"
	"time"
)

// Sentinel errors for queue management
var (
	ErrQueueCannotBeCalled    = errors.New("queue cannot be called")
	ErrQueueCannotBeCompleted = errors.New("only called queue can be completed")
	ErrQueueCannotBeSkipped   = errors.New("cannot skip this queue")
)

type Queue struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	QueueNumber int       `gorm:"uniqueIndex:idx_zone_queue_number;not null" json:"queue_number"`
	ZoneID      uint      `gorm:"uniqueIndex:idx_zone_queue_number;index;not null" json:"zone_id"`
	Zone        Zone      `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status      string    `gorm:"index;default:waiting" json:"status"` // waiting, called, completed, skipped, cancelled
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type QueueRepository interface {
	FindZoneByID(id uint) (*Zone, error)
	CreateWithNextQueueNumber(queue *Queue) error
	FindByQueueNumber(queueNumber int) (*Queue, error)
	FindByID(id uint) (*Queue, error)
	FindByUserID(userID uint) ([]Queue, error)
	UpdateStatus(id uint, status string) error
	GetByZoneID(zoneID uint) ([]Queue, error)
}

type QueueService interface {
	BookQueue(userID, zoneID uint) (*Queue, error)
	GetQueueByNumber(queueNumber int, userID uint) (*Queue, error)
	GetQueueHistory(userID uint) ([]Queue, error)
	CancelQueue(id, userID uint) error
	GetQueuesByZone(zoneID uint) ([]Queue, error)
	CallQueue(id uint) (*Queue, error)
	CompleteQueue(id uint) (*Queue, error)
	SkipQueue(id uint) (*Queue, error)
}
