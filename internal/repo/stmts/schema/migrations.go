package schema

import (
	"context"
	"fmt"

	"github.com/rqlite/gorqlite"
)

func Migrate(ctx context.Context, conn *gorqlite.Connection) error {
	// for now, we just create the table schema.
	if err := createSchema.Write(ctx, conn, nil); err != nil {
		return fmt.Errorf("failed to create table schema: %w", err)
	}

	return nil
}