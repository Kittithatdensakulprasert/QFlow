package domain

import (
	"context"
	"errors"
	"time"
)

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
	Status      string    `gorm:"index;default:waiting" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type QueueRepository interface {
	FindZoneByID(ctx context.Context, id uint) (*Zone, error)
	CreateWithNextQueueNumber(ctx context.Context, queue *Queue) error
	FindByQueueNumber(ctx context.Context, queueNumber int) (*Queue, error)
	FindByID(ctx context.Context, id uint) (*Queue, error)
	FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]Queue, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	GetByZoneID(ctx context.Context, zoneID uint) ([]Queue, error)
}

type QueueService interface {
	BookQueue(ctx context.Context, userID, zoneID uint) (*Queue, error)
	GetQueueByNumber(ctx context.Context, queueNumber int, userID uint) (*Queue, error)
	GetQueueHistory(ctx context.Context, userID uint, page, limit int) ([]Queue, int64, error)
	CancelQueue(ctx context.Context, id, userID uint) error
	GetQueuesByZone(ctx context.Context, zoneID uint) ([]Queue, error)
	CallQueue(ctx context.Context, id uint) (*Queue, error)
	CompleteQueue(ctx context.Context, id uint) (*Queue, error)
	SkipQueue(ctx context.Context, id uint) (*Queue, error)
}
