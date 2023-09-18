package stmts

import "github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"

var (
	GetPhoneNumbersByUserID = Statement[models.PhoneNumber]{
		Query: `SELECT * FROM user_phone_numbers WHERE user_id = ?`,
		Args:  []string{"user_id"},
	}

	GetUserPrimaryPhoneNumber = Statement[models.PhoneNumber]{
		Query: `SELECT * FROM user_phone_numbers WHERE user_id = ? AND is_primary = TRUE`,
		Args:  []string{"user_id"},
	}

	GetPhoneNumberByID = Statement[models.PhoneNumber]{
		Query: `SELECT * FROM user_phone_numbers WHERE user_id = ? AND id = ?`,
		Args:  []string{"user_id", "id"},
	}

	CreateUserPhoneNumber = Statement[any]{
		Query: `INSERT INTO user_phone_numbers (id, user_id, phone_number, is_primary, verified) VALUES (?, ?, ?, ?, ?)`,
		Args:  []string{"id", "user_id", "phone_number", "is_primary", "verified"},
	}

	DeleteUserPhoneNumber = Statement[any]{
		Query: `DELETE FROM user_phone_numbers WHERE user_id = ? AND id = ?`,
		Args:  []string{"user_id", "id"},
	}

	MarkPhoneNumberVerified = Statement[any]{
		Query: `UPDATE user_phone_numbers SET verified = ? WHERE id = ? AND user_id = ?`,
		Args:  []string{"verified", "id", "user_id"},
	}

	MarkPhoneNumberAsPrimary = Statement[any]{
		Query: `UPDATE user_phone_numbers SET is_primary = (id == ?) WHERE user_id = ?`,
		Args:  []string{"id", "user_id"},
	}

	MarkPhoneNumberAsVerified = Statement[any]{
		Query: `UPDATE user_phone_numbers SET verified = TRUE WHERE user_id = ? AND id = ?`,
		Args:  []string{"user_id", "id"},
	}
)
