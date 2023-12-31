package webauthn

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gofrs/uuid"
	"github.com/mileusna/useragent"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
)

func (svc *Service) BeginRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := log.L(ctx)

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
		log.L(ctx),
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

	sessionID, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := svc.Cache.PutKeyTTL(ctx, sessionID.String(), session, session.Expires.Sub(time.Now())); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "registration_session",
		Value:    sessionID.String(),
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

	log.L(ctx).Infof("%+v", response)

	cookie := middleware.FindCookie("registration_session", r.Header)
	if cookie == nil {
		http.Error(w, "cookie not found", http.StatusBadRequest)
		return
	}

	var session webauthn.SessionData
	if err := svc.Cache.GetAndDeleteKey(ctx, cookie.Value, &session); err != nil {
		log.L(ctx).Errorf("failed to find webauthn registration session: %s", err)
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
		log.L(ctx),
		svc.Datastore,
		user,
	)

	cred, err := svc.web.CreateCredential(webauthnUser, session, response)
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
