package service_test

import (
	"errors"
	"qflow/internal/domain"
	"qflow/internal/service"
	"testing"

	"gorm.io/gorm"
)

type mockProviderRepo struct {
	providers  map[uint]domain.Provider
	categories map[uint]domain.Category
	zones      map[uint]domain.Zone
	queues     []domain.Queue
	nextPID    uint
	nextZID    uint
	repoErr    error
	countCalls int
	batchCalls int
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

func (m *mockProviderRepo) CreateProvider(provider *domain.Provider) error {
	if m.repoErr != nil {
		return m.repoErr
	}
	provider.ID = m.nextPID
	m.nextPID++
	m.providers[provider.ID] = *provider
	return nil
}

func (m *mockProviderRepo) FindProviders() ([]domain.Provider, error) {
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	result := []domain.Provider{}
	for _, provider := range m.providers {
		result = append(result, provider)
	}
	return result, nil
}

func (m *mockProviderRepo) FindCategoryByID(id uint) (*domain.Category, error) {
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	category, ok := m.categories[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return &category, nil
}

func (m *mockProviderRepo) FindProviderByID(id uint) (*domain.Provider, error) {
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	provider, ok := m.providers[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return &provider, nil
}

func (m *mockProviderRepo) CreateZone(zone *domain.Zone) error {
	if m.repoErr != nil {
		return m.repoErr
	}
	zone.ID = m.nextZID
	m.nextZID++
	m.zones[zone.ID] = *zone
	return nil
}

func (m *mockProviderRepo) FindZonesByProviderID(providerID uint) ([]domain.Zone, error) {
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

func (m *mockProviderRepo) FindZoneByID(id uint) (*domain.Zone, error) {
	if m.repoErr != nil {
		return nil, m.repoErr
	}
	zone, ok := m.zones[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return &zone, nil
}

func (m *mockProviderRepo) UpdateZone(zone *domain.Zone) error {
	if m.repoErr != nil {
		return m.repoErr
	}
	m.zones[zone.ID] = *zone
	return nil
}

func (m *mockProviderRepo) CountQueuesByZoneID(zoneID uint) (int, error) {
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

func (m *mockProviderRepo) CountQueuesByZoneIDs(zoneIDs []uint) (map[uint]int, error) {
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

	provider, err := svc.CreateProvider(" Bangkok Clinic ", 1)
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

	_, err := svc.CreateProvider("Bangkok Clinic", 99)
	if !errors.Is(err, service.ErrProviderCategoryNotFound) {
		t.Fatalf("expected category not found error, got %v", err)
	}
}

func TestProviderServiceCreateProviderRequiresName(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)

	_, err := svc.CreateProvider(" ", 0)
	if !errors.Is(err, service.ErrProviderNameRequired) {
		t.Fatalf("expected provider name required error, got %v", err)
	}
}

func TestProviderServiceCreateZoneRequiresExistingProvider(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)

	_, err := svc.CreateZone(99, "Counter A")
	if !errors.Is(err, service.ErrProviderNotFound) {
		t.Fatalf("expected provider not found error, got %v", err)
	}
}

func TestProviderServiceGetZonesCountsQueues(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)
	provider, _ := svc.CreateProvider("Bangkok Clinic", 0)
	zone, _ := svc.CreateZone(provider.ID, "Counter A")
	repo.queues = append(repo.queues,
		domain.Queue{ID: 1, ZoneID: zone.ID},
		domain.Queue{ID: 2, ZoneID: zone.ID},
		domain.Queue{ID: 3, ZoneID: 999},
	)

	zones, err := svc.GetZones(provider.ID)
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
	provider, _ := svc.CreateProvider("Bangkok Clinic", 0)
	zoneA, _ := svc.CreateZone(provider.ID, "Counter A")
	zoneB, _ := svc.CreateZone(provider.ID, "Counter B")
	repo.queues = append(repo.queues,
		domain.Queue{ID: 1, ZoneID: zoneA.ID},
		domain.Queue{ID: 2, ZoneID: zoneA.ID},
		domain.Queue{ID: 3, ZoneID: zoneB.ID},
		domain.Queue{ID: 4, ZoneID: 999},
	)

	zones, err := svc.GetZones(provider.ID)
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
	provider, _ := svc.CreateProvider("Bangkok Clinic", 0)
	zone, _ := svc.CreateZone(provider.ID, "Counter A")

	updated, err := svc.ToggleZone(zone.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.IsOpen {
		t.Fatalf("expected zone to be closed: %+v", updated)
	}
}
