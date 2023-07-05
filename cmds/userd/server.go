package main

import (
	"context"
	"embed"
	"encoding/json"
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

	publicListenerMux := http.NewServeMux()

	commonService := common.New(repo, cfg)

	// Setup Auth
	authService, err := auth.NewService(repo, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth service: %w", err)
	}
	path, handler := idmv1connect.NewAuthServiceHandler(
		authService,
		interceptors,
	)
	publicListenerMux.Handle(path, handler)

	// Setup Self-Service
	selfserviceService, err := selfservice.NewService(cfg, repo, commonService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize self-service service: %w", err)
	}
	path, handler = idmv1connect.NewSelfServiceServiceHandler(
		selfserviceService,
		interceptors,
	)
	publicListenerMux.Handle(path, handler)

	// User service
	userService, err := users.NewService(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize users service: %w", err)
	}
	path, handler = idmv1connect.NewUserServiceHandler(
		userService,
		interceptors,
	)
	publicListenerMux.Handle(path, handler)

	publicListenerMux.Handle("/config.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")

		if err := enc.Encode(map[string]any{
			"domain":                    cfg.Domain,
			"loginURL":                  cfg.LoginRedirectURL,
			"siteName":                  cfg.SiteName,
			"siteNameUrl":               cfg.SiteNameURL,
			"registrationRequiresToken": cfg.RegistrationRequiresToken,
			"features":                  cfg.FeatureMap,
		}); err != nil {
			middleware.L(r.Context()).Errorf("failed to encode service config: %s", err)
		}
	}))

	c := cors.New(cors.Options{
		AllowedOrigins: append([]string{
			fmt.Sprintf("http://%s", cfg.PublicListenAddr),
			fmt.Sprintf("http://%s", cfg.Domain),
		}, cfg.AllowedOrigins...),
		AllowCredentials: true,
		AllowedHeaders:   []string{"Connect-Protocol-Version", "Content-Type", "Authentication"},
		Debug:            os.Getenv("DEBUG") != "",
	})

	staticFilesHandler, err := getStaticFilesHandler(cfg.StaticFiles)
	if err != nil {
		return nil, err
	}

	publicListenerMux.Handle("/", staticFilesHandler)
	publicListenerMux.Handle("/validate", auth.NewForwardAuthHandler(cfg, repo))
	publicListenerMux.Handle("/avatar/", users.NewAvatarHandler(repo))

	webauthnHandler, err := webauthn.New(cfg, authService, repo)
	if err != nil {
		return nil, err
	}

	publicListenerMux.Handle("/webauthn/", http.StripPrefix("/webauthn", webauthnHandler))

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

func startServer(repo *repo.Repo, cfg config.Config) error {
	reg, err := getProtoRegistry()
	if err != nil {
		return err
	}

	validator, err := protovalidate.New()
	if err != nil {
		return err
	}

	publicServer, err := setupPublicServer(repo, cfg, reg, validator)
	if err != nil {
		return err
	}

	adminServer := &http.Server{
		Addr: cfg.AdminListenAddr,
		Handler: h2c.NewHandler(
			nil, // FIXME
			&http2.Server{},
		),
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
