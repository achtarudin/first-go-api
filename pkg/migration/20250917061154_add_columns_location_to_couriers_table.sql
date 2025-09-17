-- +goose Up
ALTER TABLE couriers ADD COLUMN location POINT AFTER phone;

UPDATE couriers SET location = POINT(IFNULL(longitude, 0), IFNULL(latitude, 0)) WHERE location IS NULL;

ALTER TABLE couriers MODIFY COLUMN location POINT NOT NULL;

ALTER TABLE couriers ADD SPATIAL INDEX(location);

-- +goose Down
ALTER TABLE couriers DROP INDEX location;

ALTER TABLE couriers DROP COLUMN location;
