package service

import (
	"context"
	"errors"
	"qflow/internal/domain"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrProviderNameRequired     = errors.New("provider name is required")
	ErrZoneNameRequired         = errors.New("zone name is required")
	ErrProviderCategoryNotFound = errors.New("category not found")
	ErrProviderNotFound         = errors.New("provider not found")
	ErrProviderZoneNotFound     = errors.New("zone not found")
)

type providerService struct {
	repo domain.ProviderRepository
}

func NewProviderService(repo domain.ProviderRepository) domain.ProviderService {
	return &providerService{repo: repo}
}

func (s *providerService) CreateProvider(ctx context.Context, name string, categoryID uint) (*domain.Provider, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrProviderNameRequired
	}

	provider := &domain.Provider{Name: name}
	if categoryID > 0 {
		if _, err := s.repo.FindCategoryByID(ctx, categoryID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrProviderCategoryNotFound
			}
			return nil, err
		}
		provider.CategoryID = &categoryID
	}
	if err := s.repo.CreateProvider(ctx, provider); err != nil {
		return nil, err
	}

	return provider, nil
}

func (s *providerService) GetProviders(ctx context.Context) ([]domain.Provider, error) {
	return s.repo.FindProviders(ctx)
}

func (s *providerService) CreateZone(ctx context.Context, providerID uint, name string) (*domain.Zone, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrZoneNameRequired
	}

	if _, err := s.repo.FindProviderByID(ctx, providerID); err != nil {
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
	if err := s.repo.CreateZone(ctx, zone); err != nil {
		return nil, err
	}

	count, err := s.repo.CountQueuesByZoneID(ctx, zone.ID)
	if err != nil {
		return nil, err
	}
	zone.QueueCount = count

	return zone, nil
}

func (s *providerService) GetZones(ctx context.Context, providerID uint) ([]domain.Zone, error) {
	if _, err := s.repo.FindProviderByID(ctx, providerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProviderNotFound
		}
		return nil, err
	}

	zones, err := s.repo.FindZonesByProviderID(ctx, providerID)
	if err != nil {
		return nil, err
	}

	zoneIDs := make([]uint, 0, len(zones))
	for i := range zones {
		zoneIDs = append(zoneIDs, zones[i].ID)
	}

	queueCounts, err := s.repo.CountQueuesByZoneIDs(ctx, zoneIDs)
	if err != nil {
		return nil, err
	}

	for i := range zones {
		zones[i].QueueCount = queueCounts[zones[i].ID]
	}

	return zones, nil
}

func (s *providerService) ToggleZone(ctx context.Context, id uint) (*domain.Zone, error) {
	zone, err := s.repo.FindZoneByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProviderZoneNotFound
		}
		return nil, err
	}

	zone.IsOpen = !zone.IsOpen
	if err := s.repo.UpdateZone(ctx, zone); err != nil {
		return nil, err
	}

	count, err := s.repo.CountQueuesByZoneID(ctx, zone.ID)
	if err != nil {
		return nil, err
	}
	zone.QueueCount = count

	return zone, nil
}
