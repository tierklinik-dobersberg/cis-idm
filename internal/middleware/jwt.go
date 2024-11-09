package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/apis/pkg/server"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

var (
	ErrNoToken         = errors.New("no authentication token")
	ErrInvalidAPIToken = errors.New("invalid API token")
	ErrTokenRejected   = errors.New("authentication token has been rejected")
	ErrTokenExpired    = errors.New("token has expired")
)

const APITokenPrefix = "it."

func AuthenticateRequest(cfg config.Config, ds *repo.Queries, req *http.Request) (*jwt.Claims, error) {
	ctx := req.Context()

	token := TokenFromContext(ctx)

	if token == "" {
		return nil, ErrNoToken
	}

	// check if this is an API token
	if strings.HasPrefix(token, APITokenPrefix) {
		res, err := ds.GetUserForAPIToken(ctx, token)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrInvalidAPIToken
			}

			return nil, err
		}

		// make sure that token is stil valid
		if res.UserApiToken.ExpiresAt.Valid && res.UserApiToken.ExpiresAt.Time.After(time.Now()) {
			return nil, ErrTokenExpired
		}

		userRoles, err := ds.GetRolesForUser(ctx, res.User.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to query user roles: %w", err)
		}

		roleIds := make([]string, len(userRoles))
		for rIdx, r := range userRoles {
			roleIds[rIdx] = r.ID
		}

		// construct claims for the user
		claims := jwt.Claims{
			ID:          res.UserApiToken.ID,
			IssuedAt:    time.Now().Unix(),
			NotBefore:   time.Now().Unix(),
			Subject:     res.User.ID,
			Name:        res.User.Username,
			DisplayName: res.User.DisplayName,
			Scopes: []jwt.Scope{
				jwt.ScopeAccess,
			},
			AppMetadata: &jwt.AppMetadata{
				TokenVersion: "1",
				Authorization: &jwt.Authorization{
					Roles: roleIds,
				},
				LoginKind: jwt.LoginKindAPI,
			},
		}

		if res.UserApiToken.ExpiresAt.Valid {
			claims.ExpiresAt = res.UserApiToken.ExpiresAt.Time.Unix()
		}

		return &claims, nil
	}

	// first, try to parse the token as a JWT and if that worked, immediately
	// return the claims
	claims, tokenErr := jwt.ParseAndVerify([]byte(cfg.JWT.Secret), token)

	// immediately abort if the JWT is invalid
	if tokenErr != nil {
		var verr *gojwt.ValidationError

		if errors.As(tokenErr, &verr) {
			switch {
			case verr.Errors&gojwt.ValidationErrorExpired > 0:
				return nil, ErrTokenExpired
			}
		}

		return nil, tokenErr
	}

	isRejected, err := ds.IsTokenRejected(ctx, claims.ID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check if token has been rejected: %w", err)
	}

	if !isRejected && claims.AppMetadata != nil && claims.AppMetadata.ParentTokenID != "" {
		isRejected, err = ds.IsTokenRejected(ctx, claims.AppMetadata.ParentTokenID)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to check if parent token has been rejected: %w", err)
		}
	}

	if isRejected {
		return nil, ErrTokenRejected
	}

	return claims, nil
}

func NewJWTMiddleware(cfg config.Config, repo *repo.Queries, next http.Handler, skipVerifyFunc func(r *http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		header := r.Header.Get("Authorization")

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
			cookie := FindCookie(cfg.JWT.AccessTokenCookieName, r.Header)
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

		// check if we should skip token verification on this endpoint
		if skipVerifyFunc != nil && skipVerifyFunc(r) {
			next.ServeHTTP(w, r)

			return
		}

		// try to authenticate the request
		claims, err := AuthenticateRequest(cfg, repo, r)

		if err != nil {
			l.Errorf("failed to authenticate request: %s", err)
		}

		if claims != nil {
			l.Debugf("adding claims for user %s (name=%q) to request context", claims.Subject, claims.Name)

			var loginKind jwt.LoginKind

			if claims.AppMetadata != nil {
				loginKind = claims.AppMetadata.LoginKind
			}

			ctx = ContextWithClaims(ctx, claims)
			ctx = log.WithLogger(ctx, log.L(ctx).WithFields(logrus.Fields{
				"jwt:sub":         claims.Subject,
				"jwt:name":        claims.Name,
				"jwt:tokenSource": tokenSource,
				"jwt:kind":        loginKind,
			}))

			r = r.WithContext(ctx)
		} else {
			l.Debug("no claims found, request is unauthenticated")
		}

		// call through to the next handler
		next.ServeHTTP(w, r)
	}
}
