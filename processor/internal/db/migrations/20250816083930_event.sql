-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    event_type VARCHAR(100) NOT NULL CHECK (char_length(event_type) > 0), 
    client_id VARCHAR(100) NOT NULL CHECK (char_length(client_id) > 0), 
    payload VARCHAR(1000) NOT NULL CHECK (char_length(payload) > 0), 
    "timestamp" TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

