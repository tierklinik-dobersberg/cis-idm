package stmts

import "github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"

var (
	CreateEMail = Statement[any]{
		Query: `INSERT INTO user_emails (
			id,
			user_id,
			address,
			verified
		)
		VALUES (?, ?, ?, ?)`,
		Args: []string{"id", "user_id", "address", "verified"},
	}

	GetEmailsForUserByID = Statement[models.EMail]{
		Query: `SELECT * FROM user_emails WHERE user_id = ?`,
		Args:  []string{"user_id"},
	}

	DeleteEMailFromUser = Statement[any]{
		Query: `DELETE FROM user_emails WHERE id = ? AND user_id = ?`,
		Args:  []string{"id", "user_id"},
	}

	MarkEmailVerified = Statement[any]{
		Query: `UPDATE user_emails SET verified = ? WHERE id = ? AND user_id = ?`,
		Args:  []string{"verified", "id", "user_id"},
	}
)
