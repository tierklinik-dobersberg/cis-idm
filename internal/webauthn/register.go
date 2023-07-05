package webauthn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/mileusna/useragent"
	"github.com/tierklinik-dobersberg/cis-idm/internal/auth"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
)

func New(cfg config.Config, authService *auth.AuthService, repo *repo.Repo) (http.Handler, error) {
	mux := http.NewServeMux()

	wconfig := &webauthn.Config{
		RPDisplayName: cfg.SiteName,
		RPID:          cfg.Domain,
		RPOrigins: []string{
			cfg.PublicURL,
		},
	}

	w, err := webauthn.New(wconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create webauthn instance: %w", err)
	}

	mux.Handle("/registration/begin", beginRegistrationHandler(cfg, authService, w, repo))
	mux.Handle("/registration/finish", finishRegistrationHandler(cfg, w, repo))
	mux.Handle("/login/begin/", beginLoginHandler(cfg, w, repo))
	mux.Handle("/login/finish", finishLoginHandler(cfg, authService, w, repo))

	return mux, nil
}

func beginRegistrationHandler(cfg config.Config, auth *auth.AuthService, web *webauthn.WebAuthn, datastore *repo.Repo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := middleware.L(ctx)

		l.Infof("received request to begin webauthn registration")

		var user models.User
		claims := middleware.ClaimsFromContext(ctx)
		if claims == nil {
			var payload struct {
				Username string `json:"username"`
				Token    string `json:"token"`
			}

			dec := json.NewDecoder(r.Body)
			if err := dec.Decode(&payload); err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			// a user is performing an initial registration
			userModel, err := auth.CreateUser(ctx, models.User{
				Username: payload.Username,
			}, payload.Token)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)

				return
			}

			user = *userModel

		} else {
			// an existing user is adding a new device
			var err error
			user, err = datastore.GetUserByID(ctx, claims.Subject)
			if err != nil {
				http.Error(w, "not found", http.StatusNotFound)

				return
			}
		}

		webauthnUser := repo.NewWebAuthnUser(
			ctx,
			middleware.L(ctx),
			datastore,
			user,
		)

		exclusions := []protocol.CredentialDescriptor{}
		for _, cred := range webauthnUser.WebAuthnCredentials() {
			exclusions = append(exclusions, cred.Descriptor())
		}

		options, session, err := web.BeginRegistration(webauthnUser,
			webauthn.WithExclusions(exclusions),
			webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementPreferred),
			webauthn.WithConveyancePreference(protocol.PreferIndirectAttestation),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		sessionID, err := datastore.SaveWebauthnSession(ctx, session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "registration_session",
			Value:    sessionID,
			Secure:   *cfg.SecureCookie,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		})

		jsonResponse(w, options, http.StatusOK)
	})
}

func finishRegistrationHandler(cfg config.Config, web *webauthn.WebAuthn, datastore *repo.Repo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		response, err := protocol.ParseCredentialCreationResponseBody(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		middleware.L(ctx).Infof("%+v", response)

		cookie := middleware.FindCookie("registration_session", r.Header)
		if cookie == nil {
			http.Error(w, "cookie not found", http.StatusBadRequest)
			return
		}

		session, err := datastore.GetWebauthnSession(ctx, cookie.Value)
		if err != nil {
			http.Error(w, "session not found: "+err.Error(), http.StatusNotFound)

			return
		}

		user, err := datastore.GetUserByID(ctx, string(session.UserID))
		if err != nil {
			http.Error(w, "user not found: "+err.Error(), http.StatusNotFound)

			return
		}

		webauthnUser := repo.NewWebAuthnUser(
			ctx,
			middleware.L(ctx),
			datastore,
			user,
		)

		cred, err := web.CreateCredential(webauthnUser, *session, response)
		if err != nil {
			http.Error(w, "failed to create credentials: "+err.Error(), http.StatusInternalServerError)

			return
		}

		ua := useragent.Parse(r.UserAgent())

		if err := datastore.AddWebauthnCred(ctx, user.ID, *cred, ua); err != nil {
			http.Error(w, "failed to create credentials: "+err.Error(), http.StatusInternalServerError)

			return
		}

		jsonResponse(w, "Success", http.StatusOK)
	})
}

