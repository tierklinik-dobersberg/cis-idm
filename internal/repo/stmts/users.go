package stmts

import "github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"

var (
	GetUserByName = Statement[models.User]{
		Query: `SELECT * FROM users WHERE username = ?`,
		Args:  []string{"username"},
	}

	GetUserByEMail = Statement[struct {
		models.User  `mapstructure:",squash"`
		MailVerified bool `mapstructure:"verified"`
	}]{
		Query: `SELECT users.*, user_emails.verified FROM users
			JOIN user_emails ON user_emails.user_id = users.id WHERE user_emails.address = ?`,
		Args: []string{"mail"},
	}

	GetUserByID = Statement[models.User]{
		Query: `SELECT * FROM users WHERE id = ?`,
		Args:  []string{"id"},
	}

	GetAllUsers = Statement[models.User]{
		Query: `SELECT * FROM users`,
	}

	DeleteUser = Statement[any]{
		Query: `DELETE FROM users WHERE id = ?`,
		Args:  []string{"id"},
	}

	CreateUser = Statement[any]{
		Query: `INSERT INTO users (
			id,
			username,
			display_name,
			first_name,
			last_name,
			extra,
			avatar,
			birthday,
			password
		)
		VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?
		)`,
		Args: []string{
			"id",
			"username",
			"display_name",
			"first_name",
			"last_name",
			"extra",
			"avatar",
			"birthday",
			"password",
		},
	}

	UpdateUser = Statement[any]{
		Query: `UPDATE users SET
			username = ?,
			display_name = ?,
			first_name = ?,
			last_name = ?,
			extra = ?,
			avatar = ?,
			birthday = ?
		WHERE id = ?
			`,
		Args: []string{
			"username",
			"display_name",
			"first_name",
			"last_name",
			"extra",
			"avatar",
			"birthday",
			"id",
		},
	}

	SetUserPassword = Statement[any]{
		Query: `UPDATE users SET password = ? WHERE id = ?`,
		Args:  []string{"password", "id"},
	}
)
