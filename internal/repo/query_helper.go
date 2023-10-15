package repo

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/rqlite/gorqlite"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func Query[T any](ctx context.Context, stmt stmts.Statement[T], conn *gorqlite.Connection, args any) ([]T, error) {
	pStmt, err := stmt.Prepare(args)
	if err != nil {
		return nil, err
	}

	queryResult, err := conn.QueryOneParameterizedContext(ctx, pStmt)
	if err != nil {
		if queryResult.Err != nil {
			return nil, queryResult.Err
		}

		return nil, err
	}

	if os.Getenv("DEBUG_SQL") != "" {
		log.L(ctx).
			Infof(pStmt.Query)
	}

	typeOf := reflect.TypeOf(stmt.Result)

	// if T is nil than that statement is not expected to return data
	if typeOf == nil {
		return nil, nil
	}

	results := make([]T, 0, queryResult.NumRows())
	for queryResult.Next() {
		m, err := queryResult.Map()
		if err != nil {
			return results, err
		}

		if os.Getenv("DEBUG_SQL") != "" {
			log.L(ctx).
				WithField("response", m).
				Infof("DEBUG: response")
		}

		newObj := reflect.New(typeOf)
		obj := newObj.Interface().(*T)
		if err := mapstructure.WeakDecode(m, obj); err != nil {
			return results, err
		}

		results = append(results, *obj)
	}

	return results, nil
}

func QueryOne[T any](ctx context.Context, stmt stmts.Statement[T], conn *gorqlite.Connection, args any) (T, error) {
	results, err := Query(ctx, stmt, conn, args)
	if err != nil {
		return stmt.Result, err
	}

	if len(results) == 0 {
		return stmt.Result, stmts.ErrNoResults
	}

	if len(results) > 1 {
		return stmt.Result, fmt.Errorf("query returned more than one result")
	}

	return results[0], nil
}
