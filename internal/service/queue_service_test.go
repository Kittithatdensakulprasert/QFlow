package service

import (
	"errors"
	"testing"

	"qflow/internal/domain"

	"gorm.io/gorm"
)

// ===================== MOCK REPOSITORY =====================

type mockQueueRepo struct {
	queues map[uint]*domain.Queue
	nextID uint
	zones  map[uint]*domain.Zone
}

func newMockRepo() *mockQueueRepo {
	return &mockQueueRepo{
		nextID: 7,
		queues: map[uint]*domain.Queue{
			1: {ID: 1, QueueNumber: 1, ZoneID: 10, UserID: 99, Status: "waiting"},
			2: {ID: 2, QueueNumber: 2, ZoneID: 10, UserID: 99, Status: "called"},
			3: {ID: 3, QueueNumber: 3, ZoneID: 10, UserID: 99, Status: "completed"},
			4: {ID: 4, QueueNumber: 4, ZoneID: 10, UserID: 99, Status: "skipped"},
			5: {ID: 5, QueueNumber: 5, ZoneID: 10, UserID: 88, Status: "waiting"}, // คนอื่น
			6: {ID: 6, QueueNumber: 6, ZoneID: 10, UserID: 99, Status: "cancelled"},
		},
		zones: map[uint]*domain.Zone{
			10: {ID: 10, IsOpen: true},
			20: {ID: 20, IsOpen: true},
			99: {ID: 99, IsOpen: true},
		},
	}
}

