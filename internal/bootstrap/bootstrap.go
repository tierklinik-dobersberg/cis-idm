package bootstrap

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

func Bootstrap(ctx context.Context, cfg config.Config, userRepo *repo.Queries) error {
	_, err := bootstrapRole(ctx, userRepo, "idm_superuser", "Super-user management role", true)
	if err != nil {
		return err
	}

	// Create all bootstrap roles
	for _, roleName := range cfg.Roles {
		if _, err := bootstrapRole(ctx, userRepo, roleName, "Automatically bootstrapped role", true); err != nil {
			return err
		}
	}

	return nil
}

func bootstrapRole(ctx context.Context, ds *repo.Queries, roleName, description string, deleteProtection bool) (string, error) {
	role, err := ds.GetRoleByID(ctx, roleName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			params := repo.CreateRoleParams{
				ID:              roleName,
				Name:            roleName,
				Description:     description,
				DeleteProtected: deleteProtection,
			}

			role, err = ds.CreateRole(ctx, params)
			if err != nil {
				return "", fmt.Errorf("failed to create role %s: %w", roleName, err)
			}

			logrus.
				WithField("id", role.ID).
				WithField("name", role.Name).
				Infof("bootstrap: successfully created role")
		} else {
			return "", fmt.Errorf("failed to get idm_superuser group: %w", err)
		}
	} else {
		logrus.
			WithField("id", role.ID).
			WithField("name", role.Name).
			Infof("bootstrap: role already created")
	}

	return role.ID, nil
}

// GenerateSecret returns a random secret of the given size encoded as hex.
func GenerateSecret(size int) (string, error) {
	nonce := make([]byte, size)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	return hex.EncodeToString(nonce), nil
}
