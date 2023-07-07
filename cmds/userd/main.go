package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bufbuild/protovalidate-go"
	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/bootstrap"
	"github.com/tierklinik-dobersberg/cis-idm/internal/common"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// get configuration from environment
	cfg, err := config.FromEnvironment(ctx, os.Args[1])
	if err != nil {
		logrus.Fatalf("failed to parse config from environment: %s", err)
	}
	logrus.Infof("sucessfully loaded configuration")

	providers, err := setupApp(ctx, cfg)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	if err := bootstrap.Bootstrap(ctx, cfg, providers.Datastore); err != nil {
		cancel()

		logrus.Fatalf("failed to bootstrap: %s", err)
	}

	if err := startServer(providers); err != nil {
		logrus.Fatalf("failed to start server: %s", err)
	}
}

func setupApp(ctx context.Context, cfg config.Config) (*app.Providers, error) {
	// connect to rqlite
	datastore, err := repo.New(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the current leader and print a info message about the cluster
	leader, err := datastore.Conn.Leader()
	if err != nil {
		return nil, fmt.Errorf("rqlite does not yet have a leader elected: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"cluster": cfg.DatabaseURL,
		"leader":  leader,
	}).Infof("connected to rqlite cluster")

	// Try to create/migrate the users tables.
	if err := datastore.Migrate(ctx); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	logrus.Infof("successfully migrated database")

	tmplEngine, err := tmpl.New()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare template engine: %w", err)
	}

	reg, err := getProtoRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to create proto registry: %w", err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create protovalidate.Validator: %w", err)
	}

	commonService := common.New(datastore, cfg)

	providers := &app.Providers{
		TemplateEngine: tmplEngine,
		SMSSender:      nil,
		Datastore:      datastore,
		Config:         cfg,
		Common:         commonService,
		ProtoRegistry:  reg,
		Validator:      validator,
	}

	return providers, nil
}
