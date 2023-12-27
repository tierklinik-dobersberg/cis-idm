package repo

import (
	"context"

	"github.com/rqlite/gorqlite"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts/schema"
)

type Repo struct {
	Conn *gorqlite.Connection
}

func New(endpoint string) (*Repo, error) {
	conn, err := gorqlite.Open(endpoint)
	if err != nil {
		return nil, err
	}

	return &Repo{Conn: conn}, nil
}

func (repo *Repo) Migrate(ctx context.Context) error {
	return schema.Migrate(ctx, repo.Conn)
}

