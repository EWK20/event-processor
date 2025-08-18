package db_test

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	_ "github.com/lib/pq" // Imports postgres driver
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EWK20/event-processor/processor/internal/config"
	"github.com/EWK20/event-processor/processor/internal/db"
	"github.com/EWK20/event-processor/processor/internal/models"
)

type SaveTest struct {
	input models.Event
	err   error
}

func TestSave(t *testing.T) {
	now := time.Now().UTC()

	testCases := map[string]SaveTest{
		"Save Successful": {
			input: models.Event{
				EventType: "transaction_approved",
				ClientID:  "client_123",
				Payload: map[string]any{
					"transaction_id": "txn_456",
					"amount":         120.50,
					"currency":       "GBP",
				},
				Timestamp: now,
			},
			err: nil,
		},
		"Event Type Empty Save Error": {
			input: models.Event{
				EventType: "",
				ClientID:  "client_123",
				Payload: map[string]any{
					"transaction_id": "txn_456",
					"amount":         120.50,
					"currency":       "GBP",
				},
				Timestamp: now,
			},
			err: db.ErrFailedToSave,
		},
		"Client ID Empty Save Error": {
			input: models.Event{
				EventType: "transaction_approved",
				ClientID:  "",
				Payload: map[string]any{
					"transaction_id": "txn_456",
					"amount":         120.50,
					"currency":       "GBP",
				},
				Timestamp: now,
			},
			err: db.ErrFailedToSave,
		},
	}

	for scenario, test := range testCases {
		t.Run(scenario, func(t *testing.T) {
			ctx := t.Context()

			db, teardown := setupDB(t)
			defer teardown()

			err := db.Save(ctx, test.input)

			if test.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, test.err)

				return
			}

			require.NoError(t, err)

			// verify row was inserted
			var (
				eventType string
				clientID  string
				payload   string
			)
			err = db.Conn.QueryRow(`SELECT event_type, client_id, payload FROM events ORDER BY id DESC LIMIT 1`).
				Scan(&eventType, &clientID, &payload)

			payloadJSON, _ := json.Marshal(test.input.Payload)

			require.NoError(t, err)
			assert.Equal(t, test.input.EventType, eventType)
			assert.Equal(t, test.input.ClientID, clientID)
			assert.Contains(t, string(payloadJSON), `"currency":"GBP"`)

		})
	}
}

func setupDB(t *testing.T) (*db.Database, func()) {
	t.Helper()

	ctx := t.Context()

	dbCfg := config.DB{
		User:     "user",
		Password: "password",
		Host:     "localhost",
		Port:     "6432",
		DBName:   "test",
		SSLMode:  "disable",
	}

	connStr := "user=" + dbCfg.User + " password=" + dbCfg.Password + " host=" + dbCfg.Host + " port=" + dbCfg.Port + " dbname=default" + " sslmode=" + dbCfg.SSLMode

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := conn.Ping(); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	// Drop existing database connections
	_, err = conn.ExecContext(ctx, "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '"+dbCfg.DBName+"';")
	if err != nil {
		t.Fatalf("failed to drop existing connections: %v", err)
	}

	// Drop test schema if it already exists
	_, err = conn.ExecContext(ctx, "DROP DATABASE IF EXISTS "+dbCfg.DBName)
	if err != nil {
		t.Fatalf("failed to drop schema: %v", err)
	}

	// Create test database
	_, err = conn.ExecContext(ctx, "CREATE DATABASE "+dbCfg.DBName)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	db, err := db.New(dbCfg)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.RunMigrations(); err != nil {
		t.Fatalf("failed to run database migrations: %v", err)
	}

	teardown := func() {
		// Drop existing database connections
		_, err = conn.ExecContext(ctx, "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '"+dbCfg.DBName+"';")
		if err != nil {
			t.Fatalf("failed to drop existing connections: %v", err)
		}

		// Drop test schema
		_, err = conn.ExecContext(ctx, "DROP DATABASE IF EXISTS "+dbCfg.DBName)
		if err != nil {
			t.Fatalf("failed to drop schema: %v", err)
		}
	}

	return db, teardown
}
