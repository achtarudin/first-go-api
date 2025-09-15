-- +goose Up
CREATE TABLE merchants (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address TEXT,
    user_id INT UNSIGNED,
    created_at TIMESTAMP NULL DEFAULT NULL,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    INDEX idx_merchants_user_id (user_id)
);

CREATE TABLE couriers (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    latitude DOUBLE,
    longitude DOUBLE,
    created_at TIMESTAMP NULL DEFAULT NULL,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL
    -- Index belum ditambah, bisa ditambah jika sering dicari berdasarkan phone misal
);

CREATE TABLE foods (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price INT UNSIGNED NOT NULL,
    merchant_id INT UNSIGNED,
    description TEXT,
    created_at TIMESTAMP NULL DEFAULT NULL,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    CONSTRAINT fk_foods_merchant_id FOREIGN KEY fk_foods_merchant_id (merchant_id) REFERENCES merchants(id),
    INDEX idx_foods_merchant_id (merchant_id)
);

-- +goose Down
DROP TABLE IF EXISTS foods;
DROP TABLE IF EXISTS couriers;
DROP TABLE IF EXISTS merchants;