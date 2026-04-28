package service

import (
	"errors"
	"qflow/internal/domain"

	"gorm.io/gorm"
)

var (
	ErrInvalidUserID  = errors.New("user id is required")
	ErrInvalidZoneID  = errors.New("zone id is required")
	ErrZoneNotFound   = errors.New("zone not found")
	ErrZoneClosed     = errors.New("zone is closed")
	ErrQueueNotFound  = errors.New("queue not found")
	ErrForbiddenQueue = errors.New("queue does not belong to user")
	ErrQueueFinalized = errors.New("queue cannot be cancelled")
	ErrQueueCancelled = errors.New("queue already cancelled")
)

type queueService struct {
	repo domain.QueueRepository
}

func NewQueueService(repo domain.QueueRepository) domain.QueueService {
	return &queueService{repo: repo}
}

func (s *queueService) BookQueue(userID, zoneID uint) (*domain.Queue, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if zoneID == 0 {
		return nil, ErrInvalidZoneID
	}

	zone, err := s.repo.FindZoneByID(zoneID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}
	if !zone.IsOpen {
		return nil, ErrZoneClosed
	}

	queueNumber, err := s.repo.GetNextQueueNumber(zoneID)
	if err != nil {
		return nil, err
	}

	queue := &domain.Queue{
		QueueNumber: queueNumber,
		ZoneID:      zoneID,
		UserID:      userID,
		Status:      "waiting",
	}
	if err := s.repo.Create(queue); err != nil {
		return nil, err
	}
	queue.Zone = *zone
	return queue, nil
}

func (s *queueService) GetQueueByNumber(queueNumber int) (*domain.Queue, error) {
	queue, err := s.repo.FindByQueueNumber(queueNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQueueNotFound
		}
		return nil, err
	}
	return queue, nil
}

func (s *queueService) GetQueueHistory(userID uint) ([]domain.Queue, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	return s.repo.FindByUserID(userID)
}

func (s *queueService) CancelQueue(id, userID uint) error {
	if userID == 0 {
		return ErrInvalidUserID
	}

	queue, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQueueNotFound
		}
		return err
	}

	if queue.UserID != userID {
		return ErrForbiddenQueue
	}
	if queue.Status == "cancelled" {
		return ErrQueueCancelled
	}
	if queue.Status == "completed" || queue.Status == "skipped" {
		return ErrQueueFinalized
	}

	return s.repo.UpdateStatus(id, "cancelled")
}
