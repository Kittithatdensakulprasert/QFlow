package service_test

import (
	"errors"
	"qflow/internal/domain"
	"qflow/internal/service"
	"testing"

	"gorm.io/gorm"
)

type mockQueueRepo struct {
	zones   map[uint]domain.Zone
	queues  []domain.Queue
	nextID  uint
	repoErr error
}

func newMockQueueRepo() *mockQueueRepo {
	return &mockQueueRepo{
		zones:  map[uint]domain.Zone{},
		queues: []domain.Queue{},
		nextID: 1,
	}
}

func (m *mockQueueRepo) FindZoneByID(id uint) (*domain.Zone, error) {
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	zone, ok := m.zones[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return &zone, nil
}

func (m *mockQueueRepo) GetNextQueueNumber(zoneID uint) (int, error) {
	if m.repoErr != nil {
		return 0, m.repoErr
	}
	maxQueueNumber := 0
	for _, queue := range m.queues {
		if queue.ZoneID == zoneID && queue.QueueNumber > maxQueueNumber {
			maxQueueNumber = queue.QueueNumber
		}
	}
	return maxQueueNumber + 1, nil
}

func (m *mockQueueRepo) Create(queue *domain.Queue) error {
	if m.repoErr != nil {
		return m.repoErr
	}
	queue.ID = m.nextID
	m.nextID++
	m.queues = append(m.queues, *queue)
	return nil
}

func (m *mockQueueRepo) FindByQueueNumber(queueNumber int) (*domain.Queue, error) {
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	for i := len(m.queues) - 1; i >= 0; i-- {
		if m.queues[i].QueueNumber == queueNumber {
			return &m.queues[i], nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockQueueRepo) FindByID(id uint) (*domain.Queue, error) {
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	for i := range m.queues {
		if m.queues[i].ID == id {
			return &m.queues[i], nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockQueueRepo) FindByUserID(userID uint) ([]domain.Queue, error) {
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	result := []domain.Queue{}
	for _, queue := range m.queues {
		if queue.UserID == userID {
			result = append(result, queue)
		}
	}
	return result, nil
}

func (m *mockQueueRepo) UpdateStatus(id uint, status string) error {
	if m.repoErr != nil {
		return m.repoErr
	}
	for i := range m.queues {
		if m.queues[i].ID == id {
			m.queues[i].Status = status
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func TestBookQueue(t *testing.T) {
	repo := newMockQueueRepo()
	repo.zones[1] = domain.Zone{ID: 1, IsOpen: true}
	svc := service.NewQueueService(repo)

	queue, err := svc.BookQueue(10, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if queue.ID == 0 {
		t.Fatal("expected queue ID to be assigned")
	}
	if queue.QueueNumber != 1 {
		t.Fatalf("expected queue number 1, got %d", queue.QueueNumber)
	}
	if queue.Status != "waiting" {
		t.Fatalf("expected waiting status, got %s", queue.Status)
	}
}

func TestBookQueue_ValidationAndZoneErrors(t *testing.T) {
	repo := newMockQueueRepo()
	repo.zones[1] = domain.Zone{ID: 1, IsOpen: false}
	svc := service.NewQueueService(repo)

	_, err := svc.BookQueue(0, 1)
	if !errors.Is(err, service.ErrInvalidUserID) {
		t.Fatalf("expected invalid user id error, got %v", err)
	}

	_, err = svc.BookQueue(1, 0)
	if !errors.Is(err, service.ErrInvalidZoneID) {
		t.Fatalf("expected invalid zone id error, got %v", err)
	}

	_, err = svc.BookQueue(1, 999)
	if !errors.Is(err, service.ErrZoneNotFound) {
		t.Fatalf("expected zone not found error, got %v", err)
	}

	_, err = svc.BookQueue(1, 1)
	if !errors.Is(err, service.ErrZoneClosed) {
		t.Fatalf("expected zone closed error, got %v", err)
	}
}

func TestBookQueue_PropagatesRepositoryError(t *testing.T) {
	repo := newMockQueueRepo()
	repo.zones[1] = domain.Zone{ID: 1, IsOpen: true}
	repo.repoErr = errors.New("db down")
	svc := service.NewQueueService(repo)

	_, err := svc.BookQueue(1, 1)
	if err == nil || err.Error() != "db down" {
		t.Fatalf("expected db down error, got %v", err)
	}
}

func TestGetQueueByNumber(t *testing.T) {
	repo := newMockQueueRepo()
	repo.zones[1] = domain.Zone{ID: 1, IsOpen: true}
	svc := service.NewQueueService(repo)
	_, _ = svc.BookQueue(1, 1)

	queue, err := svc.GetQueueByNumber(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if queue.QueueNumber != 1 {
		t.Fatalf("expected queue number 1, got %d", queue.QueueNumber)
	}
}

func TestGetQueueByNumber_NotFound(t *testing.T) {
	repo := newMockQueueRepo()
	svc := service.NewQueueService(repo)

	_, err := svc.GetQueueByNumber(999)
	if !errors.Is(err, service.ErrQueueNotFound) {
		t.Fatalf("expected queue not found error, got %v", err)
	}
}

func TestGetQueueHistory(t *testing.T) {
	repo := newMockQueueRepo()
	repo.zones[1] = domain.Zone{ID: 1, IsOpen: true}
	svc := service.NewQueueService(repo)

	_, _ = svc.BookQueue(1, 1)
	_, _ = svc.BookQueue(1, 1)
	_, _ = svc.BookQueue(2, 1)

	queues, err := svc.GetQueueHistory(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(queues) != 2 {
		t.Fatalf("expected 2 queues, got %d", len(queues))
	}
}

func TestGetQueueHistory_InvalidUserID(t *testing.T) {
	repo := newMockQueueRepo()
	svc := service.NewQueueService(repo)

	_, err := svc.GetQueueHistory(0)
	if !errors.Is(err, service.ErrInvalidUserID) {
		t.Fatalf("expected invalid user id error, got %v", err)
	}
}

func TestCancelQueue(t *testing.T) {
	repo := newMockQueueRepo()
	repo.zones[1] = domain.Zone{ID: 1, IsOpen: true}
	svc := service.NewQueueService(repo)
	queue, _ := svc.BookQueue(5, 1)

	err := svc.CancelQueue(queue.ID, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := repo.FindByID(queue.ID)
	if updated.Status != "cancelled" {
		t.Fatalf("expected cancelled status, got %s", updated.Status)
	}
}

func TestCancelQueue_ErrorScenarios(t *testing.T) {
	repo := newMockQueueRepo()
	repo.zones[1] = domain.Zone{ID: 1, IsOpen: true}
	svc := service.NewQueueService(repo)
	queue, _ := svc.BookQueue(5, 1)

	err := svc.CancelQueue(queue.ID, 0)
	if !errors.Is(err, service.ErrInvalidUserID) {
		t.Fatalf("expected invalid user id error, got %v", err)
	}

	err = svc.CancelQueue(999, 5)
	if !errors.Is(err, service.ErrQueueNotFound) {
		t.Fatalf("expected queue not found error, got %v", err)
	}

	err = svc.CancelQueue(queue.ID, 99)
	if !errors.Is(err, service.ErrForbiddenQueue) {
		t.Fatalf("expected forbidden queue error, got %v", err)
	}

	repo.queues[0].Status = "completed"
	err = svc.CancelQueue(queue.ID, 5)
	if !errors.Is(err, service.ErrQueueFinalized) {
		t.Fatalf("expected queue finalized error, got %v", err)
	}

	repo.queues[0].Status = "cancelled"
	err = svc.CancelQueue(queue.ID, 5)
	if !errors.Is(err, service.ErrQueueCancelled) {
		t.Fatalf("expected queue cancelled error, got %v", err)
	}
}
