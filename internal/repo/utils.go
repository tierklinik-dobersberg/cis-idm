package repo

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrTxNotSupported = errors.New("tx not supported")
)

type Transactioner interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

func (q *Queries) Tx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	if txr, ok := q.db.(Transactioner); ok {
		return txr.BeginTx(ctx, opts)
	}

	return nil, ErrTxNotSupported
}