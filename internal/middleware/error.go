package middleware

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bufbuild/connect-go"
)

func NewErrorInterceptor() connect.UnaryInterceptorFunc {
	return func(uf connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := uf(ctx, ar)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, connect.NewError(connect.CodeNotFound, err)
				}
			}

			return resp, err
		}
	}
}
