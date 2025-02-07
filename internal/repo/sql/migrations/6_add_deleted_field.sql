-- +migrate Up
ALTER TABLE users ADD deleted BOOLEAN NOT NULL;

-- +migrate Down
ALTER TABLE users DROP deleted;