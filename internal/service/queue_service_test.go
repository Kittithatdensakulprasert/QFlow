package service

import (
	"context"
	"errors"
	"testing"

	"qflow/internal/domain"
)

// ===================== MOCK REPOSITORY =====================

type mockQueueRepo struct {
	queues        map[uint]*domain.Queue
	nextID        uint
	zones         map[uint]*domain.Zone
	findZoneErr   error
	createErr     error
	findByNumErr  error
	findByIDErr   error
	findByUserErr error
	updateErr     error
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

func (m *mockQueueRepo) FindZoneByID(_ context.Context, id uint) (*domain.Zone, error) {
	if m.findZoneErr != nil {
		return nil, m.findZoneErr
	}
	z, ok := m.zones[id]
	if !ok {
		return nil, domain.ErrQueueZoneRecordNotFound
	}
	cp := *z
	return &cp, nil
}

func (m *mockQueueRepo) CreateWithNextQueueNumber(_ context.Context, q *domain.Queue) error {
	if m.createErr != nil {
		return m.createErr
	}
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

func (m *mockQueueRepo) FindByQueueNumber(_ context.Context, qn int, userID uint) (*domain.Queue, error) {
	if m.findByNumErr != nil {
		return nil, m.findByNumErr
	}
	for _, q := range m.queues {
		if q.QueueNumber == qn && q.UserID == userID {
			cp := *q
			return &cp, nil
		}
	}
	return nil, domain.ErrQueueRecordNotFound
}

func (m *mockQueueRepo) FindByID(_ context.Context, id uint) (*domain.Queue, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	q, ok := m.queues[id]
	if !ok {
		return nil, domain.ErrQueueRecordNotFound
	}
	cp := *q
	return &cp, nil
}

func (m *mockQueueRepo) FindByUserID(_ context.Context, userID uint, offset, limit int) ([]domain.Queue, error) {
	if m.findByUserErr != nil {
		return nil, m.findByUserErr
	}
	var result []domain.Queue
	for _, q := range m.queues {
		if q.UserID == userID {
			result = append(result, *q)
		}
	}
	if offset >= len(result) {
		return []domain.Queue{}, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], nil
}

func (m *mockQueueRepo) CountByUserID(_ context.Context, userID uint) (int64, error) {
	var count int64
	for _, q := range m.queues {
		if q.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (m *mockQueueRepo) UpdateStatus(_ context.Context, id uint, status string) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	q, ok := m.queues[id]
	if !ok {
		return domain.ErrQueueRecordNotFound
	}
	q.Status = status
	return nil
}

func (m *mockQueueRepo) GetByZoneID(_ context.Context, zoneID uint) ([]domain.Queue, error) {
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

func TestNewQueueService(t *testing.T) {
	svc := NewQueueService(newMockRepo())
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

// ===================== BookQueue =====================

func TestBookQueue_Success(t *testing.T) {
	svc := newService()

	queue, err := svc.BookQueue(context.Background(), 77, 20)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if queue.Status != "waiting" {
		t.Fatalf("expected 'waiting', got: %s", queue.Status)
	}
}

func TestBookQueue_ZoneNotFound(t *testing.T) {
	svc := newService()

	_, err := svc.BookQueue(context.Background(), 77, 999)

	if !errors.Is(err, ErrZoneNotFound) {
		t.Fatalf("expected ErrZoneNotFound, got: %v", err)
	}
}

func TestBookQueue_InvalidUserID(t *testing.T) {
	svc := newService()

	_, err := svc.BookQueue(context.Background(), 0, 10)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got: %v", err)
	}
}

func TestBookQueue_InvalidZoneID(t *testing.T) {
	svc := newService()

	_, err := svc.BookQueue(context.Background(), 77, 0)
	if !errors.Is(err, ErrInvalidZoneID) {
		t.Fatalf("expected ErrInvalidZoneID, got: %v", err)
	}
}

func TestBookQueue_ZoneClosed(t *testing.T) {
	svc := newService()
	svc.repo.(*mockQueueRepo).zones[20].IsOpen = false

	_, err := svc.BookQueue(context.Background(), 77, 20)
	if !errors.Is(err, ErrZoneClosed) {
		t.Fatalf("expected ErrZoneClosed, got: %v", err)
	}
}

func TestBookQueue_FindZoneGenericError(t *testing.T) {
	svc := newService()
	repo := svc.repo.(*mockQueueRepo)
	repo.findZoneErr = errors.New("db timeout")

	_, err := svc.BookQueue(context.Background(), 77, 20)
	if !errors.Is(err, repo.findZoneErr) {
		t.Fatalf("expected generic zone lookup error, got: %v", err)
	}
}

func TestBookQueue_CreateQueueError(t *testing.T) {
	svc := newService()
	repo := svc.repo.(*mockQueueRepo)
	repo.createErr = errors.New("insert failed")

	_, err := svc.BookQueue(context.Background(), 77, 20)
	if !errors.Is(err, repo.createErr) {
		t.Fatalf("expected create error, got: %v", err)
	}
}

// ===================== GetQueueByNumber =====================

func TestGetQueueByNumber_Success(t *testing.T) {
	svc := newService()

	queue, err := svc.GetQueueByNumber(context.Background(), 1, 99)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if queue.ID != 1 {
		t.Fatalf("expected queue ID 1, got: %d", queue.ID)
	}
}

func TestGetQueueByNumber_NotFound(t *testing.T) {
	svc := newService()

	_, err := svc.GetQueueByNumber(context.Background(), 999, 99)

	if !errors.Is(err, ErrQueueNotFound) {
		t.Fatalf("expected ErrQueueNotFound, got: %v", err)
	}
}

func TestGetQueueByNumber_InvalidUserID(t *testing.T) {
	svc := newService()

	_, err := svc.GetQueueByNumber(context.Background(), 1, 0)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got: %v", err)
	}
}

func TestGetQueueByNumber_GenericError(t *testing.T) {
	svc := newService()
	repo := svc.repo.(*mockQueueRepo)
	repo.findByNumErr = errors.New("query failed")

	_, err := svc.GetQueueByNumber(context.Background(), 1, 99)
	if !errors.Is(err, repo.findByNumErr) {
		t.Fatalf("expected generic query error, got: %v", err)
	}
}

func TestGetQueueByNumber_NotOwner(t *testing.T) {
	svc := newService()

	_, err := svc.GetQueueByNumber(context.Background(), 5, 99) // queue 5 เป็นของ userID=88

	if !errors.Is(err, ErrQueueNotFound) {
		t.Fatalf("expected ErrQueueNotFound, got: %v", err)
	}
}

// ===================== CancelQueue =====================

func TestCancelQueue_Success(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(context.Background(), 1, 99)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
}

func TestCancelQueue_NotOwner(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(context.Background(), 5, 99) // queue 5 เป็นของ userID=88

	if !errors.Is(err, ErrForbiddenQueue) {
		t.Fatalf("expected ErrForbiddenQueue, got: %v", err)
	}
}

func TestCancelQueue_InvalidState_Completed(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(context.Background(), 3, 99) // status=completed

	if !errors.Is(err, ErrQueueFinalized) {
		t.Fatalf("expected ErrQueueFinalized, got: %v", err)
	}
}

func TestCancelQueue_InvalidState_Called(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(context.Background(), 2, 99) // status=called

	if !errors.Is(err, ErrQueueFinalized) {
		t.Fatalf("expected ErrQueueFinalized, got: %v", err)
	}
}

func TestCancelQueue_AlreadyCancelled(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(context.Background(), 6, 99) // status=cancelled

	if !errors.Is(err, ErrQueueCancelled) {
		t.Fatalf("expected ErrQueueCancelled, got: %v", err)
	}
}

func TestCancelQueue_NotFound(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(context.Background(), 999, 99)

	if !errors.Is(err, ErrQueueNotFound) {
		t.Fatalf("expected ErrQueueNotFound, got: %v", err)
	}
}

func TestCancelQueue_InvalidUserID(t *testing.T) {
	svc := newService()

	err := svc.CancelQueue(context.Background(), 1, 0)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got: %v", err)
	}
}

func TestCancelQueue_GenericFindError(t *testing.T) {
	svc := newService()
	repo := svc.repo.(*mockQueueRepo)
	repo.findByIDErr = errors.New("query failed")

	err := svc.CancelQueue(context.Background(), 1, 99)
	if !errors.Is(err, repo.findByIDErr) {
		t.Fatalf("expected generic find error, got: %v", err)
	}
}

// ===================== CallQueue =====================

func TestCallQueue_Success(t *testing.T) {
	svc := newService()

	queue, err := svc.CallQueue(context.Background(), 1)

	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if queue.Status != "called" {
		t.Fatalf("expected 'called', got: %s", queue.Status)
	}
}

func TestCallQueue_AlreadyCalled(t *testing.T) {
	svc := newService()

	_, err := svc.CallQueue(context.Background(), 2) // status=called

	if !errors.Is(err, domain.ErrQueueCannotBeCalled) {
		t.Fatalf("expected ErrQueueCannotBeCalled, got: %v", err)
	}
}

func TestCallQueue_Completed(t *testing.T) {
	svc := newService()

	_, err := svc.CallQueue(context.Background(), 3) // status=completed

	if !errors.Is(err, domain.ErrQueueCannotBeCalled) {
		t.Fatalf("expected ErrQueueCannotBeCalled, got: %v", err)
	}
}

func TestCallQueue_NotFound(t *testing.T) {
	svc := newService()

	_, err := svc.CallQueue(context.Background(), 999)

	if !errors.Is(err, ErrQueueNotFound) {
		t.Fatalf("expected ErrQueueNotFound, got: %v", err)
	}
}

func TestCallQueue_FindByIDGenericError(t *testing.T) {
	svc := newService()
	repo := svc.repo.(*mockQueueRepo)
	repo.findByIDErr = errors.New("lookup failed")

	_, err := svc.CallQueue(context.Background(), 1)
	if !errors.Is(err, repo.findByIDErr) {
		t.Fatalf("expected generic find error, got: %v", err)
	}
}

func TestCallQueue_UpdateStatusGenericError(t *testing.T) {
	svc := newService()
	repo := svc.repo.(*mockQueueRepo)
	repo.updateErr = errors.New("update failed")

	_, err := svc.CallQueue(context.Background(), 1)
	if !errors.Is(err, repo.updateErr) {
		t.Fatalf("expected generic update error, got: %v", err)
	}
}

// ===================== CompleteQueue =====================

func TestCompleteQueue_Success(t *testing.T) {
	svc := newService()

	queue, err := svc.CompleteQueue(context.Background(), 2) // status=called

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if queue.Status != "completed" {
		t.Fatalf("expected 'completed', got: %s", queue.Status)
	}
}

func TestCompleteQueue_WaitingQueue(t *testing.T) {
	svc := newService()

	_, err := svc.CompleteQueue(context.Background(), 1) // status=waiting

	if !errors.Is(err, domain.ErrQueueCannotBeCompleted) {
		t.Fatalf("expected ErrQueueCannotBeCompleted, got: %v", err)
	}
}

func TestCompleteQueue_SkippedQueue(t *testing.T) {
	svc := newService()

	_, err := svc.CompleteQueue(context.Background(), 4) // status=skipped

	if !errors.Is(err, domain.ErrQueueCannotBeCompleted) {
		t.Fatalf("expected ErrQueueCannotBeCompleted, got: %v", err)
	}
}

func TestCompleteQueue_NotFound(t *testing.T) {
	svc := newService()

	_, err := svc.CompleteQueue(context.Background(), 999)

	if !errors.Is(err, ErrQueueNotFound) {
		t.Fatalf("expected ErrQueueNotFound, got: %v", err)
	}
}

func TestCompleteQueue_FindByIDGenericError(t *testing.T) {
	svc := newService()
	repo := svc.repo.(*mockQueueRepo)
	repo.findByIDErr = errors.New("lookup failed")

	_, err := svc.CompleteQueue(context.Background(), 2)
	if !errors.Is(err, repo.findByIDErr) {
		t.Fatalf("expected generic find error, got: %v", err)
	}
}

func TestCompleteQueue_UpdateStatusGenericError(t *testing.T) {
	svc := newService()
	repo := svc.repo.(*mockQueueRepo)
	repo.updateErr = errors.New("update failed")

	_, err := svc.CompleteQueue(context.Background(), 2)
	if !errors.Is(err, repo.updateErr) {
		t.Fatalf("expected generic update error, got: %v", err)
	}
}

// ===================== SkipQueue =====================

func TestSkipQueue_WaitingSuccess(t *testing.T) {
	svc := newService()

	queue, err := svc.SkipQueue(context.Background(), 1)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if queue.Status != "skipped" {
		t.Fatalf("expected 'skipped', got: %s", queue.Status)
	}
}

func TestSkipQueue_CalledSuccess(t *testing.T) {
	svc := newService()

	queue, err := svc.SkipQueue(context.Background(), 2)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if queue.Status != "skipped" {
		t.Fatalf("expected 'skipped', got: %s", queue.Status)
	}
}

func TestSkipQueue_CompletedFail(t *testing.T) {
	svc := newService()

	_, err := svc.SkipQueue(context.Background(), 3) // status=completed

	if !errors.Is(err, domain.ErrQueueCannotBeSkipped) {
		t.Fatalf("expected ErrQueueCannotBeSkipped, got: %v", err)
	}
}

func TestSkipQueue_CancelledFail(t *testing.T) {
	svc := newService()

	_, err := svc.SkipQueue(context.Background(), 6) // status=cancelled

	if !errors.Is(err, domain.ErrQueueCannotBeSkipped) {
		t.Fatalf("expected ErrQueueCannotBeSkipped, got: %v", err)
	}
}

func TestSkipQueue_NotFound(t *testing.T) {
	svc := newService()

	_, err := svc.SkipQueue(context.Background(), 999)

	if !errors.Is(err, ErrQueueNotFound) {
		t.Fatalf("expected ErrQueueNotFound, got: %v", err)
	}
}

func TestSkipQueue_FindByIDGenericError(t *testing.T) {
	svc := newService()
	repo := svc.repo.(*mockQueueRepo)
	repo.findByIDErr = errors.New("lookup failed")

	_, err := svc.SkipQueue(context.Background(), 1)
	if !errors.Is(err, repo.findByIDErr) {
		t.Fatalf("expected generic find error, got: %v", err)
	}
}

func TestSkipQueue_UpdateStatusGenericError(t *testing.T) {
	svc := newService()
	repo := svc.repo.(*mockQueueRepo)
	repo.updateErr = errors.New("update failed")

	_, err := svc.SkipQueue(context.Background(), 1)
	if !errors.Is(err, repo.updateErr) {
		t.Fatalf("expected generic update error, got: %v", err)
	}
}

// ===================== GetQueuesByZone =====================

func TestGetQueuesByZone_WithQueues(t *testing.T) {
	svc := newService()

	queues, err := svc.GetQueuesByZone(context.Background(), 10)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if len(queues) == 0 {
		t.Fatal("expected queues, got empty slice")
	}
}

func TestGetQueuesByZone_EmptyZone(t *testing.T) {
	svc := newService()

	queues, err := svc.GetQueuesByZone(context.Background(), 999)

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

	queues, total, err := svc.GetQueueHistory(context.Background(), 99, 1, 20)

	if err != nil {
		t.Fatal("expected success, got:", err)
	}
	if len(queues) == 0 {
		t.Fatal("expected history, got empty")
	}
	if total == 0 {
		t.Fatal("expected total > 0")
	}
}

func TestGetQueueHistory_NoHistory(t *testing.T) {
	svc := newService()

	queues, total, err := svc.GetQueueHistory(context.Background(), 777, 1, 20)

	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if len(queues) != 0 {
		t.Fatalf("expected empty history, got %d", len(queues))
	}
	if total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
}

func TestGetQueueHistory_InvalidUserID(t *testing.T) {
	svc := newService()

	_, _, err := svc.GetQueueHistory(context.Background(), 0, 1, 20)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got: %v", err)
	}
}
