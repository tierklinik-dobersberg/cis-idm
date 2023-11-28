package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/bufbuild/connect-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/apis/pkg/cors"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/apis/pkg/privacy"
	"github.com/tierklinik-dobersberg/apis/pkg/server"
	"github.com/tierklinik-dobersberg/apis/pkg/spa"
	"github.com/tierklinik-dobersberg/apis/pkg/validator"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/services/auth"
	"github.com/tierklinik-dobersberg/cis-idm/internal/services/notify"
	"github.com/tierklinik-dobersberg/cis-idm/internal/services/roles"
	"github.com/tierklinik-dobersberg/cis-idm/internal/services/selfservice"
	"github.com/tierklinik-dobersberg/cis-idm/internal/services/users"
	"github.com/tierklinik-dobersberg/cis-idm/internal/webauthn"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

//go:embed static/ui
var static embed.FS

func getProtoRegistry() (*protoregistry.Files, error) {
	reg := new(protoregistry.Files)
	for _, file := range []protoreflect.FileDescriptor{
		idmv1.File_tkd_idm_v1_auth_service_proto,
		idmv1.File_tkd_idm_v1_self_service_proto,
		idmv1.File_tkd_idm_v1_user_proto,
		idmv1.File_tkd_idm_v1_user_service_proto,
		idmv1.File_tkd_idm_v1_role_service_proto,
		idmv1.File_tkd_idm_v1_notify_service_proto,
	} {
		if err := reg.RegisterFile(file); err != nil {
			return nil, fmt.Errorf("failed to register %s at protoregistry: %w", file.Name(), err)
		}
	}

	return reg, nil
}

func getStaticFilesHandler(path string) (http.Handler, error) {
	if path == "" {
		webapp, err := fs.Sub(static, "static/ui")
		if err != nil {
			return nil, err
		}
		return spa.ServeSPA(http.FS(webapp), "index.html"), nil
	}

	if strings.HasPrefix(path, "http") {
		remote, err := url.Parse(path)
		if err != nil {
			return nil, err
		}

		handler := func(p *httputil.ReverseProxy) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Host = remote.Host
				p.ServeHTTP(w, r)
			})
		}

		return handler(httputil.NewSingleHostReverseProxy(remote)), nil
	}

	return spa.ServeSPA(http.Dir(path), "index.html"), nil
}

func setupPublicServer(providers *app.Providers) (*http.Server, error) {
	// prepare middlewares and interceptors
	loggingInterceptor := log.NewLoggingInterceptor()
	authInterceptor := middleware.NewAuthInterceptor(providers.ProtoRegistry)
	validatorInterceptor := validator.NewInterceptor(providers.Validator)

	privacyInterceptor := privacy.NewFilterInterceptor(privacy.SubjectResolverFunc(func(ctx context.Context, ar connect.AnyRequest) (string, []string, error) {
		claims := middleware.ClaimsFromContext(ctx)
		if claims == nil {
			return "", nil, nil
		}

		return claims.Subject, claims.AppMetadata.Authorization.Roles, nil
	}))

	errorInterceptor := middleware.NewErrorInterceptor()

	interceptors := connect.WithInterceptors(
		loggingInterceptor,
		authInterceptor,
		validatorInterceptor,
		privacyInterceptor,
		errorInterceptor,
	)

	// create a new servemux to handle our routes.
	serveMux := http.NewServeMux()

	// Setup and serve the AuthService
	authService := auth.NewService(providers)
	path, handler := idmv1connect.NewAuthServiceHandler(
		authService,
		interceptors,
	)
	serveMux.Handle(path, handler)

	// Setup and serve the Self-Service
	selfserviceService := selfservice.NewService(providers)

	path, handler = idmv1connect.NewSelfServiceServiceHandler(
		selfserviceService,
		interceptors,
	)
	serveMux.Handle(path, handler)

	// Setup and serve the User service.
	userService, err := users.NewService(providers)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize users service: %w", err)
	}
	path, handler = idmv1connect.NewUserServiceHandler(
		userService,
		interceptors,
	)
	serveMux.Handle(path, handler)

	roleService := roles.NewService(providers)
	path, handler = idmv1connect.NewRoleServiceHandler(
		roleService,
		interceptors,
	)
	serveMux.Handle(path, handler)

	// Notify service
	notifyService := notify.New(providers)
	path, handler = idmv1connect.NewNotifyServiceHandler(
		notifyService,
		interceptors,
	)
	serveMux.Handle(path, handler)

	// Serve basic configuration for the UI on /config.json
	serveMux.Handle("/config.json", config.NewConfigHandler(providers.Config))

	// Setup CORS middleware
	corsOpts := cors.Config{
		AllowedOrigins: append([]string{
			providers.Config.PublicURL,
			fmt.Sprintf("http://%s", providers.Config.Domain),
			fmt.Sprintf("https://%s", providers.Config.Domain),
		}, providers.Config.AllowedOrigins...),
		AllowCredentials: true,
	}

	// Get a static file handler.
	// This will either return a handler for the embed.FS, a local directory using http.Dir
	// or a reverse proxy to some other service.
	staticFilesHandler, err := getStaticFilesHandler(providers.Config.StaticFiles)
	if err != nil {
		return nil, err
	}

	serveMux.Handle("/", staticFilesHandler)

	// Setup the avatar handler. This does not use connect-go since we need
	// to response with either HTTP redirects (if the user avatar is a URL)
	// or with plain bytes and an approriate content-type if the user avatar
	// is a dataurl.
	serveMux.Handle("/avatar/", users.NewAvatarHandler(providers))

	// setup the webauthn handlers for registration and login.
	// TODO(ppacher): migrate those to connect-go/protobuf style endpoints
	// as the browser does not actually care about how this is implemented.
	webauthnHandler, err := webauthn.New(providers, authService)
	if err != nil {
		return nil, err
	}
	serveMux.Handle("/webauthn/", http.StripPrefix("/webauthn", webauthnHandler))

	// Setup the forward auth handler
	serveMux.Handle("/validate", auth.NewForwardAuthHandler(providers))

	// If we're in debug mode, add some debug endpoints
	if os.Getenv("DEBUG") != "" {
		serveMux.Handle("/debug/cpu", http.HandlerFunc(CPUProfileHandler))
	}

	// finally, return a http.Server that uses h2c for HTTP/2 support and
	// wrap the finnal handler in CORS and a JWT middleware.
	return server.CreateWithOptions(
		providers.Config.PublicListenAddr,

		middleware.NewJWTMiddleware(providers.Config, providers.Datastore, serveMux, func(r *http.Request) bool {
			// Skip JWT token verification for the /validate endpoint as
			// the ForwardAuthHanlder will take care of this on it's own due to special
			// handling of rejected or expired tokens.
			if r.URL.Path == "/validate" || strings.HasPrefix(r.URL.Path, "/webauthn") {
				return true
			}

			return false
		}),

		server.WithCORS(corsOpts),
		server.WithTrustedProxies(providers.Config.TrustedNetworks),
	)
}

