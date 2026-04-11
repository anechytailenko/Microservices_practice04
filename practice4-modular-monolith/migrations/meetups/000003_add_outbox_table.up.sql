CREATE TABLE IF NOT EXISTS outbox_events (
    id VARCHAR(36) PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL,
        processed_at TIMESTAMP
    WITH
        TIME ZONE
);