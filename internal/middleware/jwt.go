package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
)

type TokenRejector interface {
	IsTokenRejected(context.Context, string) (bool, error)
}

func NewJWTMiddleware(cfg config.Config, repo TokenRejector, next http.Handler, skipVerifyFunc func(r *http.Request) bool) http.HandlerFunc {
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

			var err error
			claims, err = jwt.ParseAndVerify([]byte(cfg.JWTSecret), token)

			if skipVerifyFunc != nil && !skipVerifyFunc(r) {
				if err == nil {
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
						log.L(ctx).WithError(err).Errorf("failed to write response to client")
					}

					return
				}
			}

			if claims != nil {
				ctx = ContextWithClaims(ctx, claims)
				ctx = log.WithLogger(ctx, log.L(ctx).WithFields(logrus.Fields{
					"jwt:sub":         claims.Subject,
					"jwt:name":        claims.Name,
					"jwt:tokenSource": tokenSource,
				}))
			}
		}

		// add the new context to the request.
		r = r.WithContext(ctx)

		// call through to the next handler
		next.ServeHTTP(w, r)
	}
}
