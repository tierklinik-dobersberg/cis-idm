package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/bufbuild/protovalidate-go"
	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/apis/pkg/discovery"
	"github.com/tierklinik-dobersberg/apis/pkg/discovery/consuldiscover"
	"github.com/tierklinik-dobersberg/apis/pkg/discovery/wellknown"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/bootstrap"
	"github.com/tierklinik-dobersberg/cis-idm/internal/cache"
	"github.com/tierklinik-dobersberg/cis-idm/internal/common"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/mailer"
	"github.com/tierklinik-dobersberg/cis-idm/internal/policy"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/sms"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configFilePath := "/etc/cisidm/config.hcl"

	if path := os.Getenv("CONFIG_FILE"); path != "" {
		configFilePath = path
	}

	// get configuration from environment
	cfg, err := config.LoadFile(configFilePath)
	if err != nil {
		logrus.Fatalf("failed to parse config file %s: %s", configFilePath, err)
	}
	logrus.Infof("sucessfully loaded configuration")

	if cfg.LogLevel != "" {
		logrus.Infof("switching log level to %q", cfg.LogLevel)

		lvl, err := logrus.ParseLevel(cfg.LogLevel)
		if err != nil {
			logrus.Fatalf("failed to parse log level %q: %w", cfg.LogLevel, err)
		}

		logrus.SetLevel(lvl)
	}

	// prepare all application providers.
	providers, err := setupAppProviders(ctx, *cfg)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	// bootstrap the application by creating required roles.
	if err := bootstrap.Bootstrap(ctx, providers.Config, providers.Datastore); err != nil {
		cancel()

		logrus.Fatalf("failed to bootstrap: %s", err)
	}

	// Register at service catalog
	catalog, err := consuldiscover.NewFromEnv()
	if err != nil {
		cancel()

		logrus.Fatalf("failed to get service catalog client: %s", err)
	}

	if err := discovery.Register(ctx, catalog, &discovery.ServiceInstance{
		Name:    string(wellknown.IdmV1ServiceScope),
		Address: cfg.Server.AdminListenAddr,
	}); err != nil {
		logrus.Errorf("failed to register service at catalog: %s", err)
	}

	// finally, start of the HTTP/2 servers...
	if err := startServer(providers); err != nil {
		logrus.Fatalf("failed to start server: %s", err)
	}
}

func setupAppProviders(ctx context.Context, cfg config.Config) (*app.Providers, error) {
	db, err := sql.Open("sqlite3_extended", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database at %q: %q", cfg.DatabaseURL, err)
	}

	// connect to rqlite
	datastore := repo.New(db)

	// Try to create/migrate the users tables.
	if n, err := repo.Migrate(ctx, db); err == nil {
		log.L(ctx).Info("successfully applied migrations", "count", n)
	} else {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	logrus.Infof("successfully migrated database")

	// Prepare the template engine used for notification templates.
	tmplEngine, err := tmpl.New(datastore)
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

	// prepare the mailer
	var mailSender mailer.Mailer
	if cfg.MailConfig != nil {
		mailSender, err = mailer.New(mailer.Account{
			Host:          cfg.MailConfig.Host,
			Port:          cfg.MailConfig.Port,
			Username:      cfg.MailConfig.Username,
			Password:      cfg.MailConfig.Password,
			From:          cfg.MailConfig.From,
			AllowInsecure: cfg.MailConfig.AllowInsecure,
			UseSSL:        cfg.MailConfig.UseSSL,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to setup mail provider: %w", err)
		}
	} else {
		mailSender = new(mailer.NoOpMailer)
	}

	cache := cache.NewInMemoryCache()
	commonService := common.New(datastore, cfg, cache)

	// prepare engine options
	options := []policy.EngineOption{}
	for _, p := range cfg.PolicyConfig.Policies {
		options = append(options, policy.WithRawPolicy(p.Name, p.Content))
	}

	if cfg.PolicyConfig.Debug {
		options = append(options, policy.WithDebug())
	}

	// prepare rego policy engine
	engine, err := policy.NewEngine(ctx, cfg.PolicyConfig.Directories, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare policy engine: %w", err)
	}

	providers := &app.Providers{
		TemplateEngine: tmplEngine,
		SMSSender:      smsProvider,
		Mailer:         mailSender,
		Datastore:      datastore,
		Config:         cfg,
		Common:         commonService,
		ProtoRegistry:  reg,
		Validator:      validator,
		Cache:          cache,
		PolicyEngine:   engine,
	}

	return providers, nil
}
