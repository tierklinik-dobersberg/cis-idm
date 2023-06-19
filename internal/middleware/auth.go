package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
	commonv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/common/v1"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var (
	claimsContextKey = struct{ s string }{s: "claims-context-key"}
)

// ContextWithClaims returns a new context.Context with claims attached.
// Use ClaimsFromContext to retrieve the claims.
func ContextWithClaims(ctx context.Context, claims *jwt.Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

// ClaimsFromContext returns the JWT claims associated with ctx.
func ClaimsFromContext(ctx context.Context) *jwt.Claims {
	claims, _ := ctx.Value(claimsContextKey).(*jwt.Claims)
	return claims
}

func NewAuthInterceptor(cfg config.Config, reg *protoregistry.Files) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			parts := strings.Split(req.Spec().Procedure, "/")

			methodDesc := getMethodDesc(reg, parts[1], parts[2])
			if methodDesc == nil {
				return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to find method descriptor for %s", req.Spec().Procedure))
			}

			l := L(ctx).WithField("method", methodDesc.FullName())

			claims := ClaimsFromContext(ctx)

			opts, ok := proto.GetExtension(methodDesc.Options(), commonv1.E_Auth).(*commonv1.AuthDecorator)

			if ok && opts != nil {
				L(ctx).Infof("checking authentication requirement: %#v", opts)
				switch opts.Require {
				case commonv1.AuthRequirement_AUTH_REQ_REQUIRED:
					l.Infof("service method requires authentication")
					if claims == nil {
						return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("not access token provided"))
					}

					// make sure the user has at least one of the required roles assigned
					if len(opts.AllowedRoles) > 0 {
						isAllowed := false

						for _, allowedRole := range opts.AllowedRoles {
							if slices.Contains(claims.AppMetadata.Authorization.Roles, allowedRole) {
								isAllowed = true
								break
							}
						}

						if !isAllowed {
							return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("access token does not include one of the required roles"))
						}
					}

				case commonv1.AuthRequirement_AUTH_REQ_UNSPECIFIED:
					// nothing to do
				default:
					l.WithField("requirement", opts.String()).Infof("unhandeled authentication requirement")
				}
			} else {
				l.Infof("not authentication requirement specified for service method")
			}

			ctx = WithLogger(ctx, l)

			return next(ctx, req)
		})
	}

	return interceptor
}

func getMethodDesc(reg *protoregistry.Files, fqServiceName string, methodName string) protoreflect.MethodDescriptor {
	serviceNameParts := strings.Split(fqServiceName, ".")
	serviceName := serviceNameParts[len(serviceNameParts)-1]

	var methodDesc protoreflect.MethodDescriptor

	reg.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		if strings.HasPrefix(fqServiceName, string(fd.FullName())) {
			serviceDesc := fd.Services().ByName(protoreflect.Name(serviceName))
			if serviceDesc != nil {

				methodDesc = serviceDesc.Methods().ByName(protoreflect.Name(methodName))
				if methodDesc != nil {
					return false
				}
			}
		}

		return true
	})

	return methodDesc
}

func FindCookie(cookieName string, headers http.Header) *http.Cookie {
	// we create a dummy http request so we can use the cookie parser
	// from the stdlib which is, unfortunately, not exported for direct
	// use.
	dummyReq := http.Request{Header: headers}

	cookies := dummyReq.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			return cookie
		}
	}

	return nil
}
