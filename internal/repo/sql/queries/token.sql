-- name: CreateRejectedToken :exec
INSERT INTO
	token_invalidation (token_id, user_id, expires_at, issued_at)
VALUES
	(?, ?, ?, ?);

-- name: IsTokenRejected :one
SELECT
	COUNT(*) > 0
FROM
	token_invalidation
WHERE
	token_id = ?;

-- name: DeleteExpiredTokens :execrows
DELETE FROM
	token_invalidation
WHERE
	expires_at < ?;

-- name: CreateRegistrationToken :exec
INSERT INTO
	registration_tokens (
		token,
		expires,
		allowed_usage,
		initial_roles,
		created_by,
		created_at
	)
VALUES
	(?, ?, ?, ?, ?, ?);

-- name: ValidateRegistrationToken :one
SELECT
	COUNT(*) > 0
FROM
	registration_tokens
WHERE
	token = ?;

-- name: GetRegistrationToken :one
SELECT
	*
FROM
	registration_tokens
WHERE
	token = ?
	AND allowed_usage > 0
	AND (
		expires IS NULL
		OR expires > ?
	);

-- name: MarkRegistrationTokenUsed :one
UPDATE
	registration_tokens
SET
	allowed_usage = (
		CASE
			WHEN allowed_usage IS NOT NULL THEN allowed_usage - 1
			ELSE NULL
		END
	)
WHERE
	token = ?
	AND (
		expires IS NULL
		OR expires > ?
	)
RETURNING *;

-- name: InsertRecoveryCodes :exec
INSERT INTO
	mfa_backup_codes (code, user_id)
VALUES
	(?, ?);

-- name: CheckAndDeleteRecoveryCode :execrows
DELETE FROM
	mfa_backup_codes
WHERE
	user_id = ?
	AND code = ?;

-- name: RemoveAllRecoveryCodes :exec
DELETE FROM
	mfa_backup_codes
WHERE
	user_id = ?;

-- name: LoadUserRecoveryCodes :many
SELECT
	*
FROM
	mfa_backup_codes
WHERE
	user_id = ?;

-- name: CreateAPIToken :exec
INSERT INTO user_api_tokens (id, token, name, user_id, expires_at) VALUES (?, ?, ?, ?, ?);

-- name: GetAPITokensForUser :many
SELECT * FROM user_api_tokens WHERE user_id = ?;

-- name: RevokeUserAPIToken :execrows
DELETE FROM user_api_tokens WHERE id = ? AND user_id = ?;

-- name: GetUserForAPIToken :one
SELECT 
    sqlc.embed(users),
    sqlc.embed(user_api_tokens)
FROM user_api_tokens
JOIN users ON user_api_tokens.user_id = users.id
WHERE user_api_tokens.token = ?;

-- name: AddRoleToToken :exec
INSERT INTO user_api_token_roles (token_id, role_id) VALUES (?, ?);

-- name: GetRolesForToken :many
SELECT roles.*
FROM user_api_tokens
JOIN user_api_token_roles ON user_api_tokens.id = user_api_token_roles.token_id
JOIN roles ON user_api_token_roles.role_id = roles.id
WHERE user_api_tokens.id = ?;
