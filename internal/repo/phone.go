package repo

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func (repo *Repo) GetUserPhoneNumbers(ctx context.Context, userID string) ([]models.PhoneNumber, error) {
	return Query(ctx, stmts.GetPhoneNumbersByUserID, repo.Conn, map[string]any{
		"user_id": userID,
	})
}

func (repo *Repo) GetUserPhoneNumberByID(ctx context.Context, id, userID string) (models.PhoneNumber, error) {
	return QueryOne(ctx, stmts.GetPhoneNumberByID, repo.Conn, map[string]any{
		"user_id": userID,
		"id":      id,
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

	if err := stmts.CreateUserPhoneNumber.Write(ctx, repo.Conn, model); err != nil {
		return model, err
	}

	return model, nil
}

func (repo *Repo) DeleteUserPhoneNumber(ctx context.Context, userID string, phoneNumberID string) error {
	return stmts.DeleteUserPhoneNumber.Write(ctx, repo.Conn, map[string]any{
		"id":      phoneNumberID,
		"user_id": userID,
	})
}

func (repo *Repo) MarkPhoneNumberAsPrimary(ctx context.Context, userID string, phoneNumberID string) error {
	return stmts.MarkPhoneNumberAsPrimary.Write(ctx, repo.Conn, map[string]any{
		"user_id": userID,
		"id":      phoneNumberID,
	})
}

func (repo *Repo) MarkPhoneNumberAsVerified(ctx context.Context, userID string, phoneNumberID string) error {
	return stmts.MarkPhoneNumberAsVerified.Write(ctx, repo.Conn, map[string]any{
		"user_id": userID,
		"id":      phoneNumberID,
	})
}