func beginLoginHandler(cfg config.Config, web *webauthn.WebAuthn, datastore *repo.Repo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var (
			options *protocol.CredentialAssertion
			session *webauthn.SessionData
		)

		pathParts := strings.Split(r.URL.Path, "/")
		userNameOrEmail := pathParts[len(pathParts)-1]

		if userNameOrEmail != "" {
			user, err := datastore.GetUserByName(ctx, userNameOrEmail)
			if err != nil {
				if cfg.FeatureEnabled(config.FeatureLoginByMail) {
					var verified bool
					user, verified, err = datastore.GetUserByEMail(ctx, userNameOrEmail)

					if err == nil && !verified {
						http.Error(w, "e-mail address not verified", http.StatusPreconditionFailed)

						return
					}
				}
			}

			webauthnUser := repo.NewWebAuthnUser(
				ctx,
				middleware.L(ctx),
				datastore,
				user,
			)

			if err != nil {
				http.Error(w, "user not found", http.StatusNotFound)

				return
			}

			options, session, err = web.BeginLogin(webauthnUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}
		} else {
			var err error

			options, session, err = web.BeginDiscoverableLogin()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}
		}

		sessionID, err := datastore.SaveWebauthnSession(ctx, session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "login_session",
			Value:    sessionID,
			Secure:   *cfg.SecureCookie,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().Add(time.Minute * 5),
			Path:     "/",
		})

		jsonResponse(w, options, http.StatusOK)
	})
}

func finishLoginHandler(cfg config.Config, auth *auth.AuthService, web *webauthn.WebAuthn, datastore *repo.Repo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		response, err := protocol.ParseCredentialRequestResponseBody(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		requestedRedirect := r.URL.Query().Get("redirect")
		if requestedRedirect != "" {
			requestedRedirect, err = auth.HandleRequestedRedirect(ctx, requestedRedirect)
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

		session, err := datastore.GetWebauthnSession(ctx, cookie.Value)
		if err != nil {
			http.Error(w, "session not found", http.StatusNotFound)

			return
		}

		var user models.User
		getUserID := func(rawID, userHandle []byte) (webauthn.User, error) {
			var err error
			user, err = datastore.GetUserByID(ctx, string(userHandle))
			if err != nil {
				return nil, fmt.Errorf("user not found")
			}

			webauthnUser := repo.NewWebAuthnUser(
				ctx,
				middleware.L(ctx),
				datastore,
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

			_, err = web.ValidateLogin(webauthnUser, *session, response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)

				return
			}
		} else {
			_, err := web.ValidateDiscoverableLogin(getUserID, *session, response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)

				return
			}
		}

		roles, err := datastore.GetUserRoles(ctx, user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		refreshToken, refreshTokenID, err := auth.CreateSignedJWT(user, roles, "", cfg.RefreshTokenTTL.AsDuration(), jwt.ScopeRefresh)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		accessToken, _, err := auth.CreateSignedJWT(user, roles, refreshTokenID, cfg.AccessTokenTTL.AsDuration(), jwt.ScopeAccess)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		auth.AddRefreshTokenCookie(w.Header(), refreshToken, cfg.RefreshTokenTTL.AsDuration())
		auth.AddAccessTokenCookie(w.Header(), accessToken, cfg.AccessTokenTTL.AsDuration())

		userResponse := make(map[string]any)

		if requestedRedirect != "" {
			userResponse["redirectTo"] = requestedRedirect
		}

		jsonResponse(w, userResponse, http.StatusOK)
	})
}

func jsonResponse(w http.ResponseWriter, body any, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	enc.Encode(body)
}
