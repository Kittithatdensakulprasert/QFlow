package repository

import (
	"context"
	"testing"

	"qflow/internal/domain"
)

func setupQueueTest(t *testing.T) (*domain.Zone, *domain.User, *queueRepository) {
	db := newTestDB(t)
	repo := NewQueueRepository(db).(*queueRepository)

	user := &domain.User{Phone: "0812345678", Name: "Test", Role: "user"}
	db.Create(user)

	p := &domain.Provider{Name: "Prov"}
	db.Create(p)

	zone := &domain.Zone{ProviderID: p.ID, Name: "Zone A", IsOpen: true}
	db.Create(zone)

	return zone, user, repo
}

// insertQueue inserts a queue directly bypassing pg_advisory_xact_lock
func insertQueue(t *testing.T, repo *queueRepository, zoneID, userID uint, num int, status string) *domain.Queue {
	t.Helper()
	q := &domain.Queue{ZoneID: zoneID, UserID: userID, QueueNumber: num, Status: status}
	if err := repo.db.Create(q).Error; err != nil {
		t.Fatalf("failed to insert queue: %v", err)
	}
	return q
}

func TestQueueRepo_FindZoneByID(t *testing.T) {
	zone, _, repo := setupQueueTest(t)
	ctx := context.Background()

	found, err := repo.FindZoneByID(ctx, zone.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.ID != zone.ID {
		t.Errorf("expected ID %d, got %d", zone.ID, found.ID)
	}
}

func TestQueueRepo_FindZoneByID_NotFound(t *testing.T) {
	_, _, repo := setupQueueTest(t)
	ctx := context.Background()

	_, err := repo.FindZoneByID(ctx, 9999)
	if err != domain.ErrQueueZoneRecordNotFound {
		t.Errorf("expected ErrQueueZoneRecordNotFound, got %v", err)
	}
}

func TestQueueRepo_FindByQueueNumber(t *testing.T) {
	zone, user, repo := setupQueueTest(t)
	ctx := context.Background()

	q := insertQueue(t, repo, zone.ID, user.ID, 1, "waiting")

	found, err := repo.FindByQueueNumber(ctx, q.QueueNumber, user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.ID != q.ID {
		t.Errorf("expected ID %d, got %d", q.ID, found.ID)
	}
}

func TestQueueRepo_FindByQueueNumber_NotFound(t *testing.T) {
	_, user, repo := setupQueueTest(t)
	ctx := context.Background()

	_, err := repo.FindByQueueNumber(ctx, 9999, user.ID)
	if err != domain.ErrQueueRecordNotFound {
		t.Errorf("expected ErrQueueRecordNotFound, got %v", err)
	}
}

func TestQueueRepo_FindByID(t *testing.T) {
	zone, user, repo := setupQueueTest(t)
	ctx := context.Background()

	q := insertQueue(t, repo, zone.ID, user.ID, 1, "waiting")

	found, err := repo.FindByID(ctx, q.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.ID != q.ID {
		t.Errorf("expected ID %d, got %d", q.ID, found.ID)
	}
}

func TestQueueRepo_FindByID_NotFound(t *testing.T) {
	_, _, repo := setupQueueTest(t)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, 9999)
	if err != domain.ErrQueueRecordNotFound {
		t.Errorf("expected ErrQueueRecordNotFound, got %v", err)
	}
}

func TestQueueRepo_FindByUserID(t *testing.T) {
	zone, user, repo := setupQueueTest(t)
	ctx := context.Background()

	insertQueue(t, repo, zone.ID, user.ID, 1, "waiting")
	insertQueue(t, repo, zone.ID, user.ID, 2, "waiting")

	queues, err := repo.FindByUserID(ctx, user.ID, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(queues) != 2 {
		t.Errorf("expected 2, got %d", len(queues))
	}
}

func TestQueueRepo_CountByUserID(t *testing.T) {
	zone, user, repo := setupQueueTest(t)
	ctx := context.Background()

	insertQueue(t, repo, zone.ID, user.ID, 1, "waiting")
	insertQueue(t, repo, zone.ID, user.ID, 2, "waiting")

	count, err := repo.CountByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func TestQueueRepo_UpdateStatus(t *testing.T) {
	zone, user, repo := setupQueueTest(t)
	ctx := context.Background()

	q := insertQueue(t, repo, zone.ID, user.ID, 1, "waiting")

	err := repo.UpdateStatus(ctx, q.ID, "called")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, _ := repo.FindByID(ctx, q.ID)
	if found.Status != "called" {
		t.Errorf("expected called, got %s", found.Status)
	}
}

func TestQueueRepo_GetByZoneID(t *testing.T) {
	zone, user, repo := setupQueueTest(t)
	ctx := context.Background()

	insertQueue(t, repo, zone.ID, user.ID, 1, "waiting")
	insertQueue(t, repo, zone.ID, user.ID, 2, "called")

	queues, err := repo.GetByZoneID(ctx, zone.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(queues) != 2 {
		t.Errorf("expected 2, got %d", len(queues))
	}
}
