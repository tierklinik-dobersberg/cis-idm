-- +migrate Up
CREATE TABLE IF NOT EXISTS registration_tokens (
    token TEXT NOT NULL PRIMARY KEY,
    expires TIMESTAMP,
    allowed_usage INTEGER,
    initial_roles TEXT NOT NULL DEFAULT '',
    created_by TEXT NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate StatementBegin
CREATE TRIGGER IF NOT EXISTS registration_token_cleanup
AFTER
UPDATE
    ON registration_tokens BEGIN
DELETE FROM
    registration_tokens
WHERE
    allowed_usage IS NOT NULL
    AND allowed_usage = 0;
END;
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE registration_tokens;
DROP TRIGGER registration_token_cleanup;