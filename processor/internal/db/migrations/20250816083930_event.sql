-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    event_type VARCHAR(100) NOT NULL CHECK (char_length(event_type) > 0), 
    client_id VARCHAR(100) NOT NULL CHECK (char_length(client_id) > 0), 
    payload VARCHAR(1000) NOT NULL CHECK (char_length(payload) > 0), 
    "timestamp" TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_events_client_id ON events (client_id);
CREATE INDEX idx_events_event_type ON events (event_type);
-- +goose StatementEnd

