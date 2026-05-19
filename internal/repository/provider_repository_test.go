package repository

import (
	"context"
	"testing"

	"qflow/internal/domain"
)

func seedCategory(t *testing.T, db interface {
	Create(ctx context.Context, c *domain.Category) error
}) *domain.Category {
	cat := &domain.Category{Name: "TestCat"}
	if err := db.Create(context.Background(), cat); err != nil {
		panic(err)
	}
	return cat
}

func TestProviderRepo_CreateAndFind(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	ctx := context.Background()

	catID := uint(1)
	err := repo.CreateProvider(ctx, &domain.Provider{Name: "ProvA", CategoryID: &catID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	providers, err := repo.FindProviders(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(providers) != 1 {
		t.Errorf("expected 1, got %d", len(providers))
	}
}

func TestProviderRepo_FindProviders_Empty(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	ctx := context.Background()

	providers, err := repo.FindProviders(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(providers) != 0 {
		t.Errorf("expected 0, got %d", len(providers))
	}
}

func TestProviderRepo_FindCategoryByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	catRepo := NewCategoryGormRepository(db)
	ctx := context.Background()

	cat := &domain.Category{Name: "Health"}
	catRepo.Create(ctx, cat)

	found, err := repo.FindCategoryByID(ctx, cat.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Name != "Health" {
		t.Errorf("expected Health, got %s", found.Name)
	}
}

func TestProviderRepo_FindCategoryByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	ctx := context.Background()

	_, err := repo.FindCategoryByID(ctx, 9999)
	if err != domain.ErrProviderCategoryRecordNotFound {
		t.Errorf("expected ErrProviderCategoryRecordNotFound, got %v", err)
	}
}

func TestProviderRepo_FindProviderByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	ctx := context.Background()

	p := &domain.Provider{Name: "ProvB"}
	repo.CreateProvider(ctx, p)

	found, err := repo.FindProviderByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Name != "ProvB" {
		t.Errorf("expected ProvB, got %s", found.Name)
	}
}

func TestProviderRepo_FindProviderByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	ctx := context.Background()

	_, err := repo.FindProviderByID(ctx, 9999)
	if err != domain.ErrProviderRecordNotFound {
		t.Errorf("expected ErrProviderRecordNotFound, got %v", err)
	}
}

func TestProviderRepo_CreateAndFindZone(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	ctx := context.Background()

	p := &domain.Provider{Name: "ProvC"}
	repo.CreateProvider(ctx, p)

	zone := &domain.Zone{ProviderID: p.ID, Name: "Zone A", IsOpen: true}
	err := repo.CreateZone(ctx, zone)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	zones, err := repo.FindZonesByProviderID(ctx, p.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(zones) != 1 {
		t.Errorf("expected 1, got %d", len(zones))
	}
}

func TestProviderRepo_FindZoneByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	ctx := context.Background()

	p := &domain.Provider{Name: "ProvD"}
	repo.CreateProvider(ctx, p)
	zone := &domain.Zone{ProviderID: p.ID, Name: "Zone B"}
	repo.CreateZone(ctx, zone)

	found, err := repo.FindZoneByID(ctx, zone.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Name != "Zone B" {
		t.Errorf("expected Zone B, got %s", found.Name)
	}
}

func TestProviderRepo_FindZoneByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	ctx := context.Background()

	_, err := repo.FindZoneByID(ctx, 9999)
	if err != domain.ErrProviderZoneRecordNotFound {
		t.Errorf("expected ErrProviderZoneRecordNotFound, got %v", err)
	}
}

func TestProviderRepo_UpdateZone(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	ctx := context.Background()

	p := &domain.Provider{Name: "ProvE"}
	repo.CreateProvider(ctx, p)
	zone := &domain.Zone{ProviderID: p.ID, Name: "Zone C", IsOpen: true}
	repo.CreateZone(ctx, zone)

	zone.IsOpen = false
	err := repo.UpdateZone(ctx, zone)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, _ := repo.FindZoneByID(ctx, zone.ID)
	if found.IsOpen {
		t.Error("expected IsOpen=false")
	}
}

func TestProviderRepo_CountQueuesByZoneID(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	ctx := context.Background()

	p := &domain.Provider{Name: "ProvF"}
	repo.CreateProvider(ctx, p)
	zone := &domain.Zone{ProviderID: p.ID, Name: "Zone D"}
	repo.CreateZone(ctx, zone)

	// Add queues directly
	db.Create(&domain.Queue{ZoneID: zone.ID, UserID: 1, QueueNumber: 1, Status: "waiting"})
	db.Create(&domain.Queue{ZoneID: zone.ID, UserID: 2, QueueNumber: 2, Status: "waiting"})

	count, err := repo.CountQueuesByZoneID(ctx, zone.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func TestProviderRepo_CountQueuesByZoneIDs(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	ctx := context.Background()

	p := &domain.Provider{Name: "ProvG"}
	repo.CreateProvider(ctx, p)
	z1 := &domain.Zone{ProviderID: p.ID, Name: "Z1"}
	z2 := &domain.Zone{ProviderID: p.ID, Name: "Z2"}
	repo.CreateZone(ctx, z1)
	repo.CreateZone(ctx, z2)

	db.Create(&domain.Queue{ZoneID: z1.ID, UserID: 1, QueueNumber: 1, Status: "waiting"})
	db.Create(&domain.Queue{ZoneID: z2.ID, UserID: 2, QueueNumber: 1, Status: "waiting"})
	db.Create(&domain.Queue{ZoneID: z2.ID, UserID: 3, QueueNumber: 2, Status: "waiting"})

	counts, err := repo.CountQueuesByZoneIDs(ctx, []uint{z1.ID, z2.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if counts[z1.ID] != 1 {
		t.Errorf("expected 1 for z1, got %d", counts[z1.ID])
	}
	if counts[z2.ID] != 2 {
		t.Errorf("expected 2 for z2, got %d", counts[z2.ID])
	}
}

func TestProviderRepo_CountQueuesByZoneIDs_Empty(t *testing.T) {
	db := newTestDB(t)
	repo := NewProviderRepository(db)
	ctx := context.Background()

	counts, err := repo.CountQueuesByZoneIDs(ctx, []uint{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(counts) != 0 {
		t.Errorf("expected empty map, got %v", counts)
	}
}
