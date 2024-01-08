package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/tierklinik-dobersberg/apis/pkg/log"
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

type TransactionOptions func(opts *sql.TxOptions)

func ReadOnly() TransactionOptions {
	return func(tx *sql.TxOptions) {
		tx.ReadOnly = true
	}
}

func Isolation(lvl sql.IsolationLevel) TransactionOptions {
	return func(opts *sql.TxOptions) {
		opts.Isolation = lvl
	}
}

func RunInTransaction[R any](ctx context.Context, q *Queries, fn func(tx *Queries) (R, error), opts ...TransactionOptions) (R, error) {
	txOpts := new(sql.TxOptions)

	for _, fn := range opts {
		fn(txOpts)
	}

	tx, err := q.Tx(ctx, txOpts)
	if err != nil {
		return *new(R), fmt.Errorf("failed to prepare transaction: %w", err)
	}

	ds := q.WithTx(tx)
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.L(ctx).Errorf("failed to rollback transaction: %s", err)
		}
	}()

	res, err := fn(ds)
	if err != nil {
		return res, err
	}

	if err := tx.Commit(); err != nil {
		return res, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return res, err
}
