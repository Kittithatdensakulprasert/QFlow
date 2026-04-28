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

// ===================== Queue Booking =====================

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

	queue := &domain.Queue{
		ZoneID: zoneID,
		UserID: userID,
		Status: "waiting",
	}
	if err := s.repo.CreateWithNextQueueNumber(queue); err != nil {
		return nil, err
	}
	queue.Zone = *zone
	return queue, nil
}

func (s *queueService) GetQueueByNumber(queueNumber int, userID uint) (*domain.Queue, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}

	queue, err := s.repo.FindByQueueNumber(queueNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQueueNotFound
		}
		return nil, err
	}
	if queue.UserID != userID {
		return nil, ErrForbiddenQueue
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
	if queue.Status == "completed" || queue.Status == "skipped" || queue.Status == "called" {
		return ErrQueueFinalized
	}

	return s.repo.UpdateStatus(id, "cancelled")
}

// ===================== Queue Management =====================

func (s *queueService) GetQueuesByZone(zoneID uint) ([]domain.Queue, error) {
	return s.repo.GetByZoneID(zoneID)
}

func (s *queueService) CallQueue(id uint) (*domain.Queue, error) {
	queue, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQueueNotFound
		}
		return nil, err
	}
	if queue.Status != "waiting" {
		return nil, domain.ErrQueueCannotBeCalled
	}

	if err := s.repo.UpdateStatus(id, "called"); err != nil {
		return nil, err
	}

	queue.Status = "called"
	return queue, nil
}

func (s *queueService) CompleteQueue(id uint) (*domain.Queue, error) {
	queue, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQueueNotFound
		}
		return nil, err
	}
	if queue.Status != "called" {
		return nil, domain.ErrQueueCannotBeCompleted
	}

	if err := s.repo.UpdateStatus(id, "completed"); err != nil {
		return nil, err
	}

	queue.Status = "completed"
	return queue, nil
}

func (s *queueService) SkipQueue(id uint) (*domain.Queue, error) {
	queue, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQueueNotFound
		}
		return nil, err
	}
	if queue.Status != "waiting" && queue.Status != "called" {
		return nil, domain.ErrQueueCannotBeSkipped
	}

	if err := s.repo.UpdateStatus(id, "skipped"); err != nil {
		return nil, err
	}

	queue.Status = "skipped"
	return queue, nil
}
