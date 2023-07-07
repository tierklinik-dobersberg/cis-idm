package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bufbuild/protovalidate-go"
	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/bootstrap"
	"github.com/tierklinik-dobersberg/cis-idm/internal/cache"
	"github.com/tierklinik-dobersberg/cis-idm/internal/common"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/sms"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configFilePath := "/etc/cisidm/config.yml"

	if path := os.Getenv("CONFIG_FILE"); path != "" {
		configFilePath = path
	}

	// get configuration from environment
	cfg, err := config.FromEnvironment(ctx, configFilePath)
	if err != nil {
		logrus.Fatalf("failed to parse config from environment: %s", err)
	}
	logrus.Infof("sucessfully loaded configuration")

	// prepare all application providers.
	providers, err := setupAppProviders(ctx, cfg)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	// bootstrap the application by creating required roles.
	if err := bootstrap.Bootstrap(ctx, cfg, providers.Datastore); err != nil {
		cancel()

		logrus.Fatalf("failed to bootstrap: %s", err)
	}

	// finally, start of the HTTP/2 servers...
	if err := startServer(providers); err != nil {
		logrus.Fatalf("failed to start server: %s", err)
	}
}

func setupAppProviders(ctx context.Context, cfg config.Config) (*app.Providers, error) {
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

	// Prepare the template engine used for notification templates.
	tmplEngine, err := tmpl.New()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare template engine: %w", err)
	}

	// Prepare the proto registry used for reflection operations.
	reg, err := getProtoRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to create proto registry: %w", err)
	}

	// Prepare a new validator used to validate incoming requests.
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create protovalidate.Validator: %w", err)
	}

	// prepare a new Twilio SMS provider.
	var smsProvider sms.Sender

	if cfg.Twilio != nil {
		smsProvider, err = sms.New(sms.Account{
			From:        cfg.Twilio.From,
			AccountSid:  cfg.Twilio.AccountSid,
			AccessToken: cfg.Twilio.AccessToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to setup sms provider: %w", err)
		}
	} else {
		smsProvider = sms.NewNoopProvider()
	}

	commonService := common.New(datastore, cfg)

	providers := &app.Providers{
		TemplateEngine: tmplEngine,
		SMSSender:      smsProvider,
		Datastore:      datastore,
		Config:         cfg,
		Common:         commonService,
		ProtoRegistry:  reg,
		Validator:      validator,
		Cache:          cache.NewInMemoryCache(), // TODO(ppacher): support redis here for HA
	}

	return providers, nil
}
