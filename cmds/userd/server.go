package main

import (
	"net/http"

	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/userd/v1/userdv1connect"
	"github.com/tierklinik-dobersberg/cis-userd/internal/auth"
	"github.com/tierklinik-dobersberg/cis-userd/internal/config"
	"github.com/tierklinik-dobersberg/cis-userd/internal/repo"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func startServer(repo *repo.UserRepo, cfg config.Config) error {
	authService := auth.NewService(repo, cfg)

	mux := http.NewServeMux()
	path, handler := userdv1connect.NewAuthServiceHandler(authService)

	mux.Handle(path, handler)

	return http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
