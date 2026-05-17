package service_test

import (
	"context"
	"errors"
	"qflow/internal/domain"
	"qflow/internal/repository"
	"qflow/internal/service"
	"testing"
)

type mockProviderRepo struct {
	providers         map[uint]domain.Provider
	categories        map[uint]domain.Category
	zones             map[uint]domain.Zone
	queues            []domain.Queue
	nextPID           uint
	nextZID           uint
	repoErr           error
	findCategoryErr   error
	createProviderErr error
	findProvidersErr  error
	findProviderErr   error
	createZoneErr     error
	findZonesErr      error
	findZoneErr       error
	updateZoneErr     error
	countByZoneErr    error
	countByZonesErr   error
	countCalls        int
	batchCalls        int
}

func newMockProviderRepo() *mockProviderRepo {
	return &mockProviderRepo{
		providers:  map[uint]domain.Provider{},
		categories: map[uint]domain.Category{},
		zones:      map[uint]domain.Zone{},
		nextPID:    1,
		nextZID:    1,
	}
}

func (m *mockProviderRepo) CreateProvider(_ context.Context, provider *domain.Provider) error {
	if m.createProviderErr != nil {
		return m.createProviderErr
	}
	if m.repoErr != nil {
		return m.repoErr
	}
	provider.ID = m.nextPID
	m.nextPID++
	m.providers[provider.ID] = *provider
	return nil
}

func (m *mockProviderRepo) FindProviders(_ context.Context) ([]domain.Provider, error) {
	if m.findProvidersErr != nil {
		return nil, m.findProvidersErr
	}
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	result := []domain.Provider{}
	for _, provider := range m.providers {
		result = append(result, provider)
	}
	return result, nil
}

func (m *mockProviderRepo) FindCategoryByID(_ context.Context, id uint) (*domain.Category, error) {
	if m.findCategoryErr != nil {
		return nil, m.findCategoryErr
	}
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	category, ok := m.categories[id]
	if !ok {
		return nil, repository.ErrProviderCategoryRecordNotFound
	}
	return &category, nil
}

func (m *mockProviderRepo) FindProviderByID(_ context.Context, id uint) (*domain.Provider, error) {
	if m.findProviderErr != nil {
		return nil, m.findProviderErr
	}
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	provider, ok := m.providers[id]
	if !ok {
		return nil, repository.ErrProviderRecordNotFound
	}
	return &provider, nil
}

func (m *mockProviderRepo) CreateZone(_ context.Context, zone *domain.Zone) error {
	if m.createZoneErr != nil {
		return m.createZoneErr
	}
	if m.repoErr != nil {
		return m.repoErr
	}
	zone.ID = m.nextZID
	m.nextZID++
	m.zones[zone.ID] = *zone
	return nil
}

func (m *mockProviderRepo) FindZonesByProviderID(_ context.Context, providerID uint) ([]domain.Zone, error) {
	if m.findZonesErr != nil {
		return nil, m.findZonesErr
	}
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	result := []domain.Zone{}
	for _, zone := range m.zones {
		if zone.ProviderID == providerID {
			result = append(result, zone)
		}
	}
	return result, nil
}

func (m *mockProviderRepo) FindZoneByID(_ context.Context, id uint) (*domain.Zone, error) {
	if m.findZoneErr != nil {
		return nil, m.findZoneErr
	}
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	zone, ok := m.zones[id]
	if !ok {
		return nil, repository.ErrProviderZoneRecordNotFound
	}
	return &zone, nil
}

func (m *mockProviderRepo) UpdateZone(_ context.Context, zone *domain.Zone) error {
	if m.updateZoneErr != nil {
		return m.updateZoneErr
	}
	if m.repoErr != nil {
		return m.repoErr
	}
	m.zones[zone.ID] = *zone
	return nil
}

