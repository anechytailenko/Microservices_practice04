CREATE TABLE IF NOT EXISTS meetups (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    capacity INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL
);