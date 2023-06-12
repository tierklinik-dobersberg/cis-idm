package repo

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func (repo *Repo) GetUserByID(ctx context.Context, id string) (models.User, error) {
	return QueryOne(ctx, stmts.GetUserByID, repo.Conn, map[string]any{"id": id})
}

func (repo *Repo) GetUserRoles(ctx context.Context, userID string) ([]models.Role, error) {
	return Query(ctx, stmts.GetRolesForUser, repo.Conn, models.RoleAssignment{
		UserID: userID,
	})
}

func (repo *Repo) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	if user.ID == "" {
		userID, err := uuid.NewV4()
		if err != nil {
			return user, err
		}

		user.ID = userID.String()
	}

	if err := stmts.CreateUser.Write(ctx, repo.Conn, user); err != nil {
		return user, err
	}

	return user, nil
}

func (repo *Repo) UpdateUser(ctx context.Context, user models.User) error {
	return stmts.UpdateUser.Write(ctx, repo.Conn, user)
}

func (repo *Repo) SetUserPassword(ctx context.Context, userID string, password string) error {
	return stmts.SetUserPassword.Write(ctx, repo.Conn, models.User{
		ID:       userID,
		Password: password,
	})
}