func (m *mockProviderRepo) CountQueuesByZoneID(_ context.Context, zoneID uint) (int, error) {
	if m.countByZoneErr != nil {
		return 0, m.countByZoneErr
	}
	if m.repoErr != nil {
		return 0, m.repoErr
	}
	m.countCalls++
	count := 0
	for _, queue := range m.queues {
		if queue.ZoneID == zoneID {
			count++
		}
	}
	return count, nil
}

func (m *mockProviderRepo) CountQueuesByZoneIDs(_ context.Context, zoneIDs []uint) (map[uint]int, error) {
	if m.countByZonesErr != nil {
		return nil, m.countByZonesErr
	}
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	m.batchCalls++
	counts := make(map[uint]int, len(zoneIDs))
	for _, zoneID := range zoneIDs {
		counts[zoneID] = 0
	}
	for _, queue := range m.queues {
		if _, ok := counts[queue.ZoneID]; ok {
			counts[queue.ZoneID]++
		}
	}
	return counts, nil
}

func TestProviderServiceCreateProvider(t *testing.T) {
	repo := newMockProviderRepo()
	repo.categories[1] = domain.Category{ID: 1, Name: "Clinic"}
	svc := service.NewProviderService(repo)

	provider, err := svc.CreateProvider(context.Background(), " Bangkok Clinic ", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider.ID == 0 || provider.Name != "Bangkok Clinic" || provider.CategoryID == nil || *provider.CategoryID != 1 {
		t.Fatalf("unexpected provider: %+v", provider)
	}
}

func TestProviderServiceCreateProviderRequiresExistingCategory(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)

	_, err := svc.CreateProvider(context.Background(), "Bangkok Clinic", 99)
	if !errors.Is(err, service.ErrProviderCategoryNotFound) {
		t.Fatalf("expected category not found error, got %v", err)
	}
}

func TestProviderServiceCreateProviderRequiresName(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)

	_, err := svc.CreateProvider(context.Background(), " ", 0)
	if !errors.Is(err, service.ErrProviderNameRequired) {
		t.Fatalf("expected provider name required error, got %v", err)
	}
}

func TestProviderServiceCreateZoneRequiresExistingProvider(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)

	_, err := svc.CreateZone(context.Background(), 99, "Counter A")
	if !errors.Is(err, service.ErrProviderNotFound) {
		t.Fatalf("expected provider not found error, got %v", err)
	}
}

func TestProviderServiceCreateZoneRequiresName(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)

	_, err := svc.CreateZone(context.Background(), 1, " ")
	if !errors.Is(err, service.ErrZoneNameRequired) {
		t.Fatalf("expected zone name required error, got %v", err)
	}
}

func TestProviderServiceGetZonesCountsQueues(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)
	provider, _ := svc.CreateProvider(context.Background(), "Bangkok Clinic", 0)
	zone, _ := svc.CreateZone(context.Background(), provider.ID, "Counter A")
	repo.queues = append(repo.queues,
		domain.Queue{ID: 1, ZoneID: zone.ID},
		domain.Queue{ID: 2, ZoneID: zone.ID},
		domain.Queue{ID: 3, ZoneID: 999},
	)

	zones, err := svc.GetZones(context.Background(), provider.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(zones) != 1 {
		t.Fatalf("expected 1 zone, got %d", len(zones))
	}
	if zones[0].QueueCount != 2 {
		t.Fatalf("expected queue count 2, got %d", zones[0].QueueCount)
	}
	if repo.batchCalls != 1 {
		t.Fatalf("expected batch count lookup once, got %d", repo.batchCalls)
	}
	if repo.countCalls != 1 {
		t.Fatalf("expected single-zone count only from CreateZone, got %d", repo.countCalls)
	}
}

