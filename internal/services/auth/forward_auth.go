package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/apis/pkg/spa"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func parseXForwardedForHeader(r *http.Request) []net.IP {
	h := r.Header.Get("X-Forwarded-For")
	if h == "" {
		return nil
	}

	result := make([]net.IP, 0)

	ips := strings.Split(h, ",")
	for _, ip := range ips {
		i := net.ParseIP(strings.TrimSpace(ip))
		if i == nil {
			log.L(r.Context()).Errorf("received invalid x-forwarded-for header: %s", h)

			return nil
		}

		result = append(result, i)
	}

	return result
}

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

		ips := parseXForwardedForHeader(r)
		if len(ips) > 0 {
			l = l.WithField("clientIP", ips[0].String())
		}

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
		}

		fae, required, err := providers.Config.AuthRequiredForURL(method, requestURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token := middleware.TokenFromContext(ctx)

		// if authentication is required but there's not even a token, redirect to
		// the login URL or deny access
		if token == "" && required {
			handleRedirect(w, r, providers.Config.LoginRedirectURL, redirectUrl)

			return
		}

		var (
			userId              string
			primaryMail         string
			primaryMailVerified bool
			displayName         string
			roles               []string
		)

		// first, try to parse the token as a JWT
		claims, tokenErr := jwt.ParseAndVerify([]byte(providers.Config.JWTSecret), token)

		if tokenErr == nil {
			// we have a valid JWT here so load the user, primary mail
			userId = claims.Subject
			user, err := providers.Datastore.GetUserByID(ctx, claims.Subject)
			if err != nil {
				l.Errorf("failed to find user by ID: %s", claims.Subject)

				http.Error(w, "access token subject not found", http.StatusForbidden)

				return
			}

			mail, err := providers.Datastore.GetUserPrimaryMail(ctx, claims.Subject)
			if err == nil {
				primaryMail = mail.Address
				primaryMailVerified = mail.Verified
			} else {
				if !errors.Is(err, stmts.ErrNoResults) {
					l.Errorf("failed to get primary user mail: %s", err)
				} else {
					l.Debugf("user does not have a primary mail address configured")
				}
			}

			displayName = user.DisplayName
			if displayName == "" {
				displayName = user.Username
			}

			if claims.AppMetadata != nil && claims.AppMetadata.Authorization != nil {
				roles = claims.AppMetadata.Authorization.Roles
			}
		} else {
			// check if the token has been expired, and if, redirect the user to
			// the RefreshRedirectURL.
			if verr := new(gojwt.ValidationError); errors.As(tokenErr, verr) {
				switch {
				case (verr.Errors&gojwt.ValidationErrorExpired) > 0 && required:
					handleRedirect(w, r, providers.Config.RefreshRedirectURL, redirectUrl)

					return

				case (verr.Errors&gojwt.ValidationErrorMalformed) > 0 && fae != nil:
					// this seems to don't event be a JWT, so try to verify using static tokens
					// from the ForwardAuthEntry.
					for _, staticToken := range fae.Tokens {
						if staticToken.Tokens == token {
							userId = staticToken.SubjectID
							roles = staticToken.Roles

							break
						}
					}

				default:
					// forbidden, there's something wrong with the JWT
					handleRedirect(w, r, "", "")

					return
				}
			}
		}

		if required && userId == "" {
			// auth is required but we failed to authenticate the request.
			// also, there was a token present so the user tried to authenticate.
			// don't redirect to the login-page here (if the token would just have been expired,
			// the user would have been redirected to the refresh-page already)
			handleRedirect(w, r, "", "")

			return
		}

		// Add forward-auth headers
		if displayName != "" {
			w.Header().Add("X-Remote-User", displayName)
		}
		if userId != "" {
			w.Header().Add("X-Remote-User-ID", userId)
			w.Header().Add("X-Remote-Avatar-URL", fmt.Sprintf("%s/avatar/%s", providers.Config.PublicURL, userId))
		}
		if primaryMail != "" {
			w.Header().Add("X-Remote-Mail", primaryMail)
			w.Header().Add("X-Remote-Mail-Verified", strconv.FormatBool(primaryMailVerified))
		}
		if len(roles) > 0 {
			for _, role := range roles {
				w.Header().Add("X-Remote-Role", role)
			}
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
