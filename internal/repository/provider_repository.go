package repository

import (
	"qflow/internal/domain"

	"gorm.io/gorm"
)

type providerRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) domain.ProviderRepository {
	return &providerRepository{db: db}
}

func (r *providerRepository) CreateProvider(provider *domain.Provider) error {
	return r.db.Create(provider).Error
}

func (r *providerRepository) FindProviders() ([]domain.Provider, error) {
	var providers []domain.Provider
	err := r.db.Preload("Category").Order("id asc").Find(&providers).Error
	return providers, err
}

func (r *providerRepository) FindProviderByID(id uint) (*domain.Provider, error) {
	var provider domain.Provider
	err := r.db.First(&provider, id).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *providerRepository) CreateZone(zone *domain.Zone) error {
	return r.db.Create(zone).Error
}

func (r *providerRepository) FindZonesByProviderID(providerID uint) ([]domain.Zone, error) {
	var zones []domain.Zone
	err := r.db.Where("provider_id = ?", providerID).Order("id asc").Find(&zones).Error
	return zones, err
}

func (r *providerRepository) FindZoneByID(id uint) (*domain.Zone, error) {
	var zone domain.Zone
	err := r.db.First(&zone, id).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *providerRepository) UpdateZone(zone *domain.Zone) error {
	return r.db.Save(zone).Error
}

func (r *providerRepository) CountQueuesByZoneID(zoneID uint) (int, error) {
	var count int64
	err := r.db.Model(&domain.Queue{}).Where("zone_id = ?", zoneID).Count(&count).Error
	return int(count), err
}
