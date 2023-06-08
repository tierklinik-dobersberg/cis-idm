package middleware

import (
	"context"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
)

var loggerContextKey = struct{ s string }{s: "logger"}

func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func L(ctx context.Context) *logrus.Entry {
	l, _ := ctx.Value(loggerContextKey).(*logrus.Entry)
	if l == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return l
}

func NewLoggingInterceptor() connect.UnaryInterceptorFunc {
	return func(uf connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			resp, err := uf(ctx, ar)

			l := L(ctx).WithFields(logrus.Fields{
				"duration": time.Since(start),
			})

			if err != nil {
				l = l.WithError(err)
			}

			l.Infof("connect request handled")

			return resp, err
		}
	}
}
