-- +goose Up
CREATE TABLE sessions (
    id VARCHAR(64) PRIMARY KEY,
    member_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_agent VARCHAR(512),
    ip_address VARCHAR(45),
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE,
    INDEX idx_sessions_member_id (member_id),
    INDEX idx_sessions_expires_at (expires_at)
);

-- +goose Down
DROP TABLE IF EXISTS sessions;
