package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/apis/pkg/spa"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
)

func NewForwardAuthHandler(providers *app.Providers) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		method := r.Header.Get("x-forwarded-method")
		u := &url.URL{
			Scheme: r.Header.Get("x-forwarded-proto"),
			Host:   r.Header.Get("x-forwarded-host"),
			Path:   r.Header.Get("x-forwarded-uri"),
		}
		requestURL := u.String()

		l := log.L(ctx).
			WithField("method", method).
			WithField("host", u.Host).
			WithField("path", u.Path)

		var redirectUrl = requestURL

		// Skip access checks for CORS preflight requests.
		if strings.ToUpper(method) == "OPTIONS" {
			w.WriteHeader(http.StatusOK)

			return
		}

		// if the request does not accept text/html it's likely a XHR or fetch request.
		// In this case, redirect the user to the referrer/origin rather thant the x-forwarded-uri
		if !strings.Contains(r.Header.Get("Accept"), "text/html") {
			redirectUrl = r.Referer()

			if _, err := url.Parse(redirectUrl); err != nil {
				l.Errorf("failed to parse redirect URL: %s", redirectUrl)
				redirectUrl = ""
			}

			if redirectUrl == "" {
				if origin, err := url.Parse(r.Header.Get("Origin")); err == nil {
					redirectUrl = origin.String()
				}
			}

			l.Debugf("redirect URL is set to %q", redirectUrl)
		}

		// auhenticate the request
		reqCopy := r.Clone(ctx)
		reqCopy.URL = u
		reqCopy.Method = method

		claims, isAllowed, err := middleware.AuthenticateRequest(providers.Config, providers.Datastore, reqCopy)
		if !isAllowed {
			var (
				verr = new(gojwt.ValidationError)
			)

			switch {
			case err == nil:
				l.Debugf("request not allowed")
				handleRedirect(w, r, providers.Config.LoginRedirectURL, redirectUrl)

			case errors.As(err, verr) && (verr.Errors&gojwt.ValidationErrorExpired) > 0:
				l.Debugf("request not allowed: JWT token expired")
				handleRedirect(w, r, providers.Config.RefreshRedirectURL, redirectUrl)

			default:
				l.Debugf("request not allowed: %s", err)
				handleRedirect(w, r, "", "")
			}

			return
		}

		if err != nil {
			l.Errorf("request not allowed due to errors: %s", err)
			handleRedirect(w, r, providers.Config.LoginRedirectURL, redirectUrl)

			return
		}

		if claims != nil {
			w.Header().Add("X-Remote-User-ID", claims.Subject)
			w.Header().Add("X-Remote-User", claims.Name)
			w.Header().Add("X-Remote-Avatar-URL", fmt.Sprintf("%s/avatar/%s", providers.Config.PublicURL, claims.Subject))

			if claims.DisplayName != "" {
				w.Header().Add("X-Remote-User-Display-Name", claims.DisplayName)
			}

			if claims.Email != "" {
				w.Header().Add("X-Remote-Mail", claims.Email)
			}

			if claims.AppMetadata != nil && claims.AppMetadata.Authorization != nil {
				for _, r := range claims.AppMetadata.Authorization.Roles {
					w.Header().Add("X-Remote-Role", r)
				}
			}

			l.Infof("request by user %s (name=%q) is allowed", claims.Subject, claims.Name)
		} else {
			l.Infof("anonymous request is allowed")
		}

		w.WriteHeader(http.StatusOK)
	})
}

func handleRedirect(w http.ResponseWriter, r *http.Request, baseUrl string, url string) {
	spa.SetCORSHeaders(w, r)

	if url == "" || baseUrl == "" {
		http.Error(w, "not allowed", http.StatusForbidden)

		return
	}

	encodedRedirectURL := base64.URLEncoding.EncodeToString([]byte(url))
	targetUrl := fmt.Sprintf(baseUrl, encodedRedirectURL)

	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		http.Redirect(w, r, targetUrl, http.StatusFound)

		return
	}

	blob, _ := json.Marshal(map[string]any{
		"location": targetUrl,
	})

	w.WriteHeader(http.StatusForbidden)
	w.Write(blob)
}
