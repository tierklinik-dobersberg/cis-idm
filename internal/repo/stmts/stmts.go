package stmts

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/rqlite/gorqlite"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
)

type Statement[R any] struct {
	Query  string
	Args   []string
	Result R
}

var (
	ErrInvalidArgCount = errors.New("invalid number of arguments")
	ErrMissingArgument = errors.New("missing argument value")
	ErrNoResults       = errors.New("no results")
	ErrNoRowsAffected  = errors.New("no rows affected")
)

func (stmt Statement[R]) Prepare(args any) (gorqlite.ParameterizedStatement, error) {
	argMap := make(map[string]any)

	if args != nil {
		if err := mapstructure.Decode(args, &argMap); err != nil {
			return gorqlite.ParameterizedStatement{}, err
		}
	}

	argList := make([]any, 0, len(argMap))
	for _, argName := range stmt.Args {
		argValue, ok := argMap[argName]
		if !ok {
			return gorqlite.ParameterizedStatement{}, fmt.Errorf("%s: %w", argName, ErrMissingArgument)
		}
		argList = append(argList, argValue)

	}

	if len(argList) < len(stmt.Args) {
		return gorqlite.ParameterizedStatement{}, ErrInvalidArgCount
	}

	return gorqlite.ParameterizedStatement{
		Query:     stmt.Query,
		Arguments: argList,
	}, nil
}

func (stmt Statement[R]) Write(ctx context.Context, conn *gorqlite.Connection, args any) error {
	pStmt, err := stmt.Prepare(args)
	if err != nil {
		return err
	}

	if os.Getenv("DEBUG") != "" {
		middleware.L(ctx).
			Infof(pStmt.Query)
	}

	writeResult, err := conn.WriteOneParameterizedContext(ctx, pStmt)
	if err != nil {
		if writeResult.Err != nil {
			return writeResult.Err
		}
		return err
	}

	if writeResult.Err != nil {
		return writeResult.Err
	}

	if writeResult.RowsAffected == 0 {
		return ErrNoRowsAffected
	}

	return nil
}
