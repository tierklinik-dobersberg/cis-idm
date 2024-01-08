// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: permissions.sql

package repo

import (
	"context"
)

const assignPermissionToRole = `-- name: AssignPermissionToRole :exec
INSERT OR IGNORE INTO
    role_permissions (permission, role_id)
VALUES (?, ?)
`

type AssignPermissionToRoleParams struct {
	Permission string
	RoleID     string
}

func (q *Queries) AssignPermissionToRole(ctx context.Context, arg AssignPermissionToRoleParams) error {
	_, err := q.db.ExecContext(ctx, assignPermissionToRole, arg.Permission, arg.RoleID)
	return err
}

const deleteAllRolePermissions = `-- name: DeleteAllRolePermissions :exec
DELETE FROM
    role_permissions
WHERE role_id = ?
`

func (q *Queries) DeleteAllRolePermissions(ctx context.Context, roleID string) error {
	_, err := q.db.ExecContext(ctx, deleteAllRolePermissions, roleID)
	return err
}

const getRolePermissions = `-- name: GetRolePermissions :many
SELECT permission FROM role_permissions WHERE role_id = ?
`

func (q *Queries) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getRolePermissions, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		items = append(items, permission)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRolesWithPermission = `-- name: GetRolesWithPermission :many
SELECT permission, role_id FROM role_permissions WHERE permission LIKE ?
`

func (q *Queries) GetRolesWithPermission(ctx context.Context, permission string) ([]RolePermission, error) {
	rows, err := q.db.QueryContext(ctx, getRolesWithPermission, permission)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RolePermission
	for rows.Next() {
		var i RolePermission
		if err := rows.Scan(&i.Permission, &i.RoleID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unassignPermissionFromRole = `-- name: UnassignPermissionFromRole :execrows
DELETE FROM
    role_permissions
WHERE
    role_id = ? AND permission = ?
`

type UnassignPermissionFromRoleParams struct {
	RoleID     string
	Permission string
}

func (q *Queries) UnassignPermissionFromRole(ctx context.Context, arg UnassignPermissionFromRoleParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, unassignPermissionFromRole, arg.RoleID, arg.Permission)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}