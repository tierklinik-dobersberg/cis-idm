package bootstrap

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
	"golang.org/x/crypto/bcrypt"
)

func Bootstrap(ctx context.Context, cfg config.Config, userRepo *repo.Repo) (*models.User, error) {
	superuserRoleID, err := bootstrapRole(ctx, userRepo, "idm_superuser", "Super-user management role", true)
	if err != nil {
		return nil, err
	}

	// Create all bootstrap roles
	for _, roleName := range cfg.BootstrapRoles {
		if _, err := bootstrapRole(ctx, userRepo, roleName, "Automatically bootstrapped role", true); err != nil {
			return nil, err
		}
	}

	// ensure there is at least one user in idm_superuser
	users, err := userRepo.GetUsersByRole(ctx, superuserRoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users in idm_superuser group: %w", err)
	}

	if len(users) == 0 {
		password, err := generateSecret(16)
		if err != nil {
			return nil, fmt.Errorf("failed to generate secure password: %w", err)
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to generate hashed password: %w", err)
		}

		user := models.User{
			ID:       "admin",
			Username: "admin",
			Password: string(hashed),
		}

		user, err = userRepo.CreateUser(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to create new super user: %w", err)
		}
		logrus.WithFields(logrus.Fields{
			"id":       user.ID,
			"name":     user.Username,
			"password": password,
		}).Infof("bootstrap: created new super user")

		if err := userRepo.AssignRoleToUser(ctx, user.ID, superuserRoleID); err != nil {
			return nil, fmt.Errorf("failed to add new user to idm_superuser group: %w", err)
		}

		return &user, nil
	} else {
		logrus.WithFields(logrus.Fields{"count": len(users)}).Infof("bootstrap: found users in idm_superuser group")
		for _, user := range users {
			logrus.WithField("id", user.ID).Infof("idm_superuser: %s", user.Username)
		}
	}

	return nil, nil
}

func bootstrapRole(ctx context.Context, repo *repo.Repo, roleName, description string, deleteProtection bool) (string, error) {
	role, err := repo.GetRoleByID(ctx, roleName)
	if err != nil {
		if errors.Is(err, stmts.ErrNoResults) {

			role = models.Role{
				ID:              roleName,
				Name:            roleName,
				Description:     description,
				DeleteProtected: deleteProtection,
			}

			role, err = repo.CreateRole(ctx, role)
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

// generateSecret returns a random secret of the given size encoded as hex.
func generateSecret(size int) (string, error) {
	nonce := make([]byte, size)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	return hex.EncodeToString(nonce), nil
}
