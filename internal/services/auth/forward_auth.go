package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/apis/pkg/server"
	"github.com/tierklinik-dobersberg/apis/pkg/spa"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/policy"
)

func NewForwardAuthHandler(providers *app.Providers) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		parsedForwardedUri, err := url.ParseRequestURI(r.Header.Get("x-forwarded-uri"))
		if err != nil {
			log.L(ctx).Errorf("failed to parse X-Forwareded-URI %q: %s", r.Header.Get("x-forwarded-uri"), err)
		}

		method := r.Header.Get("x-forwarded-method")
		u := &url.URL{
			Scheme:   r.Header.Get("x-forwarded-proto"),
			Host:     r.Header.Get("x-forwarded-host"),
			Path:     parsedForwardedUri.Path,
			RawPath:  parsedForwardedUri.RawPath,
			RawQuery: parsedForwardedUri.RawQuery,
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

		// try to authenticate the request.
		claims, authErr := middleware.AuthenticateRequest(providers.Config, providers.Datastore, reqCopy)

		// prepare the input for the rego policy query
		input := ForwardAuthInput{
			Method:   reqCopy.Method,
			Path:     reqCopy.URL.Path,
			Query:    reqCopy.URL.Query(),
			Host:     reqCopy.Host,
			Headers:  reqCopy.Header,
			ClientIP: server.RealIPFromContext(ctx).String(),
		}

		// If we got valid JWT claims, resolve the SubjectInput
		if claims != nil {
			kind := jwt.LoginKindInvalid
			if claims.AppMetadata != nil {
				kind = claims.AppMetadata.LoginKind
			}

			input.Subject, err = policy.NewSubjectInput(ctx, providers.Datastore, providers.Config.PermissionTree(), claims.Subject, kind, claims.ID)
			if err != nil {
				l.Errorf("failed to resolve subject input: %s", err)

				// clear out the subject and let rego policies still evaluate the request.
				input.Subject = nil
			}
		}

		// Execute rego policies to find a decision
		var result ForwardAuthPolicyResult

		query := providers.Config.PolicyConfig.ForwardAuthQuery
		if err := providers.PolicyEngine.QueryOne(ctx, query, input, &result); err != nil {
			l.Errorf("failed to evaluate rego policies: %s; request will be denied", err)

			handleRedirect(w, r, "", "")

			return
		}

		// Regardless of if the request is permitted, add all headers from
		// the policy to the response
		if len(result.Headers) > 0 {
			for key, values := range result.Headers {
				for _, val := range values {
					w.Header().Add(key, val)
				}
			}
		}

		// evalute the result
		if !result.Allow {
			// The request has been denied by policy, now figure out how to reply:

			if authErr != nil {
				l = l.WithField("token_error", authErr)
			}

			if result.StatusCode > 0 {
				l = l.WithField("status_code", result.StatusCode)
			}

			l.Infof("request has been denied by policy")

			switch {
			// If a status code has been assigned than we directly reply with
			// this code. This is useful if a request should be denied even if
			// it is authenticated.
			case result.StatusCode > 0:
				w.WriteHeader(result.StatusCode)

			// If there wasn't even a token or the token has been rejected,
			// redirect to the login page.
			case errors.Is(authErr, middleware.ErrNoToken),
				errors.Is(authErr, middleware.ErrTokenRejected):

				handleRedirect(w, r, providers.Config.UserInterface.LoginRedirectURL, redirectUrl)

			// If the token has been expired, redirect to the refresh token page
			case errors.Is(authErr, middleware.ErrTokenExpired):
				handleRedirect(w, r, providers.Config.UserInterface.RefreshRedirectURL, redirectUrl)

			// We got a valid token but our rego policies denied the request. Respond without any
			// redirection
			case authErr == nil:

				handleRedirect(w, r, "", "")

			// request was denied by rego policies and we do have some invalid token at hand.
			// redirect the user to the login page.
			default:
				handleRedirect(w, r, providers.Config.UserInterface.LoginRedirectURL, redirectUrl)
			}

			return
		}

		l.Infof("request has been allowed by policy")

		// If we got an authenticated subject, add those headers as well
		// TODO(ppacher): make the default header names configurable
		if sub := input.Subject; sub != nil {
			w.Header().Add("X-Remote-User-ID", sub.ID)
			w.Header().Add("X-Remote-User", sub.Username)
			w.Header().Add("X-Remote-Avatar-URL", fmt.Sprintf("%s/avatar/%s", providers.Config.UserInterface.PublicURL, sub.ID))

			if claims.DisplayName != "" {
				w.Header().Add("X-Remote-User-Display-Name", sub.DisplayName)
			}

			if claims.Email != "" {
				w.Header().Add("X-Remote-Mail", sub.Email)
			}

			for _, r := range sub.Roles {
				w.Header().Add("X-Remote-Role", r.ID)
			}

			// all all permissions from all roles to the headers.
			for _, p := range sub.Permissions {
				w.Header().Add("X-Remote-Permission", p)
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
