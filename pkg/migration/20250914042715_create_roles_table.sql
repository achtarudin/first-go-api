-- +goose Up
CREATE TABLE roles (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP NULL DEFAULT NULL,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_roles_id (id),
    INDEX idx_roles_name (name)
    
);

INSERT INTO roles (name, created_at, updated_at) VALUES
('admin', NOW(), NOW()),
('merchant', NOW(), NOW()),
('courier', NOW(), NOW()),
('customer', NOW(), NOW());

-- +goose Down
DROP TABLE IF EXISTS roles;
