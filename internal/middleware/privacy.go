package middleware

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/tierklinik-dobersberg/apis/pkg/privacy"
	"google.golang.org/protobuf/proto"
)

func NewPrivacyFilterInterceptor() connect.UnaryInterceptorFunc {
	return func(uf connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
			req, err := uf(ctx, ar)
			if err != nil {
				return nil, err
			}

			claims := ClaimsFromContext(ctx)
			if claims == nil {
				return req, nil
			}

			if err := privacy.FilterAllowedFields(req.Any().(proto.Message), claims.Subject, claims.AppMetadata.Authorization.Roles); err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}

			return req, nil
		}
	}
}
