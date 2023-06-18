package repo

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func (repo *Repo) GetUserPhoneNumbers(ctx context.Context, userID string) ([]models.PhoneNumber, error) {
	return Query(ctx, stmts.GetUserPhoneNumbersByID, repo.Conn, map[string]any{
		"user_id": userID,
	})
}

func (repo *Repo) AddUserPhoneNumber(ctx context.Context, model models.PhoneNumber) (models.PhoneNumber, error) {
	if model.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return model, err
		}

		model.ID = id.String()
	}

	if err := stmts.CreatePhoneNumber.Write(ctx, repo.Conn, model); err != nil {
		return model, err
	}

	return model, nil
}
