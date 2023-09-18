package webauthn

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gofrs/uuid"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
)

func (svc *Service) BeginLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		options *protocol.CredentialAssertion
		session *webauthn.SessionData
	)

	pathParts := strings.Split(r.URL.Path, "/")
	userNameOrEmail := pathParts[len(pathParts)-1]

	if userNameOrEmail != "" {
		user, err := svc.Datastore.GetUserByName(ctx, userNameOrEmail)
		if err != nil {
			if svc.Config.FeatureEnabled(config.FeatureLoginByMail) {
				var verified bool
				user, verified, err = svc.Datastore.GetUserByEMail(ctx, userNameOrEmail)

				if err == nil && !verified {
					http.Error(w, "e-mail address not verified", http.StatusPreconditionFailed)

					return
				}
			}
		}

		webauthnUser := repo.NewWebAuthnUser(
			ctx,
			log.L(ctx),
			svc.Datastore,
			user,
		)

		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)

			return
		}

		options, session, err = svc.web.BeginLogin(webauthnUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	} else {
		var err error

		options, session, err = svc.web.BeginDiscoverableLogin()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}

	sessionID, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := svc.Cache.PutKeyTTL(ctx, sessionID.String(), session, time.Until(session.Expires)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "login_session",
		Value:    sessionID.String(),
		Secure:   *svc.Config.SecureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Minute * 5),
		Path:     "/",
	})

	jsonResponse(w, options, http.StatusOK)
}

func (svc *Service) FinishLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	response, err := protocol.ParseCredentialRequestResponseBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestedRedirect := r.URL.Query().Get("redirect")
	if requestedRedirect != "" {
		requestedRedirect, err = svc.HandleRequestedRedirect(ctx, requestedRedirect)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}
	}

	cookie := middleware.FindCookie("login_session", r.Header)
	if cookie == nil {
		http.Error(w, "cookie not found", http.StatusBadRequest)
		return
	}

	var session webauthn.SessionData
	if err := svc.Cache.GetAndDeleteKey(ctx, cookie.Value, &session); err != nil {
		log.L(ctx).Errorf("failed to get webauthn login session for key %s: %s", cookie.Value, err)
		http.Error(w, "session not found: "+err.Error(), http.StatusNotFound)

		return
	}

	var user models.User
	getUserID := func(rawID, userHandle []byte) (webauthn.User, error) {
		var err error
		user, err = svc.Datastore.GetUserByID(ctx, string(userHandle))
		if err != nil {
			return nil, fmt.Errorf("user not found")
		}

		webauthnUser := repo.NewWebAuthnUser(
			ctx,
			log.L(ctx),
			svc.Datastore,
			user,
		)

		return webauthnUser, nil
	}

	if len(session.UserID) > 0 {
		webauthnUser, err := getUserID(nil, session.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)

			return
		}

		_, err = svc.web.ValidateLogin(webauthnUser, session, response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)

			return
		}
	} else {
		_, err := svc.web.ValidateDiscoverableLogin(getUserID, session, response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)

			return
		}
	}

	roles, err := svc.Datastore.GetUserRoles(ctx, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// Generate and add refresh and access tokens

	_, refreshTokenID, err := svc.AddRefreshToken(user, roles, w.Header())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if _, _, err := svc.AddAccessToken(user, roles, 0, refreshTokenID, w.Header()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// Prepare the response.

	userResponse := make(map[string]any)
	if requestedRedirect != "" {
		userResponse["redirectTo"] = requestedRedirect
	}

	jsonResponse(w, userResponse, http.StatusOK)
}
