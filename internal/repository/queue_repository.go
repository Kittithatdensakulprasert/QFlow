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

func (r *queueRepository) GetNextQueueNumber(zoneID uint) (int, error) {
	var last domain.Queue
	err := r.db.Order("queue_number desc").First(&last).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, nil
		}
		return 0, err
	}
	return last.QueueNumber + 1, nil
}

func (r *queueRepository) Create(queue *domain.Queue) error {
	return r.db.Create(queue).Error
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
