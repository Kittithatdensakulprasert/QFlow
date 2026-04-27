package domain

type Notification struct {
	ID      uint   `json:"id"`
	UserID  uint   `json:"user_id"`
	Message string `json:"message"`
	IsRead  bool   `json:"is_read"`
}

type NotificationRepository interface {
	// TODO: define methods
}

type NotificationService interface {
	// TODO: define methods
}
