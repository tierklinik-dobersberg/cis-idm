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
			log.L(ctx).Errorf("failed to parse X-Forwarded-URI %q: %s", r.Header.Get("x-forwarded-uri"), err)
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
		if s := providers.Config.ForwardAuth.AllowCORSPreflight; s == nil || *s {
			if strings.ToUpper(method) == "OPTIONS" && r.Header.Get("Origin") != "" && r.Header.Get("Access-Control-Request-Method") != "" {
				w.WriteHeader(http.StatusOK)

				return
			}
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
		reqCopy.Host = u.Host

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

		query := providers.Config.ForwardAuth.RegoQuery
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

		l = l.WithField("policyResult", result)

		var isAllowed bool
		if providers.Config.ForwardAuth.Default == "deny" {
			isAllowed = result.Allow
		} else {
			isAllowed = !result.Deny
		}

		// evaluate the result
		if !isAllowed {
			// The request has been denied by policy, now figure out how to reply:

			if authErr != nil {
				l = l.WithField("token_error", authErr)
			}

			l.Infof("request has been denied by policy")

			switch {
			// If a status code has been assigned than we directly reply with
			// this code. This is useful if a request should be denied even if
			// it is authenticated.
			case result.StatusCode > 0:
				w.WriteHeader(result.StatusCode)
				if _, err := w.Write([]byte(result.ResponseBody)); err != nil {
					l.Errorf("failed to write response body: %s", err)
				}

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

		if result.AssignSubject != "" {
			l.Infof("loading subject overwrite")
			input.Subject, err = policy.NewSubjectInput(ctx, providers.Datastore, providers.Config.PermissionTree(), result.AssignSubject, "", "")
			if err != nil {
				l.Errorf("failed to overwrite request subject: %s", err)
				handleRedirect(w, r, "", "")
				return
			}
		}

		// If we got an authenticated subject, add those headers as well
		if sub := input.Subject; sub != nil {
			fwCfg := providers.Config.ForwardAuth

			if h := fwCfg.UserIDHeader; *h != "" {
				w.Header().Add(*h, sub.ID)
			}

			if h := fwCfg.UsernameHeader; *h != "" {
				w.Header().Add(*h, sub.Username)
			}

			if h := fwCfg.AvatarURLHeader; *h != "" {
				w.Header().Add(*h, fmt.Sprintf("%s/avatar/%s", providers.Config.UserInterface.PublicURL, sub.ID))
			}

			if h := fwCfg.DisplayNameHeader; *h != "" {
				if sub.DisplayName != "" {
					w.Header().Add(*h, sub.DisplayName)
				}
			}

			if h := fwCfg.MailHeader; *h != "" {
				if sub.Email != "" {
					w.Header().Add(*h, sub.Email)
				}
			}

			if h := fwCfg.RoleHeader; *h != "" {
				for _, r := range sub.Roles {
					w.Header().Add(*h, r.ID)
				}
			}

			if h := fwCfg.ResolvedPermissionHeader; *h != "" {
				// all all permissions from all roles to the headers.
				for _, p := range sub.Permissions {
					w.Header().Add(*h, p)
				}
			}

			l.Infof("request by user %s (name=%q) is allowed", sub.ID, sub.Username)
		} else {
			l.Infof("anonymous request is allowed")
		}

		w.WriteHeader(http.StatusOK)
	})
}

func handleRedirect(w http.ResponseWriter, r *http.Request, baseUrl string, url string) {
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
