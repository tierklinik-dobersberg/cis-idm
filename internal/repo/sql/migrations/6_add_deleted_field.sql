-- +migrate Up
ALTER TABLE users ADD deleted BOOLEAN NOT NULL DEFAULT false;

-- +migrate Down
ALTER TABLE users DROP deleted;