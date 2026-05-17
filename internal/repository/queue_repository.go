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

var (
	ErrQueueZoneRecordNotFound = errors.New("queue zone record not found")
	ErrQueueRecordNotFound     = errors.New("queue record not found")
)

func NewQueueRepository(db *gorm.DB) domain.QueueRepository {
	return &queueRepository{db: db}
}

func (r *queueRepository) FindZoneByID(ctx context.Context, id uint) (*domain.Zone, error) {
	var zone domain.Zone
	err := r.db.WithContext(ctx).First(&zone, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrQueueZoneRecordNotFound
	}
	return &zone, err
}

func (r *queueRepository) CreateWithNextQueueNumber(ctx context.Context, queue *domain.Queue) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", int64(1001)).Error; err != nil { //nolint:gosec // G201: hardcoded lock key, not user input
			return err
		}

		var last domain.Queue
		err := tx.Order("queue_number desc").First(&last).Error
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

func (r *queueRepository) FindByQueueNumber(ctx context.Context, queueNumber int) (*domain.Queue, error) {
	var queue domain.Queue
	err := r.db.WithContext(ctx).
		Preload("Zone").
		Where("queue_number = ?", queueNumber).
		Order("created_at desc").
		First(&queue).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrQueueRecordNotFound
	}
	return &queue, err
}

func (r *queueRepository) FindByID(ctx context.Context, id uint) (*domain.Queue, error) {
	var queue domain.Queue
	err := r.db.WithContext(ctx).Preload("Zone").First(&queue, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrQueueRecordNotFound
	}
	return &queue, err
}

func (r *queueRepository) FindByUserID(ctx context.Context, userID uint) ([]domain.Queue, error) {
	var queues []domain.Queue
	err := r.db.WithContext(ctx).
		Preload("Zone").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&queues).Error
	return queues, err
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
