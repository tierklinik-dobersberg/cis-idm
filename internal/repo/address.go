package repo

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func (repo *Repo) GetUserAddresses(ctx context.Context, userID string) ([]models.Address, error) {
	return Query(ctx, stmts.GetUserAddressesByID, repo.Conn, map[string]any{
		"user_id": userID,
	})
}

func (repo *Repo) GetAddressesByID(ctx context.Context, userID, id string) (models.Address, error) {
	return QueryOne(ctx, stmts.GetAddressByID, repo.Conn, map[string]any{
		"user_id": userID,
		"id":      id,
	})
}

func (repo *Repo) AddUserAddress(ctx context.Context, address models.Address) (models.Address, error) {
	if address.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return address, err
		}

		address.ID = id.String()
	}

	if err := stmts.CreateAddress.Write(ctx, repo.Conn, address); err != nil {
		return address, err
	}

	return address, nil
}

func (repo *Repo) UpdateUserAddress(ctx context.Context, address models.Address) error {
	if address.ID == "" {
		return fmt.Errorf("missing address id")
	}

	return stmts.UpdateAddress.Write(ctx, repo.Conn, address)
}

func (repo *Repo) DeleteUserAddress(ctx context.Context, userID, addressID string) error {
	return stmts.DeleteAddress.Write(ctx, repo.Conn, map[string]any{
		"id":      addressID,
		"user_id": userID,
	})
}
