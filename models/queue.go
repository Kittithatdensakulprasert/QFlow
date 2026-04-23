package models

type Queue struct {
	ID          int    `json:"id"`
	QueueNumber int    `json:"queueNumber"`
	Status      string `json:"status"`
}
