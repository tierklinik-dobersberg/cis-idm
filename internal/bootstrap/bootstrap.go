package bootstrap

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/cis-userd/internal/repo"
	"github.com/tierklinik-dobersberg/cis-userd/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-userd/internal/repo/stmts"
	"golang.org/x/crypto/bcrypt"
)

func Bootstrap(ctx context.Context, userRepo *repo.UserRepo) (*models.User, error) {
	// Bootstrap system groups
	var (
		group models.Group
		err   error
	)

	if group, err = userRepo.GetGroupByID(ctx, "idm_superuser"); err != nil {
		if errors.Is(err, stmts.ErrNoResults) {

			group = models.Group{
				ID:          "idm_superuser",
				Name:        "idm_superuser",
				Description: "Internal management group for super users",
			}

			group, err = userRepo.CreateGroup(ctx, group)
			if err != nil {
				return nil, fmt.Errorf("failed to create idm_superuser group: %w", err)
			}

			logrus.WithField("id", group.ID).Infof("bootstrap: creating idm_superuser group")
		} else {
			return nil, fmt.Errorf("failed to get idm_superuser group: %w", err)
		}
	} else {
		logrus.WithFields(logrus.Fields{"id": group.ID}).Infof("bootstrap: idm_superuser group exists")
	}

	// ensure there is at least one user in idm_superuser
	users, err := userRepo.GetUsersInGroup(ctx, group.ID)
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

		if err := userRepo.AddGroupMembership(ctx, user.ID, group.ID); err != nil {
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

// generateSecret returns a random secret of the given size encoded as hex.
func generateSecret(size int) (string, error) {
	nonce := make([]byte, size)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	return hex.EncodeToString(nonce), nil
}
