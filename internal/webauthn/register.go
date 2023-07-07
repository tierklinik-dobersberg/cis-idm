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
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/auth"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
)

type Service struct {
	*app.Providers

	authService *auth.AuthService

	web *webauthn.WebAuthn
}

func New(providers *app.Providers, authService *auth.AuthService) (http.Handler, error) {
	mux := http.NewServeMux()

	wconfig := &webauthn.Config{
		RPDisplayName: providers.Config.SiteName,
		RPID:          providers.Config.Domain,
		RPOrigins: []string{
			providers.Config.PublicURL,
		},
	}

	w, err := webauthn.New(wconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create webauthn instance: %w", err)
	}

	instance := &Service{
		authService: authService,
		Providers:   providers,
		web:         w,
	}

	mux.Handle("/registration/begin", http.HandlerFunc(instance.BeginRegistrationHandler))
	mux.Handle("/registration/finish", http.HandlerFunc(instance.FinishRegistrationHandler))
	mux.Handle("/login/begin/", http.HandlerFunc(instance.BeginLoginHandler))
	mux.Handle("/login/finish", http.HandlerFunc(instance.FinishLoginHandler))

	return mux, nil
}

func (svc *Service) BeginRegistrationHandler(w http.ResponseWriter, r *http.Request) {
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
		userModel, err := svc.authService.CreateUser(ctx, models.User{
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
		user, err = svc.Datastore.GetUserByID(ctx, claims.Subject)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)

			return
		}
	}

	webauthnUser := repo.NewWebAuthnUser(
		ctx,
		middleware.L(ctx),
		svc.Datastore,
		user,
	)

	exclusions := []protocol.CredentialDescriptor{}
	for _, cred := range webauthnUser.WebAuthnCredentials() {
		exclusions = append(exclusions, cred.Descriptor())
	}

	options, session, err := svc.web.BeginRegistration(webauthnUser,
		webauthn.WithExclusions(exclusions),
		webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementPreferred),
		webauthn.WithConveyancePreference(protocol.PreferIndirectAttestation),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	sessionID, err := svc.Datastore.SaveWebauthnSession(ctx, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "registration_session",
		Value:    sessionID,
		Secure:   *svc.Config.SecureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	jsonResponse(w, options, http.StatusOK)
}

func (svc *Service) FinishRegistrationHandler(w http.ResponseWriter, r *http.Request) {
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

	session, err := svc.Datastore.GetWebauthnSession(ctx, cookie.Value)
	if err != nil {
		http.Error(w, "session not found: "+err.Error(), http.StatusNotFound)

		return
	}

	user, err := svc.Datastore.GetUserByID(ctx, string(session.UserID))
	if err != nil {
		http.Error(w, "user not found: "+err.Error(), http.StatusNotFound)

		return
	}

	webauthnUser := repo.NewWebAuthnUser(
		ctx,
		middleware.L(ctx),
		svc.Datastore,
		user,
	)

	cred, err := svc.web.CreateCredential(webauthnUser, *session, response)
	if err != nil {
		http.Error(w, "failed to create credentials: "+err.Error(), http.StatusInternalServerError)

		return
	}

	ua := useragent.Parse(r.UserAgent())

	if err := svc.Datastore.AddWebauthnCred(ctx, user.ID, *cred, ua); err != nil {
		http.Error(w, "failed to create credentials: "+err.Error(), http.StatusInternalServerError)

		return
	}

	jsonResponse(w, "Success", http.StatusOK)
}

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
			middleware.L(ctx),
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

	sessionID, err := svc.Datastore.SaveWebauthnSession(ctx, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "login_session",
		Value:    sessionID,
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

	session, err := svc.Datastore.GetWebauthnSession(ctx, cookie.Value)
	if err != nil {
		http.Error(w, "session not found", http.StatusNotFound)

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
			middleware.L(ctx),
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

		_, err = svc.web.ValidateLogin(webauthnUser, *session, response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)

			return
		}
	} else {
		_, err := svc.web.ValidateDiscoverableLogin(getUserID, *session, response)
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

func jsonResponse(w http.ResponseWriter, body any, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	enc.Encode(body)
}
