-- name: CreateEMail :one
INSERT INTO
	user_emails (id, user_id, address, verified, is_primary)
VALUES
	(?, ?, ?, ?, ?) RETURNING *;

-- name: GetEmailsForUserByID :many
SELECT
	*
FROM
	user_emails
WHERE
	user_id = ?;

-- name: GetEmailByID :one
SELECT
	*
FROM
	user_emails
WHERE
	user_id = ?
	AND id = ?;

-- name: GetPrimaryEmailForUserByID :one
SELECT
	*
FROM
	user_emails
WHERE
	user_id = ?
	and is_primary = true
LIMIT
	1;

-- name: DeleteEMailFromUser :execrows
DELETE FROM
	user_emails
WHERE
	id = ?
	AND user_id = ?;

-- name: MarkEmailVerified :execrows
UPDATE
	user_emails
SET
	verified = ?
WHERE
	id = ?
	AND user_id = ?;

-- name: MarkEmailAsPrimary :execrows
UPDATE
	user_emails
SET
	is_primary = (id == ?)
WHERE
	user_id = ?;

-- name: MarkEmailAsVerified :execrows
UPDATE
	user_emails
SET
	verified = true
WHERE
	user_id = ?
	AND id = ?;