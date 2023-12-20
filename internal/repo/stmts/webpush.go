package stmts

import "github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"

var (
	GetWebPushSubscriptionsForUser = Statement[models.Webpush]{
		Query: `SELECT * FROM webpush_subscriptions WHERE user_id = ?`,
		Args:  []string{"user_id"},
	}

	DeleteWebPushSubscriptionForToken = Statement[any]{
		Query: `DELETE FROM webpush_subscriptions WHERE token_id = ?`,
		Args:  []string{"token_id"},
	}

	DeleteWebPushSubscriptionByID = Statement[any]{
		Query: `DELETE FROM webpush_subscriptions WHERE id = ?`,
		Args:  []string{"id"},
	}

	CreateWebPushSubscriptionForUser = Statement[any]{
		Query: `INSERT OR REPLACE INTO webpush_subscriptions (
			id,
			user_id,
			user_agent,
			endpoint,
			auth,
			key,
			token_id
		)
		VALUES ( 
			?, ?, ?, ?, ?, ?, ?
		)`,
		Args: []string{
			"id",
			"user_id",
			"user_agent",
			"endpoint",
			"auth",
			"key",
			"token_id",
		},
	}
)
