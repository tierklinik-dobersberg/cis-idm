package stmts

import "github.com/tierklinik-dobersberg/cis-userd/internal/repo/models"

var (
	GetGroupByName = Statement[models.Group]{
		Query: `SELECT * FROM groups WHERE name = ?`,
		Args:  []string{"name"},
	}

	GetGroupByID = Statement[models.Group]{
		Query: `SELECT * FROM groups WHERE id = ?`,
		Args:  []string{"id"},
	}

	CreateGroup = Statement[any]{
		Query: `INSERT INTO groups (id, name, description) VALUES (?, ?, ?)`,
		Args:  []string{"id", "name", "description"},
	}

	AddGroupMembership = Statement[any]{
		Query: `INSERT INTO group_memberships (user_id, group_id) VALUES (?, ?)`,
		Args:  []string{"user_id", "group_id"},
	}

	GetUserGroupMemberships = Statement[models.Group]{
		Query: `SELECT
				groups.id as id, groups.name as name, groups.description as description
			FROM group_memberships
			JOIN groups ON groups.id = group_id
			WHERE user_id = ?`,
		Args: []string{"user_id"},
	}

	GetUsersInGroup = Statement[models.User]{
		Query: `SELECT * FROM group_memberships
		JOIN users ON users.id = user_id
		WHERE group_id = ?`,
		Args: []string{"group_id"},
	}
)
