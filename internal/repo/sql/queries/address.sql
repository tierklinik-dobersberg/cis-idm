-- name: GetUserAddresses :many
SELECT
    *
FROM
    user_addresses
WHERE
    user_id = ?;

-- name: GetUserAddress :one
SELECT
    *
FROM
    user_addresses
WHERE
    user_id = ?
    AND id = ?;

-- name: CreateUserAddress :one
INSERT INTO
    user_addresses (
        id,
        user_id,
        city_code,
        city_name,
        street,
        extra
    )
VALUES
    (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateUserAddress :one
UPDATE
    user_addresses
SET
    city_code = ?,
    city_name = ?,
    street = ?,
    extra = ?
WHERE
    id = ?
    AND user_id = ?
RETURNING *;

-- name: DeleteUserAddress :execrows
DELETE FROM
    user_addresses
WHERE
    id = ?
    AND user_id = ?;