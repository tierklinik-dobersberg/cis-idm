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
)
