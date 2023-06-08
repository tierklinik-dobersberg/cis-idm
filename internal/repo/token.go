package repo

import (
	"context"
	"errors"
	"time"

	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func (repo *Repo) MarkTokenRejected(ctx context.Context, token models.RejectedToken) error {
	if err := stmts.CreateRejectedToken.Write(ctx, repo.Conn, token); err != nil {
		return err
	}

	return nil
}

func (repo *Repo) IsTokenRejected(ctx context.Context, tokenID string) (bool, error) {
	_, err := QueryOne(ctx, stmts.IsTokenRejected, repo.Conn, models.RejectedToken{TokenID: tokenID})
	if err != nil {
		if errors.Is(err, stmts.ErrNoResults) {
			return false, nil
		}

		return true, err
	}

	return true, nil
}

func (repo *Repo) DeleteRejectedTokens(ctx context.Context, threshold time.Time) error {
	return stmts.DeleteExpiredTokens.Write(ctx, repo.Conn, models.RejectedToken{ExiresAt: threshold.Unix()})
}
