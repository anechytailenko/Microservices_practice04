CREATE TABLE IF NOT EXISTS inbox_events (
    event_id VARCHAR(255) PRIMARY KEY,
    processed_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS outbox_events (
<<<<<<< Updated upstream
    id VARCHAR(36) PRIMARY KEY,
=======
    id VARCHAR(255) PRIMARY KEY,
>>>>>>> Stashed changes
    event_type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_outbox_events_processed ON outbox_events (processed, created_at);