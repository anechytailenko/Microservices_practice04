CREATE TABLE IF NOT EXISTS inbox_events (
<<<<<<< Updated upstream
    event_id VARCHAR(36) PRIMARY KEY,
=======
    event_id VARCHAR(255) PRIMARY KEY,
>>>>>>> Stashed changes
    processed_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);