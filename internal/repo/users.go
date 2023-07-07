package repo

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gofrs/uuid"
	"github.com/mileusna/useragent"
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

func (repo *Repo) GetUsers(ctx context.Context) ([]models.User, error) {
	return Query(ctx, stmts.GetAllUsers, repo.Conn, nil)
}

func (repo *Repo) DeleteUser(ctx context.Context, id string) error {
	return stmts.DeleteUser.Write(ctx, repo.Conn, map[string]any{"id": id})
}

func (repo *Repo) SetUserTotpSecret(ctx context.Context, userID string, totpSecret string) error {
	return stmts.EnrollUserTOTPSecret.Write(ctx, repo.Conn, map[string]any{
		"id":          userID,
		"totp_secret": totpSecret,
	})
}

func (repo *Repo) RemoveUserTotpSecret(ctx context.Context, userID string) error {
	return stmts.RemoveUserTOTPSecret.Write(ctx, repo.Conn, map[string]any{"id": userID})
}

func (repo *Repo) AddWebauthnCred(ctx context.Context, userID string, cred webauthn.Credential, ua useragent.UserAgent) error {
	blob, err := json.Marshal(cred)
	if err != nil {
		return err
	}

	return stmts.AddWebauthnCred.Write(ctx, repo.Conn, map[string]any{
		"id":            hex.EncodeToString(cred.ID),
		"user_id":       userID,
		"cred":          string(blob),
		"client_name":   ua.Name,
		"client_os":     ua.OS,
		"client_device": ua.Device,
		"cred_type":     cred.Authenticator.Attachment,
	})
}

func (repo *Repo) GetPasskeys(ctx context.Context, userID string) ([]models.Passkey, error) {
	return Query(ctx, stmts.GetWebauthnCreds, repo.Conn, map[string]any{
		"user_id": userID,
	})
}

func (repo *Repo) GetWebauthnCreds(ctx context.Context, userID string) ([]webauthn.Credential, error) {
	res, err := Query(ctx, stmts.GetWebauthnCreds, repo.Conn, map[string]any{"user_id": userID})
	if err != nil {
		return nil, err
	}

	creds := make([]webauthn.Credential, len(res))
	for idx, r := range res {
		var c webauthn.Credential
		if err := json.Unmarshal([]byte(r.Cred), &c); err != nil {
			return nil, fmt.Errorf("failed to parse webauthn credentials")
		}

		creds[idx] = c
	}

	return creds, nil
}

func (repo *Repo) RemoveWebauthnCred(ctx context.Context, userID, id string) error {
	return stmts.RemoveWebauthnCred.Write(ctx, repo.Conn, map[string]any{
		"user_id": userID,
		"id":      id,
	})
}
