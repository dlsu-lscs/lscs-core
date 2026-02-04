-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    id VARCHAR(20) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT
);
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE member_roles (
    member_id INT NOT NULL,
    role_id VARCHAR(20) NOT NULL,
    granted_by INT,
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (member_id, role_id),
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES members(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (id, name, description) VALUES
    ('ADMIN', 'Administrator', 'Full system access, can manage all members and settings'),
    ('MODERATOR', 'Moderator', 'Can moderate content and manage basic member issues');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS member_roles;
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
