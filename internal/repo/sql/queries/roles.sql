-- name: GetRoleByName :one
SELECT
	*
FROM
	roles
WHERE
	name = ?;

-- name: GetRoleByID :one
SELECT
	*
FROM
	roles
WHERE
	id = ?;

-- name: GetRoles :many
SELECT
	*
FROM
	roles;

-- name: CreateRole :one
INSERT INTO
	roles (id, name, description, delete_protected)
VALUES
	(?, ?, ?, ?)
RETURNING *;

-- name: UpdateRole :one
UPDATE
	roles
SET
	name = ?,
	description = ?,
	delete_protected = ?
WHERE
	id = ?
RETURNING *;

-- name: AssignRoleToUser :exec
INSERT INTO
	role_assignments (user_id, role_id)
VALUES
	(?, ?);

-- name: UnassignRoleFromUser :execrows
DELETE FROM
	role_assignments
WHERE
	user_id = ?
	AND role_id = ?;

-- name: GetRolesForUser :many
SELECT
	roles.*
FROM
	role_assignments
	JOIN roles ON roles.id = role_id
WHERE
	user_id = ?;

-- name: GetUsersByRole :many
SELECT
	*
FROM
	role_assignments
	JOIN users ON users.id = user_id
WHERE
	role_id = ?;

-- name: DeleteRole :execrows
DELETE FROM
	roles
WHERE
	id = ?;