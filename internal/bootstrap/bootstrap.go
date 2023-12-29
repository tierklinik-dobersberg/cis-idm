package bootstrap

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

func Bootstrap(ctx context.Context, cfg config.Config, userRepo *repo.Queries) error {
	superuserRoleID, err := bootstrapRole(ctx, userRepo, "idm_superuser", "Super-user management role", true)
	if err != nil {
		return err
	}

	// Create all bootstrap roles
	for _, roleName := range cfg.BootstrapRoles {
		if _, err := bootstrapRole(ctx, userRepo, roleName, "Automatically bootstrapped role", true); err != nil {
			return err
		}
	}

	// ensure there is at least one user in idm_superuser
	users, err := userRepo.GetUsersByRole(ctx, superuserRoleID)
	if err != nil {
		return fmt.Errorf("failed to retrieve users in idm_superuser group: %w", err)
	}

	// if there are not users with the idm_superuser role create a new registration token
	if len(users) == 0 {
		tokenValue, err := GenerateSecret(8)
		if err != nil {
			return err
		}
		blobs, _ := json.Marshal([]string{"idm_superuser"})

		token := repo.CreateRegistrationTokenParams{
			Token:        tokenValue,
			InitialRoles: string(blobs),
			AllowedUsage: sql.NullInt64{
				Int64: 1,
				Valid: true,
			},
		}

		if err := userRepo.CreateRegistrationToken(ctx, token); err != nil {
			return err
		}

		logrus.WithField("token", tokenValue).Infof("Please bootstrap the superuser account using the provided registration token.")

		return nil
	}

	logrus.WithFields(logrus.Fields{"count": len(users)}).Infof("bootstrap: found users in idm_superuser group")
	for _, user := range users {
		logrus.WithField("id", user.ID).Infof("idm_superuser: %s", user.Username)
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
