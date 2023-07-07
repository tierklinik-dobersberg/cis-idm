package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tierklinik-dobersberg/cis-idm/internal/bootstrap"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
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

	// connect to rqlite
	userRepo, err := repo.New(cfg.DatabaseURL)
	if err != nil {
		cancel()
		logrus.Fatalf("failed to connect to database: %s", err)
	}

	// Get the current leader and print a info message about the cluster
	leader, err := userRepo.Conn.Leader()
	if err != nil {
		cancel()
		logrus.Fatalf("rqlite does not yet have a leader elected: %s", err)
	}

	logrus.WithFields(logrus.Fields{
		"cluster": cfg.DatabaseURL,
		"leader":  leader,
	}).Infof("connected to rqlite cluster")

	// Try to create/migrate the users tables.
	if err := userRepo.Migrate(ctx); err != nil {
		cancel()
		logrus.Fatalf("failed to migrate database: %s", err)
	}

	logrus.Infof("successfully migrated database")

	if err := bootstrap.Bootstrap(ctx, cfg, userRepo); err != nil {
		cancel()

		logrus.Fatalf("failed to bootstrap: %s", err)
	}

	if err := startServer(userRepo, cfg); err != nil {
		logrus.Fatalf("failed to start server: %s", err)
	}
}