func TestProviderServiceGetZonesCountsQueuesInOneBatch(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)
	provider, _ := svc.CreateProvider(context.Background(), "Bangkok Clinic", 0)
	zoneA, _ := svc.CreateZone(context.Background(), provider.ID, "Counter A")
	zoneB, _ := svc.CreateZone(context.Background(), provider.ID, "Counter B")
	repo.queues = append(repo.queues,
		domain.Queue{ID: 1, ZoneID: zoneA.ID},
		domain.Queue{ID: 2, ZoneID: zoneA.ID},
		domain.Queue{ID: 3, ZoneID: zoneB.ID},
		domain.Queue{ID: 4, ZoneID: 999},
	)

	zones, err := svc.GetZones(context.Background(), provider.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(zones) != 2 {
		t.Fatalf("expected 2 zones, got %d", len(zones))
	}

	countsByZoneID := map[uint]int{}
	for _, zone := range zones {
		countsByZoneID[zone.ID] = zone.QueueCount
	}
	if countsByZoneID[zoneA.ID] != 2 || countsByZoneID[zoneB.ID] != 1 {
		t.Fatalf("unexpected counts: %+v", countsByZoneID)
	}
	if repo.batchCalls != 1 {
		t.Fatalf("expected batch count lookup once, got %d", repo.batchCalls)
	}
}

func TestProviderServiceToggleZone(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)
	provider, _ := svc.CreateProvider(context.Background(), "Bangkok Clinic", 0)
	zone, _ := svc.CreateZone(context.Background(), provider.ID, "Counter A")

	updated, err := svc.ToggleZone(context.Background(), zone.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.IsOpen {
		t.Fatalf("expected zone to be closed: %+v", updated)
	}
}

func TestProviderServiceGetZones_ProviderNotFound(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)

	_, err := svc.GetZones(context.Background(), 999)
	if !errors.Is(err, service.ErrProviderNotFound) {
		t.Fatalf("expected provider not found error, got %v", err)
	}
}

func TestProviderServiceToggleZone_ZoneNotFound(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)

	_, err := svc.ToggleZone(context.Background(), 999)
	if !errors.Is(err, service.ErrProviderZoneNotFound) {
		t.Fatalf("expected zone not found error, got %v", err)
	}
}

func TestProviderServiceCreateProvider_CreateProviderError(t *testing.T) {
	repo := newMockProviderRepo()
	repo.createProviderErr = errors.New("create provider failed")
	svc := service.NewProviderService(repo)

	_, err := svc.CreateProvider(context.Background(), "Bangkok Clinic", 0)
	if !errors.Is(err, repo.createProviderErr) {
		t.Fatalf("expected create provider error, got %v", err)
	}
}

func TestProviderServiceCreateProvider_FindCategoryGenericError(t *testing.T) {
	repo := newMockProviderRepo()
	repo.findCategoryErr = errors.New("db timeout")
	svc := service.NewProviderService(repo)

	_, err := svc.CreateProvider(context.Background(), "Bangkok Clinic", 1)
	if !errors.Is(err, repo.findCategoryErr) {
		t.Fatalf("expected find category error, got %v", err)
	}
}

func TestProviderServiceGetProviders_GenericError(t *testing.T) {
	repo := newMockProviderRepo()
	repo.findProvidersErr = errors.New("list providers failed")
	svc := service.NewProviderService(repo)

	_, err := svc.GetProviders(context.Background())
	if !errors.Is(err, repo.findProvidersErr) {
		t.Fatalf("expected get providers error, got %v", err)
	}
}

func TestProviderServiceCreateZone_FindProviderGenericError(t *testing.T) {
	repo := newMockProviderRepo()
	repo.findProviderErr = errors.New("db down")
	svc := service.NewProviderService(repo)

	_, err := svc.CreateZone(context.Background(), 1, "Counter A")
	if !errors.Is(err, repo.findProviderErr) {
		t.Fatalf("expected find provider error, got %v", err)
	}
}

