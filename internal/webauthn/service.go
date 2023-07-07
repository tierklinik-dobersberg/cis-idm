package webauthn

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/auth"
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

func jsonResponse(w http.ResponseWriter, body any, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	enc.Encode(body)
}
