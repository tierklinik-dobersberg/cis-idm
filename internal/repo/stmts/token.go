package stmts

import "github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"

var (
	CreateRejectedToken = Statement[any]{
		Query: `INSERT INTO token_invalidation (token_id, user_id, expires_at, issued_at) VALUES (?, ?, ?, ?)`,
		Args:  []string{"token_id", "user_id", "expires_at", "issued_at"},
	}

	IsTokenRejected = Statement[models.RejectedToken]{
		Query: `SELECT * FROM token_invalidation WHERE token_id = ?`,
		Args:  []string{"token_id"},
	}

	DeleteExpiredTokens = Statement[any]{
		Query: `DELETE FROM token_invalidation WHERE expires_at < ?`,
		Args:  []string{"expires_at"},
	}

	CreateRegistrationToken = Statement[any]{
		Query: `INSERT INTO registration_tokens (
			token,
			expires,
			allowed_usage,
			initial_roles,
			created_by,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?)`,
		Args: []string{"token", "expires", "allowed_usage", "initial_roles", "created_by", "created_at"},
	}

	GetRegistrationToken = Statement[models.RegistrationToken]{
		Query: `SELECT * FROM registration_tokens WHERE token = ? AND allowed_usage > 0 AND (expires IS NULL OR expires > ?)`,
		Args:  []string{"token", "expires"},
	}

	MarkRegistrationTokenUsed = Statement[models.RegistrationToken]{
		Query: `UPDATE registration_tokens SET allowed_usage = (
			CASE 
				WHEN allowed_usage IS NOT NULL THEN allowed_usage - 1
				ELSE NULL
			END
		) WHERE token = ? AND (expires IS NULL OR expires > ?) RETURNING *`,
		Args: []string{"token", "expires"},
	}

	InsertRecoveryCodes = Statement[any]{
		Query: `INSERT INTO mfa_backup_codes (code, user_id) VALUES (?, ?)`,
		Args:  []string{"code", "user_id"},
	}

	CheckAndDeleteRecoveryCode = Statement[any]{
		Query: `DELETE FROM mfa_backup_codes WHERE user_id = ? AND code = ?`,
		Args:  []string{"user_id", "code"},
	}

	RemoveAllRecoveryCodes = Statement[any]{
		Query: `DELETE FROM mfa_backup_codes WHERE user_id = ?`,
		Args:  []string{"user_id"},
	}

	LoadUserRecoveryCodes = Statement[any]{
		Query: `SELECT * FROM mfa_backup_codes WHERE user_id = ?`,
		Args:  []string{"user_id"},
	}
)
