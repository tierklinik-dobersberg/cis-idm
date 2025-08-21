package repo

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
)

//go:embed sql/migrations/*
var dbMigrations embed.FS

const sqlDialect = "sqlite3"

func Migrate(ctx context.Context, db *sql.DB) (int, error) {
	migrations := migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "sql/migrations",
	}

	n, err := migrate.Exec(db, sqlDialect, migrations, migrate.Up)
	if err != nil {
		return n, err
	}

	records, err := migrate.GetMigrationRecords(db, sqlDialect)
	if err != nil {
		log.L(ctx).Error("failed to get migration records", "error", err)

		return n, nil
	}

	if len(records) > 0 {
		log.L(ctx).Info("applied database migrations:")
		for _, r := range records {
			log.L(ctx).Info(" ✓ "+r.Id, "applied_at", r.AppliedAt)
		}
	} else {
		return 0, fmt.Errorf("failed to get any migration records")
	}

	return n, nil
}
