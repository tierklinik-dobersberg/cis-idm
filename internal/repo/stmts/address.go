package stmts

import "github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"

var (
	GetUserAddressesByID = Statement[models.Address]{
		Query: `SELECT * FROM user_addresses WHERE user_id = ?`,
		Args:  []string{"user_id"},
	}

	GetAddressByID = Statement[models.Address]{
		Query: `SELECT * FROM user_addresses WHERE user_id = ? AND id = ?`,
		Args:  []string{"user_id", "id"},
	}

	CreateAddress = Statement[any]{
		Query: `INSERT INTO user_addresses (
			id,
			user_id,
			city_code,
			city_name,
			street,
			extra
		)
		VALUES (?, ?, ?, ?, ?, ?)`,
		Args: []string{"id", "user_id", "city_code", "city_name", "street", "extra"},
	}

	UpdateAddress = Statement[any]{
		Query: `UPDATE user_addresses SET city_code = ?, city_name = ?, street = ?, extra = ? WHERE id = ? AND user_id = ?`,
		Args:  []string{"city_code", "city_name", "street", "extra", "id", "user_id"},
	}

	DeleteAddress = Statement[any]{
		Query: `DELETE FROM user_addresses WHERE id = ? AND user_id = ?`,
		Args:  []string{"id", "user_id"},
	}
)
