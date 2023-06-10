package repo

import (
	"context"

	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func (repo *Repo) GetUserAddresses(ctx context.Context, userID string) ([]models.Address, error) {
	return Query(ctx, stmts.GetUserAddressesByID, repo.Conn, map[string]any{
		"user_id": userID,
	})
}
