package repo

import (
	"context"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
)

type WebauthnUser struct {
	models.User
	repo *Repo
	ctx  context.Context
	log  *logrus.Entry
}

func NewWebAuthnUser(ctx context.Context, log *logrus.Entry, repo *Repo, usr models.User) webauthn.User {
	return &WebauthnUser{
		User: usr,
		repo: repo,
		ctx:  ctx,
		log:  log.WithField("username", usr.Username).WithField("user_id", usr.ID),
	}
}

func (usr *WebauthnUser) WebAuthnID() []byte {
	return []byte(usr.ID)
}

func (usr *WebauthnUser) WebAuthnName() string {
	return usr.Username
}

func (usr *WebauthnUser) WebAuthnDisplayName() string {
	if usr.DisplayName != "" {
		return usr.DisplayName
	}

	return usr.Username
}

func (usr *WebauthnUser) WebAuthnCredentials() []webauthn.Credential {
	usr.log.Infof("searching for webauthn credentials")

	res, err := usr.repo.GetWebauthnCreds(usr.ctx, usr.ID)
	if err != nil {
		usr.log.Errorf("failed to fetch webauthn credentials: %s", err)

		return nil
	}

	if len(res) == 0 {
		usr.log.Errorf("user does not have any webauthn credentials yet.")
	}

	return res
}

func (usr *WebauthnUser) WebAuthnIcon() string {
	return usr.Avatar
}