func (m *mockQueueRepo) FindZoneByID(id uint) (*domain.Zone, error) {
	z, ok := m.zones[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *z
	return &cp, nil
}

func (m *mockQueueRepo) CreateWithNextQueueNumber(q *domain.Queue) error {
	maxQN := 0
	for _, existing := range m.queues {
		if existing.QueueNumber > maxQN {
			maxQN = existing.QueueNumber
		}
	}
	q.QueueNumber = maxQN + 1
	q.ID = m.nextID
	m.nextID++
	cp := *q
	m.queues[q.ID] = &cp
	return nil
}

func (m *mockQueueRepo) FindByQueueNumber(qn int) (*domain.Queue, error) {
	for _, q := range m.queues {
		if q.QueueNumber == qn {
			cp := *q
			return &cp, nil
		}
	}
	return nil, ErrQueueNotFound
}

func (m *mockQueueRepo) FindByID(id uint) (*domain.Queue, error) {
	q, ok := m.queues[id]
	if !ok {
		return nil, ErrQueueNotFound
	}
	cp := *q
	return &cp, nil
}

func (m *mockQueueRepo) FindByUserID(userID uint) ([]domain.Queue, error) {
	var result []domain.Queue
	for _, q := range m.queues {
		if q.UserID == userID {
			result = append(result, *q)
		}
	}
	return result, nil
}

func (m *mockQueueRepo) UpdateStatus(id uint, status string) error {
	q, ok := m.queues[id]
	if !ok {
		return ErrQueueNotFound
	}
	q.Status = status
	return nil
}

func (m *mockQueueRepo) GetByZoneID(zoneID uint) ([]domain.Queue, error) {
	var result []domain.Queue
	for _, q := range m.queues {
		if q.ZoneID == zoneID {
			result = append(result, *q)
		}
	}
	return result, nil
}

// ===================== SERVICE =====================

func newService() *queueService {
	return &queueService{repo: newMockRepo()}
}

// ===================== BookQueue =====================

func TestBookQueue_Success(t *testing.T) {
	svc := newService()

	queue, err := svc.BookQueue(77, 20)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if queue.Status != "waiting" {
		t.Fatalf("expected 'waiting', got: %s", queue.Status)
	}
}

func TestBookQueue_ZoneNotFound(t *testing.T) {
	svc := newService()

	_, err := svc.BookQueue(77, 999)

	if !errors.Is(err, ErrZoneNotFound) {
		t.Fatalf("expected ErrZoneNotFound, got: %v", err)
	}
}

// ===================== GetQueueByNumber =====================

func TestGetQueueByNumber_Success(t *testing.T) {
	svc := newService()

	queue, err := svc.GetQueueByNumber(1, 99)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if queue.ID != 1 {
		t.Fatalf("expected queue ID 1, got: %d", queue.ID)
	}
}

func TestGetQueueByNumber_NotFound(t *testing.T) {
	svc := newService()

	_, err := svc.GetQueueByNumber(999, 99)

	if !errors.Is(err, ErrQueueNotFound) {
		t.Fatalf("expected ErrQueueNotFound, got: %v", err)
	}
}

func TestGetQueueByNumber_Forbidden(t *testing.T) {
	svc := newService()

	_, err := svc.GetQueueByNumber(5, 99) // queue 5 เป็นของ userID=88

	if !errors.Is(err, ErrForbiddenQueue) {
		t.Fatalf("expected ErrForbiddenQueue, got: %v", err)
	}
}

// ===================== CancelQueue =====================

func TestCancelQueue_Success(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(1, 99)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
}

func TestCancelQueue_NotOwner(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(5, 99) // queue 5 เป็นของ userID=88

	if !errors.Is(err, ErrForbiddenQueue) {
		t.Fatalf("expected ErrForbiddenQueue, got: %v", err)
	}
}

func TestCancelQueue_InvalidState_Completed(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(3, 99) // status=completed

	if !errors.Is(err, ErrQueueFinalized) {
		t.Fatalf("expected ErrQueueFinalized, got: %v", err)
	}
}

func TestCancelQueue_InvalidState_Called(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(2, 99) // status=called

	if !errors.Is(err, ErrQueueFinalized) {
		t.Fatalf("expected ErrQueueFinalized, got: %v", err)
	}
}

func TestCancelQueue_AlreadyCancelled(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(6, 99) // status=cancelled

	if !errors.Is(err, ErrQueueCancelled) {
		t.Fatalf("expected ErrQueueCancelled, got: %v", err)
	}
}

func TestCancelQueue_NotFound(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(999, 99)

	if !errors.Is(err, ErrQueueNotFound) {
		t.Fatalf("expected ErrQueueNotFound, got: %v", err)
	}
}

// ===================== CallQueue =====================

func TestCallQueue_Success(t *testing.T) {
	svc := newService()

	queue, err := svc.CallQueue(1)

	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if queue.Status != "called" {
		t.Fatalf("expected 'called', got: %s", queue.Status)
	}
}

func TestCallQueue_AlreadyCalled(t *testing.T) {
	svc := newService()

	_, err := svc.CallQueue(2) // status=called

	if !errors.Is(err, domain.ErrQueueCannotBeCalled) {
		t.Fatalf("expected ErrQueueCannotBeCalled, got: %v", err)
	}
}

func TestCallQueue_Completed(t *testing.T) {
	svc := newService()

	_, err := svc.CallQueue(3) // status=completed

	if !errors.Is(err, domain.ErrQueueCannotBeCalled) {
		t.Fatalf("expected ErrQueueCannotBeCalled, got: %v", err)
	}
}

func TestCallQueue_NotFound(t *testing.T) {
	svc := newService()

	_, err := svc.CallQueue(999)

	if !errors.Is(err, ErrQueueNotFound) {
		t.Fatalf("expected ErrQueueNotFound, got: %v", err)
	}
}

// ===================== CompleteQueue =====================

func TestCompleteQueue_Success(t *testing.T) {
	svc := newService()

	queue, err := svc.CompleteQueue(2) // status=called

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if queue.Status != "completed" {
		t.Fatalf("expected 'completed', got: %s", queue.Status)
	}
}

func TestCompleteQueue_WaitingQueue(t *testing.T) {
	svc := newService()

	_, err := svc.CompleteQueue(1) // status=waiting

	if !errors.Is(err, domain.ErrQueueCannotBeCompleted) {
		t.Fatalf("expected ErrQueueCannotBeCompleted, got: %v", err)
	}
}

func TestCompleteQueue_SkippedQueue(t *testing.T) {
	svc := newService()

	_, err := svc.CompleteQueue(4) // status=skipped

	if !errors.Is(err, domain.ErrQueueCannotBeCompleted) {
		t.Fatalf("expected ErrQueueCannotBeCompleted, got: %v", err)
	}
}

func TestCompleteQueue_NotFound(t *testing.T) {
	svc := newService()

	_, err := svc.CompleteQueue(999)

	if !errors.Is(err, ErrQueueNotFound) {
		t.Fatalf("expected ErrQueueNotFound, got: %v", err)
	}
}

// ===================== SkipQueue =====================

func TestSkipQueue_WaitingSuccess(t *testing.T) {
	svc := newService()

	queue, err := svc.SkipQueue(1)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if queue.Status != "skipped" {
		t.Fatalf("expected 'skipped', got: %s", queue.Status)
	}
}

func TestSkipQueue_CalledSuccess(t *testing.T) {
	svc := newService()

	queue, err := svc.SkipQueue(2)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if queue.Status != "skipped" {
		t.Fatalf("expected 'skipped', got: %s", queue.Status)
	}
}

func TestSkipQueue_CompletedFail(t *testing.T) {
	svc := newService()

	_, err := svc.SkipQueue(3) // status=completed

	if !errors.Is(err, domain.ErrQueueCannotBeSkipped) {
		t.Fatalf("expected ErrQueueCannotBeSkipped, got: %v", err)
	}
}

func TestSkipQueue_CancelledFail(t *testing.T) {
	svc := newService()

	_, err := svc.SkipQueue(6) // status=cancelled

	if !errors.Is(err, domain.ErrQueueCannotBeSkipped) {
		t.Fatalf("expected ErrQueueCannotBeSkipped, got: %v", err)
	}
}

func TestSkipQueue_NotFound(t *testing.T) {
	svc := newService()

	_, err := svc.SkipQueue(999)

	if !errors.Is(err, ErrQueueNotFound) {
		t.Fatalf("expected ErrQueueNotFound, got: %v", err)
	}
}

// ===================== GetQueuesByZone =====================

func TestGetQueuesByZone_WithQueues(t *testing.T) {
	svc := newService()

	queues, err := svc.GetQueuesByZone(10)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if len(queues) == 0 {
		t.Fatal("expected queues, got empty slice")
	}
}

func TestGetQueuesByZone_EmptyZone(t *testing.T) {
	svc := newService()

	queues, err := svc.GetQueuesByZone(999)

	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if len(queues) != 0 {
		t.Fatalf("expected empty slice, got %d", len(queues))
	}
}

// ===================== GetQueueHistory =====================

func TestGetQueueHistory_Success(t *testing.T) {
	svc := newService()

	queues, err := svc.GetQueueHistory(99)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if len(queues) == 0 {
		t.Fatal("expected history, got empty")
	}
}

func TestGetQueueHistory_NoHistory(t *testing.T) {
	svc := newService()

	queues, err := svc.GetQueueHistory(777)

	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if len(queues) != 0 {
		t.Fatalf("expected empty history, got %d", len(queues))
	}
}
