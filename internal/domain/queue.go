package domain

type Queue struct {
	ID          uint   `json:"id"`
	QueueNumber int    `json:"queue_number"`
	ZoneID      uint   `json:"zone_id"`
	UserID      uint   `json:"user_id"`
	Status      string `json:"status"` // waiting, called, completed, skipped, cancelled
}

type QueueRepository interface {
	// TODO: define methods
}

type QueueService interface {
	// TODO: define methods
}
