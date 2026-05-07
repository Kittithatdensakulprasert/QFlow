package repository

import (
	"qflow/internal/domain"

	"gorm.io/gorm"
)

type queueRepository struct {
	db *gorm.DB
}

func NewQueueRepository(db *gorm.DB) domain.QueueRepository {
	return &queueRepository{db: db}
}

func (r *queueRepository) FindZoneByID(id uint) (*domain.Zone, error) {
	var zone domain.Zone
	if err := r.db.First(&zone, id).Error; err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *queueRepository) CreateWithNextQueueNumber(queue *domain.Queue) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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

func (r *queueRepository) FindByQueueNumber(queueNumber int) (*domain.Queue, error) {
	var queue domain.Queue
	err := r.db.
		Preload("Zone").
		Where("queue_number = ?", queueNumber).
		Order("created_at desc").
		First(&queue).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *queueRepository) FindByID(id uint) (*domain.Queue, error) {
	var queue domain.Queue
	err := r.db.Preload("Zone").First(&queue, id).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *queueRepository) FindByUserID(userID uint) ([]domain.Queue, error) {
	var queues []domain.Queue
	err := r.db.
		Preload("Zone").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&queues).Error
	return queues, err
}

func (r *queueRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.Queue{}).Where("id = ?", id).Update("status", status).Error
}

func (r *queueRepository) GetByZoneID(zoneID uint) ([]domain.Queue, error) {
	var queues []domain.Queue
	err := r.db.
		Preload("User").
		Where("zone_id = ?", zoneID).
		Order("queue_number asc").
		Find(&queues).Error
	return queues, err
}
