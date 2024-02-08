package repo

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/mattn/go-sqlite3"
)

func init() {
	sql.Register("sqlite3_extended", &sqlite3.SQLiteDriver{
		ConnectHook: func(sc *sqlite3.SQLiteConn) error {
			if err := sc.RegisterFunc("uuid", func() (string, error) {
				u, err := uuid.NewV4()
				if err != nil {
					return "", err
				}

				return u.String(), nil
			}, false); err != nil {
				return fmt.Errorf("failed to register uuid function: %w", err)
			}

			_, err := sc.Exec("PRAGMA foreign_keys = ON;", nil)
			if err != nil {
				return fmt.Errorf("failed to enable foreign key support: %w", err)
			}

			return nil
		},
	})
}
