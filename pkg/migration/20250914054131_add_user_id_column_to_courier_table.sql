-- +goose Up
ALTER TABLE couriers DROP COLUMN name;

ALTER TABLE couriers ADD COLUMN user_id INT UNSIGNED AFTER id;

ALTER TABLE couriers ADD CONSTRAINT fk_couriers_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE couriers ADD CONSTRAINT unique_couriers_user_id UNIQUE (user_id);


-- +goose Down
ALTER TABLE couriers DROP FOREIGN KEY fk_couriers_user_id;

ALTER TABLE couriers DROP INDEX unique_couriers_user_id;

ALTER TABLE couriers DROP COLUMN user_id;

ALTER TABLE couriers ADD COLUMN name VARCHAR(100) AFTER id;
