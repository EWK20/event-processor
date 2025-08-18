package processor_test

import (
	"context"
	"time"

	"github.com/EWK20/event-processor/processor/internal/models"
)

type FakeDB struct {
	events []models.Event
}

func NewFakeDB() *FakeDB {
	return &FakeDB{
		events: []models.Event{
			{
				ID:        1,
				EventType: "transaction_approved",
				ClientID:  "client_123",
				Payload:   `{"amount":"176.11","currency":"GBP","transaction_id":"txn_806"}`,
				Timestamp: time.Date(2025, 8, 18, 7, 48, 48, 0, time.UTC),
			},
			{
				ID:        2,
				EventType: "user_signup",
				ClientID:  "client_456",
				Payload:   `{"username": "john_doe"}`,
				Timestamp: time.Date(2025, 8, 18, 7, 49, 0, 0, time.UTC),
			},
			{
				ID:        3,
				EventType: "transaction_approved",
				ClientID:  "client_789",
				Payload:   `{"amount":"428.54","currency":"GBP","transaction_id":"txn_912"}`,
				Timestamp: time.Date(2025, 8, 18, 8, 30, 48, 0, time.UTC),
			},
		},
	}
}

func (db *FakeDB) Save(_ context.Context, event models.Event) error {
	lastID := db.events[len(db.events)-1].ID

	newID := lastID + 1

	event.ID = newID

	db.events = append(db.events, event)

	return nil
}