func setupAdminServer(providers *app.Providers) (*http.Server, error) {
	// prepare middlewares and interceptors
	loggingInterceptor := log.NewLoggingInterceptor()
	validatorInterceptor := validator.NewInterceptor(providers.Validator)
	privacyInterceptor := privacy.NewFilterInterceptor(privacy.SubjectResolverFunc(func(ctx context.Context, ar connect.AnyRequest) (string, []string, error) {
		claims := middleware.ClaimsFromContext(ctx)
		if claims == nil {
			return "", nil, nil
		}

		return claims.Subject, claims.AppMetadata.Authorization.Roles, nil
	}))

	errorInterceptor := middleware.NewErrorInterceptor()

	interceptors := connect.WithInterceptors(
		loggingInterceptor,
		validatorInterceptor,
		privacyInterceptor,
		errorInterceptor,
	)

	serveMux := http.NewServeMux()

	// User service
	userService, err := users.NewService(providers)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize users service: %w", err)
	}
	path, handler := idmv1connect.NewUserServiceHandler(
		userService,
		interceptors,
	)
	serveMux.Handle(path, handler)

	// Role service
	roleService := roles.NewService(providers)
	path, handler = idmv1connect.NewRoleServiceHandler(
		roleService,
		interceptors,
	)
	serveMux.Handle(path, handler)

	// Notify service
	notifyService := notify.New(providers)
	path, handler = idmv1connect.NewNotifyServiceHandler(
		notifyService,
		interceptors,
	)
	serveMux.Handle(path, handler)

	serveMux.Handle("/validate", auth.NewForwardAuthHandler(providers))

	return server.CreateWithOptions(
		providers.Config.AdminListenAddr,
		middleware.NewJWTMiddleware(
			providers.Config,
			providers.Datastore,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/validate" {
					log.L(r.Context()).Infof("adding fake admin claims to request: %s", r.URL.Path)

					clientIP := r.RemoteAddr

					// associate a dummy claims object to all request that are
					// received on the admin interface
					claims := jwt.Claims{
						Subject:  clientIP,
						ID:       clientIP,
						Audience: providers.Config.Audience,
						Issuer:   providers.Config.Domain,
						Scopes:   []jwt.Scope{jwt.ScopeAccess},
						Name:     clientIP,
						AppMetadata: &jwt.AppMetadata{
							TokenVersion: "2",
							Authorization: &jwt.Authorization{
								Roles: []string{"idm_superuser"},
							},
						},
					}

					ctx := middleware.ContextWithClaims(r.Context(), &claims)

					r = r.WithContext(ctx)
				}

				serveMux.ServeHTTP(w, r)
			}),
			func(r *http.Request) bool {
				// Skip JWT token verification for the /validate endpoint as
				// the ForwardAuthHanlder will take care of this on it's own due to special
				// handling of rejected or expired tokens.
				if r.URL.Path == "/validate" || strings.HasPrefix(r.URL.Path, "/webauthn") {
					return true
				}

				return false
			},
		),

		server.WithTrustedProxies(providers.Config.TrustedNetworks),
	)
}

func startServer(providers *app.Providers) error {
	publicServer, err := setupPublicServer(providers)
	if err != nil {
		return fmt.Errorf("failed to create public server: %w", err)
	}

	adminServer, err := setupAdminServer(providers)
	if err != nil {
		return fmt.Errorf("failed to create admin server: %w", err)
	}

	return server.Serve(context.Background(), publicServer, adminServer)
}
