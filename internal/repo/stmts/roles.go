package stmts

import "github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"

var (
	GetRoleByName = Statement[models.Role]{
		Query: `SELECT * FROM roles WHERE name = ?`,
		Args:  []string{"name"},
	}

	GetRoleByID = Statement[models.Role]{
		Query: `SELECT * FROM roles WHERE id = ?`,
		Args:  []string{"id"},
	}

	CreateRole = Statement[any]{
		Query: `INSERT INTO roles (id, name, description, delete_protected) VALUES (?, ?, ?, ?)`,
		Args:  []string{"id", "name", "description", "delete_protected"},
	}

	AssignRoleToUser = Statement[any]{
		Query: `INSERT INTO role_assignments (user_id, role_id) VALUES (?, ?)`,
		Args:  []string{"user_id", "role_id"},
	}

	GetRolesForUser = Statement[models.Role]{
		Query: `SELECT
				roles.id as id, roles.name as name, roles.description as description
			FROM role_assignments
			JOIN roles ON roles.id = role_id
			WHERE user_id = ?`,
		Args: []string{"user_id"},
	}

	GetUsersByRole = Statement[models.User]{
		Query: `SELECT * FROM role_assignments
		JOIN users ON users.id = user_id
		WHERE role_id = ?`,
		Args: []string{"role_id"},
	}
)
