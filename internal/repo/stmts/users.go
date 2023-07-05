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

	EnrollUserTOTPSecret = Statement[any]{
		Query: `UPDATE users SET totp_secret = ? WHERE id = ? AND totp_secret IS NULL`,
		Args:  []string{"totp_secret", "id"},
	}

	RemoveUserTOTPSecret = Statement[any]{
		Query: `UPDATE users SET totp_secret = NULL WHERE id = ? AND totp_secret IS NOT NULL`,
		Args:  []string{"id"},
	}

	SetUserPassword = Statement[any]{
		Query: `UPDATE users SET password = ? WHERE id = ?`,
		Args:  []string{"password", "id"},
	}

	AddWebauthnCred = Statement[any]{
		Query: `INSERT INTO webauthn_creds (id, user_id, cred, client_name, client_os, client_device, cred_type) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		Args:  []string{"id", "user_id", "cred", "client_name", "client_os", "client_device", "cred_type"},
	}

	GetWebauthnCreds = Statement[models.Passkey]{
		Query: `SELECT * FROM webauthn_creds WHERE user_id = ?`,
		Args:  []string{"user_id"},
	}

	RemoveWebauthnCred = Statement[any]{
		Query: `DELETE FROM webauthn_creds WHERE user_id = ? AND id = ?`,
		Args:  []string{"user_id", "id"},
	}

	SaveWebauthnSession = Statement[any]{
		Query: `INSERT INTO webauthn_sessions (id, user_id, session) VALUES (?, ?, ?)`,
		Args:  []string{"id", "user_id", "session"},
	}

	GetWebauthnSession = Statement[struct {
		ID      string `mapstructure:"id"`
		UserID  string `mapstructure:"user_id"`
		Session string `mapstructure:"session"`
	}]{
		Query: `SELECT * FROM webauthn_sessions WHERE id = ?`,
		Args:  []string{"id"},
	}

	DeleteWebauthnSession = Statement[any]{
		Query: `DELETE FROM webauthn_sessions WHERE id = ?`,
		Args:  []string{"id"},
	}

	// FIXME(ppacher): create trigger that removes stale session entries ...
)
