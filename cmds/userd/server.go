package main

import (
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/cis-idm/internal/auth"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware/acl"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func startServer(repo *repo.Repo, cfg config.Config) error {
	reg := new(protoregistry.Files)

	for _, file := range []protoreflect.FileDescriptor{
		idmv1.File_tkd_idm_v1_auth_service_proto,
		idmv1.File_tkd_idm_v1_descriptor_proto,
		idmv1.File_tkd_idm_v1_self_service_proto,
		idmv1.File_tkd_idm_v1_user_proto,
		idmv1.File_tkd_idm_v1_user_service_proto,
	} {
		if err := reg.RegisterFile(file); err != nil {
			return fmt.Errorf("failed to register %s at protoregistry: %w", file.Name(), err)
		}
	}

	loggingInterceptor := middleware.NewLoggingInterceptor()
	authInterceptor := middleware.NewAuthInterceptor(cfg, reg)
	aclInterceptor := acl.NewInterceptor(reg)

	authService, err := auth.NewService(repo, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize auth service: %w", err)
	}

	mux := http.NewServeMux()
	path, handler := idmv1connect.NewAuthServiceHandler(
		authService,
		connect.WithInterceptors(loggingInterceptor, authInterceptor, aclInterceptor),
	)

	mux.Handle(path, handler)

	return http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(
			middleware.NewJWTMiddleware(cfg, repo, mux),
			&http2.Server{},
		),
	)
}
