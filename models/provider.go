package models

type Provider struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Zone struct {
	ID         int    `json:"id"`
	ProviderID int    `json:"providerId"`
	Name       string `json:"name"`
	IsOpen     bool   `json:"isOpen"`
	QueueCount int    `json:"queueCount"`
}
