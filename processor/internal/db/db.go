package db

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/lib/pq" // Imports postgres driver

	"github.com/EWK20/event-processor/processor/internal/config"
	"github.com/EWK20/event-processor/processor/internal/models"
	"github.com/pressly/goose/v3"
)

var (
	ErrFailedToConnectToDB    = errors.New("failed to connect to database")
	ErrFailedToPingDB         = errors.New("failed to ping database")
	ErrFailedToSave           = errors.New("failed to save data")
	ErrFailedToMarshalPayload = errors.New("failed to marshal payload data")
)

type Database struct {
	Conn *sql.DB
}

func New(cfg config.DB) (*Database, error) {
	connStr := "user=" + cfg.User + " password=" + cfg.Password + " host=" + cfg.Host + " port=" + cfg.Port + " dbname=" + cfg.DBName + " sslmode=" + cfg.SSLMode

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%w, %w", ErrFailedToConnectToDB, err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("%w, %w", ErrFailedToPingDB, err)
	}

	return &Database{
		conn,
	}, nil
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (db *Database) RunMigrations() error {
	goose.SetBaseFS(embedMigrations)

	migrationsDir := "migrations"

	// Run migrations
	if err := goose.Up(db.Conn, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func (db *Database) Save(ctx context.Context, event models.Event) error {
	query := `
	INSERT INTO events (
		event_type, client_id, payload, "timestamp"
	) VALUES (
	 	$1, $2, $3, $4
	)`

	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToMarshalPayload, err)
	}

	_, err = db.Conn.ExecContext(ctx, query, event.EventType, event.ClientID, string(payloadJSON), event.Timestamp.UTC())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSave, err)
	}

	return nil
}
