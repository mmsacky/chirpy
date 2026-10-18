-- +goose UP
ALTER TABLE users
ADD COLUMN token TEXT NOT NULL DEFAULT 'unset';

-- +goose DOWN
ALTER TABLE users
DROP COLUMN token;
