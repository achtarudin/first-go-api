-- +goose Up
CREATE TABLE transactions (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED,
    merchant_id INT UNSIGNED,
    courier_id INT UNSIGNED,
    food_id INT UNSIGNED,
    quantity INT UNSIGNED NOT NULL,
    total_price INT UNSIGNED NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP NULL DEFAULT NULL,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (merchant_id) REFERENCES merchants(id),
    FOREIGN KEY (courier_id) REFERENCES couriers(id),
    FOREIGN KEY (food_id) REFERENCES foods(id),
    INDEX idx_transactions_user_id (user_id),
    INDEX idx_transactions_merchant_id (merchant_id),
    INDEX idx_transactions_courier_id (courier_id),
    INDEX idx_transactions_food_id (food_id)
);

-- +goose Down
DROP TABLE IF EXISTS transactions;