package domain

type Provider struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Zone struct {
	ID         uint   `json:"id"`
	ProviderID uint   `json:"provider_id"`
	Name       string `json:"name"`
	IsOpen     bool   `json:"is_open"`
	QueueCount int    `gorm:"-" json:"queue_count"`
}

type ProviderRepository interface {
	// TODO: define methods
}

type ProviderService interface {
	// TODO: define methods
}
