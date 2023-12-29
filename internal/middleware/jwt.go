package middleware

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/apis/pkg/server"
	"github.com/tierklinik-dobersberg/cis-idm/internal/common"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

func AuthenticateRequest(cfg config.Config, ds *repo.Queries, req *http.Request) (*jwt.Claims, bool, error) {
	ctx := req.Context()

	l := log.L(ctx)

	ips := server.RealIPFromContext(ctx)
	if len(ips) > 0 {
		l = l.WithField("clientIP", ips.String())
	}

	token := TokenFromContext(ctx)

	// first, try to parse the token as a JWT and if that worked, immediately
	// return the claims
	claims, tokenErr := jwt.ParseAndVerify([]byte(cfg.JWTSecret), token)

	if tokenErr == nil {
		l.Debugf("found valid JWT for user %s (name=%q)", claims.Subject, claims.Name)

		mail, err := ds.GetPrimaryEmailForUserByID(ctx, claims.Subject)
		if err == nil {
			claims.Email = mail.Address
		} else {
			if !errors.Is(err, sql.ErrNoRows) {
				l.Errorf("failed to get primary user mail: %s", err)
			} else {
				l.Debugf("user does not have a primary mail address configured")
			}
		}

		return claims, true, nil
	} else if token != "" {
		l.Infof("failed to parse JWT: %s", tokenErr)
	}

	// search for a forward-auth entry that might allow this request
	l.Debugf("searching for forward-auth entry for %s %s", req.Method, req.URL.String())

	fae, required, err := cfg.AuthRequiredForURL(req.Context(), req.Method, req.URL.String())
	if err != nil {
		return nil, false, fmt.Errorf("failed to get forward auth entry: %w", err)
	}

	// Handle the token error gracefully, expired tokens are rejected, malformed tokens
	// are allowed since it might represent a static token from the configuration file.
	if verr := new(gojwt.ValidationError); errors.As(tokenErr, &verr) {
		switch {
		case (verr.Errors&gojwt.ValidationErrorExpired) > 0 && required:
			return nil, !required, verr

		case (verr.Errors & gojwt.ValidationErrorMalformed) > 0:
			// empty on purpose, this error is expected for static tokens

			if token == "" {
				tokenErr = nil
			}

		default:
			l.Infof("unexpected JWT token error: (%T) %s", verr, verr)

			return nil, !required, verr
		}
	} else {
		l.Infof("unexpected token error: (%T) %s", tokenErr, tokenErr)

		// this is some other error so return it immediately
		return nil, !required, tokenErr
	}

	// There is no forward-auth config entry for this request
	// so just return the previous token error
	if fae == nil {
		l.Debug("no forward auth entry found")

		return nil, !required, tokenErr
	}

	l.Debug("checking forward-auth entry")

	subject, isAllowed, err := fae.Allowed(req)
	if err != nil {
		// just log the error here.
		l.Errorf("failed to check forward auth entry: %s", err)
	}

	if !isAllowed {
		l.Debug("request is not allowed by forward-auth entry")

		return nil, !required, nil
	}

	if subject == "" {
		l.Debugf("anonoumouse request authenticated by forward-auth")

		return nil, true, nil
	}

	l.Debugf("request authenticated by forward-auth: loading user by id %q", subject)

	claims = &jwt.Claims{}

	l = l.WithField("subject", subject)

	user, err := ds.GetUserByID(ctx, subject)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find user: %w", err)
	}

	common.EnsureDisplayName(&user)

	claims.Subject = user.ID
	claims.Name = user.Username
	claims.DisplayName = user.DisplayName

	mail, err := ds.GetPrimaryEmailForUserByID(ctx, claims.Subject)
	if err == nil {
		claims.Email = mail.Address
	} else {
		if !errors.Is(err, sql.ErrNoRows) {
			l.Errorf("failed to get primary user mail: %s", err)
		} else {
			l.Debugf("user does not have a primary mail address configured")
		}
	}

	// get user roles and append to the claims.

	userRoles, err := ds.GetRolesForUser(ctx, claims.Subject)
	if err == nil {
		claims.AppMetadata = &jwt.AppMetadata{
			Authorization: &jwt.Authorization{},
		}
		for _, r := range userRoles {
			claims.AppMetadata.Authorization.Roles = append(claims.AppMetadata.Authorization.Roles, r.ID)
		}
	} else {
		l.Errorf("failed to get user roles: %s", err)
	}

	return claims, true, nil
}

func NewJWTMiddleware(cfg config.Config, repo *repo.Queries, next http.Handler, skipVerifyFunc func(r *http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		header := r.Header.Get("Authentication")

		var (
			token       string
			claims      *jwt.Claims
			tokenSource string
		)
		if strings.HasPrefix(header, "Bearer ") {
			token = strings.Replace(header, "Bearer ", "", 1)
			tokenSource = "header"
		} else {
			// try to get the access token from a cookie
			cookie := FindCookie(cfg.AccessTokenCookieName, r.Header)
			if cookie != nil {
				token = cookie.Value
				tokenSource = "cookie"
			}
		}

		if token != "" {
			ctx = ContextWithToken(ctx, token)

			// add the new context to the request.
			r = r.WithContext(ctx)
		}

		ips := server.RealIPFromContext(ctx)

		// Fix the request URL by adding host/scheme
		r.URL.Host = r.Host
		r.URL.Scheme = "http"

		l := log.L(ctx).
			WithField("method", r.Method).
			WithField("host", r.URL.Host).
			WithField("path", r.URL.Path).
			WithField("clientIP", ips)

		// For debugging
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			l = l.WithField("xff", xff)
		}

		if r.TLS != nil {
			r.URL.Scheme = "https"
		} else if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
			r.URL.Scheme = proto
		}

		// Update the request logger
		r = r.WithContext(log.WithLogger(ctx, l))

		// try to authenticate the request
		var err error
		claims, _, err = AuthenticateRequest(cfg, repo, r)
		if err != nil {
			l.Errorf("failed to authenticate request: %s", err)
		}

		if skipVerifyFunc == nil || !skipVerifyFunc(r) {
			if err == nil && claims != nil && claims.ID != "" {
				var isRejected bool
				isRejected, err = repo.IsTokenRejected(ctx, claims.ID)

				if err == nil && !isRejected && claims.AppMetadata != nil && claims.AppMetadata.ParentTokenID != "" {
					isRejected, err = repo.IsTokenRejected(ctx, claims.AppMetadata.ParentTokenID)
				}

				if err == nil && isRejected {
					err = fmt.Errorf("token has been rejected")
				}
			}

			if err != nil {
				w.Header().Set("Content-Type", "application/json;encoding=utf-8")
				w.WriteHeader(http.StatusForbidden)

				blob, _ := json.Marshal(map[string]any{"code": connect.CodeUnauthenticated, "message": "invalid access token", "details": err.Error()})
				if _, err := w.Write(blob); err != nil {
					l.WithError(err).Errorf("failed to write response to client")
				}

				return
			}
		}

		if claims != nil {
			l.Debugf("adding claims for user %s (name=%q) to request context", claims.Subject, claims.Name)

			ctx = ContextWithClaims(ctx, claims)
			ctx = log.WithLogger(ctx, log.L(ctx).WithFields(logrus.Fields{
				"jwt:sub":         claims.Subject,
				"jwt:name":        claims.Name,
				"jwt:tokenSource": tokenSource,
			}))

			r = r.WithContext(ctx)
		} else {
			l.Debugf("no claims found, request is unauthenticated")
		}

		// call through to the next handler
		next.ServeHTTP(w, r)
	}
}
