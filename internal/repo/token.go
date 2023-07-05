package repo

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/rqlite/gorqlite"
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
	return stmts.DeleteExpiredTokens.Write(ctx, repo.Conn, models.RejectedToken{ExpiresAt: threshold.Unix()})
}

func (repo *Repo) CreateRegistrationToken(ctx context.Context, token models.RegistrationToken) error {
	if token.CreatedAt == 0 {
		token.CreatedAt = time.Now().Unix()
	}

	return stmts.CreateRegistrationToken.Write(ctx, repo.Conn, token)
}

func (repo *Repo) ValidateRegistrationToken(ctx context.Context, token string) (models.RegistrationToken, error) {
	return QueryOne(ctx, stmts.GetRegistrationToken, repo.Conn, map[string]any{
		"token":   token,
		"expires": time.Now().Unix(),
	})
}

func (repo *Repo) MarkRegistrationTokenUsed(ctx context.Context, token string) error {
	return stmts.MarkRegistrationTokenUsed.Write(ctx, repo.Conn, map[string]any{
		"token":   token,
		"expires": time.Now().Unix(),
	})
}

func (repo *Repo) ReplaceUserRecoveryCodes(ctx context.Context, userID string, codes []string) error {
	sqlStmts := make([]gorqlite.ParameterizedStatement, len(codes)+1)
	sqlStmts[0], _ = stmts.RemoveAllRecoveryCodes.Prepare(map[string]any{"user_id": userID})

	for idx, code := range codes {
		prepared, err := stmts.InsertRecoveryCodes.Prepare(map[string]any{
			"user_id": userID,
			"code":    code,
		})
		if err != nil {
			return err
		}

		sqlStmts[idx+1] = prepared
	}

	results, err := repo.Conn.WriteParameterizedContext(ctx, sqlStmts)
	if err != nil {
		merr := new(multierror.Error)
		merr.Errors = append(merr.Errors, err)
		for _, res := range results {
			if res.Err != nil {
				merr.Errors = append(merr.Errors, res.Err)
			}
		}

		return merr.ErrorOrNil()
	}

	return nil
}

func (repo *Repo) CheckAndDeleteRecoveryCode(ctx context.Context, userID string, code string) error {
	return stmts.CheckAndDeleteRecoveryCode.Write(ctx, repo.Conn, map[string]any{
		"user_id": userID,
		"code":    code,
	})
}

func (repo *Repo) UserHasRecoveryCodes(ctx context.Context, userID string) (bool, error) {
	_, err := QueryOne(ctx, stmts.LoadUserRecoveryCodes, repo.Conn, map[string]any{"user_id": userID})
	if err != nil {
		if errors.Is(err, stmts.ErrNoResults) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
