package domain

import "time"

type Provider struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	CategoryID uint      `json:"category_id"`
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
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProviderRepository interface {
	// TODO: define methods
}

type ProviderService interface {
	// TODO: define methods
}
