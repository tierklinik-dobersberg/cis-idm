-- name: AssignPermissionToRole :exec
INSERT OR IGNORE INTO
    role_permissions (permission, role_id)
VALUES (?, ?);

-- name: UnassignPermissionFromRole :execrows
DELETE FROM
    role_permissions
WHERE
    role_id = ? AND permission = ?;

-- name: GetRolePermissions :many
SELECT permission FROM role_permissions WHERE role_id = ?;

-- name: GetRolesWithPermission :many
SELECT * FROM role_permissions WHERE permission LIKE ?;

-- name: DeleteAllRolePermissions :exec
DELETE FROM
    role_permissions
WHERE role_id = ?;