package acl

import (
	"context"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func NewInterceptor(reg *protoregistry.Files) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
			return next(ctx, ar)
		})
	}
}
