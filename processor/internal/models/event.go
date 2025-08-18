package models

import "time"

type Event struct {
	ID        int64     `json:"id"`
	EventType string    `json:"event_type"`
	ClientID  string    `json:"client_id"`
	Payload   any       `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}
