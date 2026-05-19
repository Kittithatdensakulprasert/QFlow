package repository

import (
	"context"
	"errors"
	"qflow/internal/domain"

	"gorm.io/gorm"
)

type queueRepository struct {
	db *gorm.DB
}

func NewQueueRepository(db *gorm.DB) domain.QueueRepository {
	return &queueRepository{db: db}
}

func (r *queueRepository) FindZoneByID(ctx context.Context, id uint) (*domain.Zone, error) {
	var zone domain.Zone
	err := r.db.WithContext(ctx).First(&zone, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrQueueZoneRecordNotFound
	}
	return &zone, err
}

func (r *queueRepository) CreateWithNextQueueNumber(ctx context.Context, queue *domain.Queue) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		lockKey := int64(queue.ZoneID) //#nosec G115
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", lockKey).Error; err != nil {
			return err
		}

		var last domain.Queue
		err := tx.Where("zone_id = ?", queue.ZoneID).Order("queue_number desc").First(&last).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound {
			queue.QueueNumber = 1
		} else {
			queue.QueueNumber = last.QueueNumber + 1
		}

		return tx.Create(queue).Error
	})
}

func (r *queueRepository) FindByQueueNumber(ctx context.Context, queueNumber int, userID uint) (*domain.Queue, error) {
	var queue domain.Queue
	err := r.db.WithContext(ctx).
		Preload("Zone").
		Where("queue_number = ? AND user_id = ?", queueNumber, userID).
		Order("created_at desc").
		First(&queue).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrQueueRecordNotFound
	}
	return &queue, err
}

func (r *queueRepository) FindByID(ctx context.Context, id uint) (*domain.Queue, error) {
	var queue domain.Queue
	err := r.db.WithContext(ctx).Preload("Zone").First(&queue, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrQueueRecordNotFound
	}
	return &queue, err
}

func (r *queueRepository) FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]domain.Queue, error) {
	var queues []domain.Queue
	err := r.db.WithContext(ctx).
		Preload("Zone").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&queues).Error
	return queues, err
}

func (r *queueRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&domain.Queue{}).Where("user_id = ?", userID).Count(&total).Error
	return total, err
}

func (r *queueRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Queue{}).Where("id = ?", id).Update("status", status).Error
}

func (r *queueRepository) GetByZoneID(ctx context.Context, zoneID uint) ([]domain.Queue, error) {
	var queues []domain.Queue
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("zone_id = ?", zoneID).
		Order("queue_number asc").
		Find(&queues).Error
	return queues, err
}
