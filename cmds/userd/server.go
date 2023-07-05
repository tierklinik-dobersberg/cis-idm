package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/bufbuild/protovalidate-go"
	"github.com/rs/cors"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/cis-idm/internal/auth"
	"github.com/tierklinik-dobersberg/cis-idm/internal/common"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware/acl"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/selfservice"
	"github.com/tierklinik-dobersberg/cis-idm/internal/users"
	"github.com/tierklinik-dobersberg/cis-idm/internal/webauthn"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
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
		return common.ServeSPA(http.FS(webapp), "index.html"), nil
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

	return common.ServeSPA(http.Dir(path), "index.html"), nil
}

func setupPublicServer(repo *repo.Repo, cfg config.Config, reg *protoregistry.Files, validator *protovalidate.Validator) (*http.Server, error) {
	// prepare middlewares and interceptors
	loggingInterceptor := middleware.NewLoggingInterceptor()
	authInterceptor := middleware.NewAuthInterceptor(cfg, reg)
	aclInterceptor := acl.NewInterceptor(reg)
	validatorInterceptor := middleware.NewValidationInterceptor(validator)
	privacyInterceptor := middleware.NewPrivacyFilterInterceptor()

	interceptors := connect.WithInterceptors(
		loggingInterceptor,
		authInterceptor,
		aclInterceptor,
		validatorInterceptor,
		privacyInterceptor,
	)

	// create a new servemux to handle our routes.
	publicListenerMux := http.NewServeMux()

	// create a common.Service instance that shares code for the AuthService
	// and UserService instances.
	commonService := common.New(repo, cfg)

	// Setup and serve the AuthService
	authService, err := auth.NewService(repo, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth service: %w", err)
	}
	path, handler := idmv1connect.NewAuthServiceHandler(
		authService,
		interceptors,
	)
	publicListenerMux.Handle(path, handler)

	// Setup and serve the Self-Service
	selfserviceService, err := selfservice.NewService(cfg, repo, commonService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize self-service service: %w", err)
	}
	path, handler = idmv1connect.NewSelfServiceServiceHandler(
		selfserviceService,
		interceptors,
	)
	publicListenerMux.Handle(path, handler)

	// Setup and serve the User service.
	userService, err := users.NewService(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize users service: %w", err)
	}
	path, handler = idmv1connect.NewUserServiceHandler(
		userService,
		interceptors,
	)
	publicListenerMux.Handle(path, handler)

	// Serve basic configuration for the UI on /config.json
	publicListenerMux.Handle("/config.json", config.NewConfigHandler(cfg))

	// Setup CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: append([]string{
			cfg.PublicURL,
			fmt.Sprintf("http://%s", cfg.Domain),
			fmt.Sprintf("https://%s", cfg.Domain),
		}, cfg.AllowedOrigins...),
		AllowCredentials: true,
		AllowedHeaders:   []string{"Connect-Protocol-Version", "Content-Type", "Authentication"},
		Debug:            os.Getenv("DEBUG") != "",
	})

	// Get a static file handler.
	// This will either return a handler for the embed.FS, a local directory using http.Dir
	// or a reverse proxy to some other service.
	staticFilesHandler, err := getStaticFilesHandler(cfg.StaticFiles)
	if err != nil {
		return nil, err
	}

	publicListenerMux.Handle("/", staticFilesHandler)

	// Setup the avatar handler. This does not use connect-go since we need
	// to response with either HTTP redirects (if the user avatar is a URL)
	// or with plain bytes and an approriate content-type if the user avatar
	// is a dataurl.
	publicListenerMux.Handle("/avatar/", users.NewAvatarHandler(repo))

	// setup the webauthn handlers for registration and login.
	// TODO(ppacher): migrate those to connect-go/protobuf style endpoints
	// as the browser does not actually care about how this is implemented.
	webauthnHandler, err := webauthn.New(cfg, authService, repo)
	if err != nil {
		return nil, err
	}
	publicListenerMux.Handle("/webauthn/", http.StripPrefix("/webauthn", webauthnHandler))

	// finally, return a http.Server that uses h2c for HTTP/2 support and
	// wrap the finnal handler in CORS and a JWT middleware.
	return &http.Server{
		Addr: cfg.PublicListenAddr,
		Handler: h2c.NewHandler(
			c.Handler(middleware.NewJWTMiddleware(cfg, repo, publicListenerMux, func(r *http.Request) bool {
				// Skip JWT token verification for the /validate endpoint as
				// the ForwardAuthHanlder will take care of this on it's own due to special
				// handling of rejected or expired tokens.
				if r.URL.Path == "validate" || strings.HasPrefix(r.URL.Path, "webauthn/") {
					return true
				}

				return false
			})),
			&http2.Server{},
		),
	}, nil
}

func setupAdminServer(repo *repo.Repo, cfg config.Config, reg *protoregistry.Files, validator *protovalidate.Validator) (*http.Server, error) {
	// prepare middlewares and interceptors
	loggingInterceptor := middleware.NewLoggingInterceptor()
	validatorInterceptor := middleware.NewValidationInterceptor(validator)
	privacyInterceptor := middleware.NewPrivacyFilterInterceptor()

	interceptors := connect.WithInterceptors(
		loggingInterceptor,
		validatorInterceptor,
		privacyInterceptor,
	)

	serveMux := http.NewServeMux()

	// User service
	userService, err := users.NewService(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize users service: %w", err)
	}
	path, handler := idmv1connect.NewUserServiceHandler(
		userService,
		interceptors,
	)
	serveMux.Handle(path, handler)

	serveMux.Handle("/validate", auth.NewForwardAuthHandler(cfg, repo))

	return &http.Server{
		Addr: cfg.AdminListenAddr,
		Handler: h2c.NewHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				clientIP := r.RemoteAddr

				// associate a dummy claims object to all request that are
				// received on the admin interface
				claims := jwt.Claims{
					Subject:  clientIP,
					ID:       clientIP,
					Audience: cfg.Audience,
					Issuer:   cfg.Domain,
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

				serveMux.ServeHTTP(w, r)
			}),
			&http2.Server{},
		),
	}, nil
}

func startServer(repo *repo.Repo, cfg config.Config) error {
	reg, err := getProtoRegistry()
	if err != nil {
		return fmt.Errorf("failed to create proto registry: %w", err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("failed to create protovalidate.Validator: %w", err)
	}

	publicServer, err := setupPublicServer(repo, cfg, reg, validator)
	if err != nil {
		return fmt.Errorf("failed to create public server: %w", err)
	}

	adminServer, err := setupAdminServer(repo, cfg, reg, validator)
	if err != nil {
		return fmt.Errorf("failed to create admin server: %w", err)
	}

	errgrp, ctx := errgroup.WithContext(context.Background())

	go func() {
		<-ctx.Done()

		timeOutCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		publicServer.Shutdown(timeOutCtx)
		adminServer.Shutdown(timeOutCtx)
	}()

	errgrp.Go(publicServer.ListenAndServe)
	errgrp.Go(adminServer.ListenAndServe)

	if err := errgrp.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
