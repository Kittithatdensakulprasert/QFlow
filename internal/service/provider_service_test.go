package service_test

import (
	"errors"
	"qflow/internal/domain"
	"qflow/internal/service"
	"testing"

	"gorm.io/gorm"
)

type mockProviderRepo struct {
	providers map[uint]domain.Provider
	zones     map[uint]domain.Zone
	queues    []domain.Queue
	nextPID   uint
	nextZID   uint
	repoErr   error
}

func newMockProviderRepo() *mockProviderRepo {
	return &mockProviderRepo{
		providers: map[uint]domain.Provider{},
		zones:     map[uint]domain.Zone{},
		nextPID:   1,
		nextZID:   1,
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
	count := 0
	for _, queue := range m.queues {
		if queue.ZoneID == zoneID {
			count++
		}
	}
	return count, nil
}

func TestProviderServiceCreateProvider(t *testing.T) {
	repo := newMockProviderRepo()
	svc := service.NewProviderService(repo)

	provider, err := svc.CreateProvider(" Bangkok Clinic ", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider.ID == 0 || provider.Name != "Bangkok Clinic" || provider.CategoryID == nil || *provider.CategoryID != 1 {
		t.Fatalf("unexpected provider: %+v", provider)
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
