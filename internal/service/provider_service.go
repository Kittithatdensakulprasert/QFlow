package service

import (
	"errors"
	"qflow/internal/domain"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrProviderNameRequired = errors.New("provider name is required")
	ErrZoneNameRequired     = errors.New("zone name is required")
	ErrProviderNotFound     = errors.New("provider not found")
)

type providerService struct {
	repo domain.ProviderRepository
}

func NewProviderService(repo domain.ProviderRepository) domain.ProviderService {
	return &providerService{repo: repo}
}

func (s *providerService) CreateProvider(name string, categoryID uint) (*domain.Provider, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrProviderNameRequired
	}

	provider := &domain.Provider{Name: name}
	if categoryID > 0 {
		provider.CategoryID = &categoryID
	}
	if err := s.repo.CreateProvider(provider); err != nil {
		return nil, err
	}

	return provider, nil
}

func (s *providerService) GetProviders() ([]domain.Provider, error) {
	return s.repo.FindProviders()
}

func (s *providerService) CreateZone(providerID uint, name string) (*domain.Zone, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrZoneNameRequired
	}

	if _, err := s.repo.FindProviderByID(providerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProviderNotFound
		}
		return nil, err
	}

	zone := &domain.Zone{
		ProviderID: providerID,
		Name:       name,
		IsOpen:     true,
	}
	if err := s.repo.CreateZone(zone); err != nil {
		return nil, err
	}

	count, err := s.repo.CountQueuesByZoneID(zone.ID)
	if err != nil {
		return nil, err
	}
	zone.QueueCount = count

	return zone, nil
}

func (s *providerService) GetZones(providerID uint) ([]domain.Zone, error) {
	if _, err := s.repo.FindProviderByID(providerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProviderNotFound
		}
		return nil, err
	}

	zones, err := s.repo.FindZonesByProviderID(providerID)
	if err != nil {
		return nil, err
	}

	for i := range zones {
		count, err := s.repo.CountQueuesByZoneID(zones[i].ID)
		if err != nil {
			return nil, err
		}
		zones[i].QueueCount = count
	}

	return zones, nil
}

func (s *providerService) ToggleZone(id uint) (*domain.Zone, error) {
	zone, err := s.repo.FindZoneByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}

	zone.IsOpen = !zone.IsOpen
	if err := s.repo.UpdateZone(zone); err != nil {
		return nil, err
	}

	count, err := s.repo.CountQueuesByZoneID(zone.ID)
	if err != nil {
		return nil, err
	}
	zone.QueueCount = count

	return zone, nil
}
