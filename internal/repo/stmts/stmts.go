package stmts

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/mitchellh/mapstructure"
	"github.com/rqlite/gorqlite"
)

type (
	Statement[R any] struct {
		Query  string
		Args   []string
		Result R
	}

	StatementList []Statement[any]
)

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

func (list StatementList) Prepare(args []any) ([]gorqlite.ParameterizedStatement, error) {
	result := make([]gorqlite.ParameterizedStatement, len(list))

	if args != nil {
		if len(args) != len(list) {
			return nil, fmt.Errorf("arguments length does not match statement length")
		}
	}

	for idx, s := range list {
		v, err := s.Prepare(args[idx])
		if err != nil {
			return nil, fmt.Errorf("failed to prepare statement %d: %w", idx, err)
		}

		result[idx] = v
	}

	return result, nil
}

func (list StatementList) Write(ctx context.Context, conn *gorqlite.Connection, args []any) error {
	stmts, err := list.Prepare(args)
	if err != nil {
		return err
	}

	writeResult, err := conn.WriteParameterizedContext(ctx, stmts)
	if err != nil {
		return err
	}

	merr := new(multierror.Error)
	for idx, wr := range writeResult {
		if wr.Err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("stmt#%d: %w", idx, err))
		}
	}

	return merr.ErrorOrNil()
}