package domain

import (
	"context"
	"time"
)

type Provider struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	CategoryID *uint     `gorm:"index" json:"category_id,omitempty"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Zones      []Zone    `gorm:"foreignKey:ProviderID" json:"zones,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Zone struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProviderID uint      `gorm:"index;not null" json:"provider_id"`
	Name       string    `gorm:"not null" json:"name"`
	IsOpen     bool      `gorm:"default:true" json:"is_open"`
	QueueCount int       `gorm:"-" json:"queue_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProviderRepository interface {
	CreateProvider(ctx context.Context, provider *Provider) error
	FindProviders(ctx context.Context) ([]Provider, error)
	FindCategoryByID(ctx context.Context, id uint) (*Category, error)
	FindProviderByID(ctx context.Context, id uint) (*Provider, error)
	CreateZone(ctx context.Context, zone *Zone) error
	FindZonesByProviderID(ctx context.Context, providerID uint) ([]Zone, error)
	FindZoneByID(ctx context.Context, id uint) (*Zone, error)
	UpdateZone(ctx context.Context, zone *Zone) error
	CountQueuesByZoneID(ctx context.Context, zoneID uint) (int, error)
	CountQueuesByZoneIDs(ctx context.Context, zoneIDs []uint) (map[uint]int, error)
}

type ProviderService interface {
	CreateProvider(ctx context.Context, name string, categoryID uint) (*Provider, error)
	GetProviders(ctx context.Context) ([]Provider, error)
	CreateZone(ctx context.Context, providerID uint, name string) (*Zone, error)
	GetZones(ctx context.Context, providerID uint) ([]Zone, error)
	ToggleZone(ctx context.Context, id uint) (*Zone, error)
}
