package repo

import (
	"context"
	"database/sql"
	"embed"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql/migrations/*
var dbMigrations embed.FS

func Migrate(ctx context.Context, db *sql.DB) error {
	migrations := migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       ".",
	}

	_, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}
