package stmts

import "github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"

var (
	GetUserAddressesByID = Statement[models.Address]{
		Query: `SELECT * FROM user_addresses WHERE user_id = ?`,
		Args:  []string{"user_id"},
	}
)
