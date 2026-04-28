package models

import "time"

type Queue struct {
	ID          int        `json:"id"`
	QueueNumber int        `json:"queueNumber"`
	Status      string     `json:"status"` // waiting, called, completed, skipped
	UserID      uint       `json:"userId"`
	ZoneID      uint       `json:"zoneId"`
	CreatedAt   time.Time  `json:"createdAt"`
	CalledAt    *time.Time `json:"calledAt"`
}
