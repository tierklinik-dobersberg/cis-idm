package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

func NewForwardAuthHandler(cfg config.Config, repo *repo.Repo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		method := r.Header.Get("x-forwarded-method")
		requestURL := (&url.URL{
			Scheme: r.Header.Get("x-forwarded-proto"),
			Host:   r.Header.Get("x-forwarded-host"),
			Path:   r.Header.Get("x-forwarded-uri"),
		}).String()

		encodedRequestURL := base64.URLEncoding.EncodeToString([]byte(requestURL))

		// FIXME(ppacher): make sure the JWTMiddleware does not return an error if the token is invalid
		l := middleware.L(ctx)

		claims := middleware.ClaimsFromContext(ctx)
		token := middleware.TokenFromContext(ctx)

		// parse and verify the token here so we can react to token-expired errors.
		var tokenErr error
		if token != "" {
			_, tokenErr = jwt.ParseAndVerify([]byte(cfg.JWTSecret), token)
		}

		l.WithField("has_token", token != "").
			WithField("has_claims", claims != nil).
			WithField("token_error", tokenErr).
			Infof("received forward_auth request for %s", requestURL)

		// check if authentication is required for this URL
		if required, err := cfg.AuthRequiredForURL(method, requestURL); required {

			// if there is no access token at all so redirect the user to the login
			// screen
			if token == "" {
				if cfg.LoginRedirectURL != "" {
					url := fmt.Sprintf(cfg.LoginRedirectURL, encodedRequestURL)

					http.Redirect(w, r, url, http.StatusFound)
					return
				}
			}

			// check if the token has been expired, and if, redirect the user to
			// the RefreshRedirectURL.
			if verr := new(gojwt.ValidationError); errors.As(tokenErr, verr) && (verr.Errors&gojwt.ValidationErrorExpired) > 0 {
				if cfg.RefreshRedirectURL != "" {
					url := fmt.Sprintf(cfg.RefreshRedirectURL, encodedRequestURL)

					http.Redirect(w, r, url, http.StatusFound)

					return
				}
			}

			// Authentication is required but it seems like we don't have any
			// valid JWT claims and no login/refresh URL defined. We can just respond
			// with StatusForbidden now
			if claims == nil || tokenErr != nil {
				// otherwise access is forbidden
				errorString := "forbidden"
				if token == "" {
					errorString = "no access token"
				} else if claims == nil {
					errorString = "no token claims found"
				} else if tokenErr != nil {
					errorString = tokenErr.Error()
				}

				http.Error(w, errorString, http.StatusForbidden)

				if tokenErr != nil {
					l.WithError(tokenErr).Errorf("invalid access token")
				}

				return
			}

			// we have a valid jWT access token so we can fallthrough and continue
		} else if err != nil {
			middleware.L(ctx).Errorf("failed to check request to %q for authentication requirements: %s", requestURL, err)

			http.Error(w, "internal server error: invalid forward auth configuration", http.StatusInternalServerError)

			return
		}

		// if we have valid user claims then load the user and it's primary e-mail address from
		// the repository and add some headers for the upstream server.
		if claims != nil {
			user, err := repo.GetUserByID(ctx, claims.Subject)
			if err != nil {
				middleware.L(ctx).Errorf("failed to find user by ID")

				http.Error(w, "access token subject not found", http.StatusForbidden)

				return
			}

			mail, err := repo.GetUserPrimaryMail(ctx, claims.Subject)
			if err != nil {
				middleware.L(ctx).Errorf("failed to get primary user mail: %s", err)
			} else {
				w.Header().Add("X-Remote-Mail", mail.Address)
				w.Header().Add("X-Remote-Mail-Verified", strconv.FormatBool(mail.Verified))
			}

			displayName := user.DisplayName
			if displayName == "" {
				displayName = user.Username
			}

			w.Header().Add("X-Remote-User", displayName)
			w.Header().Add("X-Remote-User-ID", user.ID)
			w.Header().Add("X-Remote-Avatar-URL", fmt.Sprintf("%s/avatar/%s", cfg.PublicURL, user.ID))
		}

		middleware.L(ctx).WithFields(logrus.Fields{
			"url": requestURL,
		}).Info("handled forward-auth request")

		w.WriteHeader(http.StatusOK)
	})
}
