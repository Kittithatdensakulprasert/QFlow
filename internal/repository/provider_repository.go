package repository

import (
	"context"
	"errors"
	"qflow/internal/domain"

	"gorm.io/gorm"
)

var (
	ErrProviderCategoryRecordNotFound = errors.New("provider category record not found")
	ErrProviderRecordNotFound         = errors.New("provider record not found")
	ErrProviderZoneRecordNotFound     = errors.New("provider zone record not found")
)

type providerRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) domain.ProviderRepository {
	return &providerRepository{db: db}
}

func (r *providerRepository) CreateProvider(ctx context.Context, provider *domain.Provider) error {
	return r.db.WithContext(ctx).Create(provider).Error
}

func (r *providerRepository) FindProviders(ctx context.Context) ([]domain.Provider, error) {
	var providers []domain.Provider
	err := r.db.WithContext(ctx).Preload("Category").Order("id asc").Find(&providers).Error
	return providers, err
}

func (r *providerRepository) FindCategoryByID(ctx context.Context, id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProviderCategoryRecordNotFound
	}
	return &category, err
}

func (r *providerRepository) FindProviderByID(ctx context.Context, id uint) (*domain.Provider, error) {
	var provider domain.Provider
	err := r.db.WithContext(ctx).First(&provider, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProviderRecordNotFound
	}
	return &provider, err
}

func (r *providerRepository) CreateZone(ctx context.Context, zone *domain.Zone) error {
	return r.db.WithContext(ctx).Create(zone).Error
}

func (r *providerRepository) FindZonesByProviderID(ctx context.Context, providerID uint) ([]domain.Zone, error) {
	var zones []domain.Zone
	err := r.db.WithContext(ctx).Where("provider_id = ?", providerID).Order("id asc").Find(&zones).Error
	return zones, err
}

func (r *providerRepository) FindZoneByID(ctx context.Context, id uint) (*domain.Zone, error) {
	var zone domain.Zone
	err := r.db.WithContext(ctx).First(&zone, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProviderZoneRecordNotFound
	}
	return &zone, err
}

func (r *providerRepository) UpdateZone(ctx context.Context, zone *domain.Zone) error {
	return r.db.WithContext(ctx).Save(zone).Error
}

func (r *providerRepository) CountQueuesByZoneID(ctx context.Context, zoneID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Queue{}).Where("zone_id = ?", zoneID).Count(&count).Error
	return int(count), err
}

func (r *providerRepository) CountQueuesByZoneIDs(ctx context.Context, zoneIDs []uint) (map[uint]int, error) {
	counts := make(map[uint]int, len(zoneIDs))
	if len(zoneIDs) == 0 {
		return counts, nil
	}

	var rows []struct {
		ZoneID uint
		Count  int
	}
	err := r.db.WithContext(ctx).Model(&domain.Queue{}).
		Select("zone_id, COUNT(*) AS count").
		Where("zone_id IN ?", zoneIDs).
		Group("zone_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		counts[row.ZoneID] = row.Count
	}
	return counts, nil
}
