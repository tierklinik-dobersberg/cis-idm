package repo

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func (repo *Repo) GetUserByName(ctx context.Context, name string) (models.User, error) {
	return QueryOne(ctx, stmts.GetUserByName, repo.Conn, map[string]any{"username": name})
}

func (repo *Repo) GetUserByEMail(ctx context.Context, name string) (models.User, bool, error) {
	result, err := QueryOne(ctx, stmts.GetUserByEMail, repo.Conn, map[string]any{"mail": name})
	if err != nil {
		return models.User{}, false, err
	}

	return result.User, result.MailVerified, nil
}

func (repo *Repo) GetRoleByName(ctx context.Context, name string) (models.Role, error) {
	return QueryOne(ctx, stmts.GetRoleByName, repo.Conn, map[string]any{
		"name": name,
	})
}

func (repo *Repo) GetRoleByID(ctx context.Context, id string) (models.Role, error) {
	return QueryOne(ctx, stmts.GetRoleByID, repo.Conn, map[string]any{
		"id": id,
	})
}

func (repo *Repo) GetRoles(ctx context.Context) ([]models.Role, error) {
	return Query(ctx, stmts.GetRoles, repo.Conn, nil)
}

func (repo *Repo) DeleteRole(ctx context.Context, id string) error {
	return stmts.DeleteRole.Write(ctx, repo.Conn, map[string]any{
		"id": id,
	})
}

func (repo *Repo) CreateRole(ctx context.Context, group models.Role) (models.Role, error) {
	if group.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return models.Role{}, err
		}

		group.ID = id.String()
	}

	if err := stmts.CreateRole.Write(ctx, repo.Conn, group); err != nil {
		return models.Role{}, err
	}

	return group, nil
}

func (repo *Repo) UpdateRole(ctx context.Context, role models.Role) error {
	return stmts.UpdateRole.Write(ctx, repo.Conn, role)
}

func (repo *Repo) AssignRoleToUser(ctx context.Context, userID string, roleID string) error {
	return stmts.AssignRoleToUser.Write(ctx, repo.Conn, models.RoleAssignment{
		UserID: userID,
		RoleID: roleID,
	})
}

func (repo *Repo) UnassignRoleFromUser(ctx context.Context, userID string, roleID string) error {
	return stmts.UnassignRoleFromUser.Write(ctx, repo.Conn, models.RoleAssignment{
		UserID: userID,
		RoleID: roleID,
	})
}

func (repo *Repo) GetUsersByRole(ctx context.Context, roleID string) ([]models.User, error) {
	return Query(ctx, stmts.GetUsersByRole, repo.Conn, models.RoleAssignment{
		RoleID: roleID,
	})
}
