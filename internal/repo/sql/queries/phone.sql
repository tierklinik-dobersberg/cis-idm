-- name: GetPhoneNumbersByUserID :many
SELECT
	*
FROM
	user_phone_numbers
WHERE
	user_id = ?;

-- name: GetUserPrimaryPhoneNumber :one
SELECT
	*
FROM
	user_phone_numbers
WHERE
	user_id = ?
	AND is_primary = TRUE;

-- name: GetPhoneNumberByID :one
SELECT
	*
FROM
	user_phone_numbers
WHERE
	user_id = ?
	AND id = ?;

-- name: CreateUserPhoneNumber :one
INSERT INTO
	user_phone_numbers (id, user_id, phone_number, is_primary, verified)
VALUES
	(?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteUserPhoneNumber :execrows
DELETE FROM
	user_phone_numbers
WHERE
	user_id = ?
	AND id = ?;

-- name: MarkPhoneNumberVerified :execrows
UPDATE
	user_phone_numbers
SET
	verified = ?
WHERE
	id = ?
	AND user_id = ?;

-- name: MarkPhoneNumberAsPrimary :execrows
UPDATE
	user_phone_numbers
SET
	is_primary = (id == ?)
WHERE
	user_id = ?;

-- name: MarkPhoneNumberAsVerified :execrows
UPDATE
	user_phone_numbers
SET
	verified = TRUE
WHERE
	user_id = ?
	AND id = ?;