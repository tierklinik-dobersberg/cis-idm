-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = ?;

-- name: GetUserByName :one
SELECT
    *
FROM
    users
WHERE
    username = ?;

-- name: SetUserExtraData :execrows
UPDATE users
SET extra = ?
WHERE id = ?;


-- name: GetUserByEMail :one
SELECT
    sqlc.embed(users),
    user_emails.verified
FROM
    users
    JOIN user_emails ON user_emails.user_id = users.id
WHERE
    user_emails.address = ?;

-- name: GetAllUsers :many
SELECT
    *
FROM
    users;

-- name: DeleteUser :execrows
DELETE FROM
    users
WHERE
    id = ?;

-- name: CreateUser :one
INSERT INTO
    users (
        id,
        username,
        display_name,
        first_name,
        last_name,
        extra,
        avatar,
        birthday,
        password
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
			username = ?,
			display_name = ?,
			first_name = ?,
			last_name = ?,
			extra = ?,
			avatar = ?,
			birthday = ?
		WHERE id = ?
        RETURNING *;