func TestProviderServiceCreateZone_CreateZoneError(t *testing.T) {
	repo := newMockProviderRepo()
	repo.providers[1] = domain.Provider{ID: 1, Name: "Clinic"}
	repo.createZoneErr = errors.New("insert zone failed")
	svc := service.NewProviderService(repo)

	_, err := svc.CreateZone(context.Background(), 1, "Counter A")
	if !errors.Is(err, repo.createZoneErr) {
		t.Fatalf("expected create zone error, got %v", err)
	}
}

func TestProviderServiceCreateZone_CountQueuesError(t *testing.T) {
	repo := newMockProviderRepo()
	repo.providers[1] = domain.Provider{ID: 1, Name: "Clinic"}
	repo.countByZoneErr = errors.New("count failed")
	svc := service.NewProviderService(repo)

	_, err := svc.CreateZone(context.Background(), 1, "Counter A")
	if !errors.Is(err, repo.countByZoneErr) {
		t.Fatalf("expected count queues error, got %v", err)
	}
}

func TestProviderServiceGetZones_FindProviderGenericError(t *testing.T) {
	repo := newMockProviderRepo()
	repo.findProviderErr = errors.New("provider lookup failed")
	svc := service.NewProviderService(repo)

	_, err := svc.GetZones(context.Background(), 1)
	if !errors.Is(err, repo.findProviderErr) {
		t.Fatalf("expected find provider error, got %v", err)
	}
}

func TestProviderServiceGetZones_FindZonesError(t *testing.T) {
	repo := newMockProviderRepo()
	repo.providers[1] = domain.Provider{ID: 1, Name: "Clinic"}
	repo.findZonesErr = errors.New("zones query failed")
	svc := service.NewProviderService(repo)

	_, err := svc.GetZones(context.Background(), 1)
	if !errors.Is(err, repo.findZonesErr) {
		t.Fatalf("expected find zones error, got %v", err)
	}
}

func TestProviderServiceGetZones_CountQueuesBatchError(t *testing.T) {
	repo := newMockProviderRepo()
	repo.providers[1] = domain.Provider{ID: 1, Name: "Clinic"}
	repo.zones[1] = domain.Zone{ID: 1, ProviderID: 1, Name: "Counter A", IsOpen: true}
	repo.countByZonesErr = errors.New("batch count failed")
	svc := service.NewProviderService(repo)

	_, err := svc.GetZones(context.Background(), 1)
	if !errors.Is(err, repo.countByZonesErr) {
		t.Fatalf("expected batch count error, got %v", err)
	}
}

func TestProviderServiceToggleZone_FindZoneGenericError(t *testing.T) {
	repo := newMockProviderRepo()
	repo.findZoneErr = errors.New("zone lookup failed")
	svc := service.NewProviderService(repo)

	_, err := svc.ToggleZone(context.Background(), 1)
	if !errors.Is(err, repo.findZoneErr) {
		t.Fatalf("expected find zone error, got %v", err)
	}
}

func TestProviderServiceToggleZone_UpdateError(t *testing.T) {
	repo := newMockProviderRepo()
	repo.zones[1] = domain.Zone{ID: 1, ProviderID: 1, Name: "Counter A", IsOpen: true}
	repo.updateZoneErr = errors.New("update failed")
	svc := service.NewProviderService(repo)

	_, err := svc.ToggleZone(context.Background(), 1)
	if !errors.Is(err, repo.updateZoneErr) {
		t.Fatalf("expected update zone error, got %v", err)
	}
}

func TestProviderServiceToggleZone_CountQueuesError(t *testing.T) {
	repo := newMockProviderRepo()
	repo.zones[1] = domain.Zone{ID: 1, ProviderID: 1, Name: "Counter A", IsOpen: true}
	repo.countByZoneErr = errors.New("count failed")
	svc := service.NewProviderService(repo)

	_, err := svc.ToggleZone(context.Background(), 1)
	if !errors.Is(err, repo.countByZoneErr) {
		t.Fatalf("expected count queues error, got %v", err)
	}
}
