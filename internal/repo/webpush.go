package repo

import (
	"context"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/gofrs/uuid"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func (repo *Repo) GetWebPushSubscriptionsForUser(ctx context.Context, userID string) ([]models.Webpush, error) {
	return Query(ctx, stmts.GetWebPushSubscriptionsForUser, repo.Conn, map[string]any{
		"user_id": userID,
	})
}

func (repo *Repo) DeleteWebPushSubscriptionForToken(ctx context.Context, tokenID string) error {
	return stmts.DeleteWebPushSubscriptionForToken.Write(ctx, repo.Conn, map[string]any{
		"token_id": tokenID,
	})
}

func (repo *Repo) DeleteWebPushSubscriptionByID(ctx context.Context, id string) error {
	return stmts.DeleteWebPushSubscriptionByID.Write(ctx, repo.Conn, map[string]any{
		"id": id,
	})
}

func (repo *Repo) CreateWebPushSubscriptionForUser(ctx context.Context, userID string, userAgent string, tokenID string, sub webpush.Subscription) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	if err := stmts.CreateWebPushSubscriptionForUser.Write(ctx, repo.Conn, map[string]any{
		"id":         id.String(),
		"user_id":    userID,
		"user_agent": userAgent,
		"endpoint":   sub.Endpoint,
		"auth":       sub.Keys.Auth,
		"key":        sub.Keys.P256dh,
		"token_id":   tokenID,
	}); err != nil {
		return "", err
	}

	return id.String(), nil
}
