package repo

import (
	"context"

	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func (repo *Repo) GetUserPhoneNumbers(ctx context.Context, userID string) ([]models.PhoneNumber, error) {
	return Query(ctx, stmts.GetUserPhoneNumbersByID, repo.Conn, map[string]any{
		"user_id": userID,
	})
}
