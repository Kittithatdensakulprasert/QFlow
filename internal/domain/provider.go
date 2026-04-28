package domain

import "time"

type Provider struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	CategoryID *uint     `json:"category_id,omitempty"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Zones      []Zone    `gorm:"foreignKey:ProviderID" json:"zones,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Zone struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProviderID uint      `gorm:"not null" json:"provider_id"`
	Name       string    `gorm:"not null" json:"name"`
	IsOpen     bool      `gorm:"default:true" json:"is_open"`
	QueueCount int       `gorm:"-" json:"queue_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProviderRepository interface {
	CreateProvider(provider *Provider) error
	FindProviders() ([]Provider, error)
	FindProviderByID(id uint) (*Provider, error)
	CreateZone(zone *Zone) error
	FindZonesByProviderID(providerID uint) ([]Zone, error)
	FindZoneByID(id uint) (*Zone, error)
	UpdateZone(zone *Zone) error
	CountQueuesByZoneID(zoneID uint) (int, error)
}

type ProviderService interface {
	CreateProvider(name string, categoryID uint) (*Provider, error)
	GetProviders() ([]Provider, error)
	CreateZone(providerID uint, name string) (*Zone, error)
	GetZones(providerID uint) ([]Zone, error)
	ToggleZone(id uint) (*Zone, error)
}
