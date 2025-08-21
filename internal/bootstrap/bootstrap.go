package bootstrap

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/tierklinik-dobersberg/apis/pkg/data"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

func Bootstrap(ctx context.Context, cfg config.Config, store *repo.Queries) error {
	tx, err := store.Tx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// warp the repository with the transaction tx
	store = store.WithTx(tx)

	// Load all system roles and create a lookup map
	allRoles, err := store.GetSystemRoles(ctx)
	if err != nil {
		return fmt.Errorf("failed to list available roles: %w", err)
	}
	roleIndex := data.IndexSlice(allRoles, func(r repo.Role) string { return r.ID })

	// bootstrap the administrator role
	if err := bootstrapRole(ctx, store, config.Role{
		ID:          "idm_superuser",
		Name:        "idm_superuser",
		Description: "The super user role",
		// NOTE: idm_superuser role does not have permissions assigned by default because
		//       the role itself is considered super-user.
	}); err != nil {
		return fmt.Errorf("role idm_superuser: %w", err)
	}
	delete(roleIndex, "idm_superuser")

	// Create all bootstrap roles
	for idx, role := range cfg.Roles {
		if role.ID == "" || role.Name == "" {
			return fmt.Errorf("invalid role definition in configuration at index %d", idx)
		}

		if err := bootstrapRole(ctx, store, role); err != nil {
			return fmt.Errorf("role %q: %w", role.ID, err)
		}

		delete(roleIndex, role.ID)
	}

	// delete all system roles that were not processed above
	for id := range roleIndex {
		if _, err := store.DeleteRole(ctx, id); err != nil {
			return fmt.Errorf("failed to remove old system role %q: %w", id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func bootstrapRole(ctx context.Context, ds *repo.Queries, role config.Role) error {
	params := repo.CreateSystemRoleParams{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}

	_, err := ds.CreateSystemRole(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	// delete all role permissions and re-create them
	if err := ds.DeleteAllRolePermissions(ctx, role.ID); err != nil {
		return fmt.Errorf("failed to delete existing role permissions: %w", err)
	}

	// re-assing all permissions from the configuration
	for _, perm := range role.Permissions {
		if err := ds.AssignPermissionToRole(ctx, repo.AssignPermissionToRoleParams{
			Permission: perm,
			RoleID:     role.ID,
		}); err != nil {
			return fmt.Errorf("failed to assign permission %q", perm)
		}
	}
	log.L(ctx).Debug("bootstrap: successfully created or  updated system role", "id", role.ID, "name", role.Name, "permissions", role.Permissions)

	return nil
}

// GenerateSecret returns a random secret of the given size encoded as hex.
func GenerateSecret(size int) (string, error) {
	nonce := make([]byte, size)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	return hex.EncodeToString(nonce), nil
}
