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

-- name: GetSystemRoles :many
SELECT
	*
FROM
	roles
WHERE origin = 'system';


-- name: CreateRole :one
INSERT INTO
	roles (id, name, description, origin, delete_protected)
VALUES
	(?, ?, ?, 'api', ?)
RETURNING *;

-- name: CreateSystemRole :one
INSERT INTO
	roles (id, name, description, origin, delete_protected)
VALUES
	(?, ?, ?, 'system', true)
ON CONFLICT(id)	DO
	UPDATE
	SET
		id = excluded.id,
		name = excluded.name,
		description = excluded.description,
		origin = 'system',
		delete_protected = true
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
INSERT OR IGNORE INTO
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