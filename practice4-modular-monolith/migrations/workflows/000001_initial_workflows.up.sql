CREATE TABLE IF NOT EXISTS workflows (
    workflow_id VARCHAR(36) PRIMARY KEY,
    type VARCHAR(100) NOT NULL,
    state VARCHAR(50) NOT NULL,
    last_error TEXT,
    created_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL,
        updated_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL
);