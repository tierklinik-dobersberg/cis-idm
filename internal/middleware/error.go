package middleware

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func NewErrorInterceptor() connect.UnaryInterceptorFunc {
	return func(uf connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := uf(ctx, ar)
			if err != nil {
				if errors.Is(err, stmts.ErrNoResults) {
					return nil, connect.NewError(connect.CodeNotFound, err)
				}
			}

			return resp, err
		}
	}
}
