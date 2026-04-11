CREATE TABLE IF NOT EXISTS notifications (
    event_id VARCHAR(36) PRIMARY KEY,
    correlation_id VARCHAR(36) NOT NULL,
    owner_user_id VARCHAR(36) NOT NULL,
    payload JSONB NOT NULL,
    received_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_notifications_owner_user_id ON notifications (owner_user_id);