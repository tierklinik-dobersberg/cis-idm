package stmts

import "github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"

var (
	GetUserPhoneNumbersByID = Statement[models.PhoneNumber]{
		Query: `SELECT * FROM user_phone_numbers WHERE user_id = ?`,
		Args:  []string{"user_id"},
	}
)
