package common

import (
	"context"

	servertiming "github.com/mitchellh/go-server-timing"
)

func Timing(ctx context.Context, name string, fn func() error) error {
	t := servertiming.FromContext(ctx)

	if t != nil {
		m := t.NewMetric(name).Start()
		defer m.Stop()
	}

	return fn()
}
