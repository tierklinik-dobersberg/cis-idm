package policy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tierklinik-dobersberg/cis-idm/internal/policy"
)

func Test_Engine_QueryOne(t *testing.T) {
	engine, err := policy.NewEngine(context.TODO(), []string{"./testdata"})
	require.NoError(t, err)

	type TestResult struct {
		Allow   bool           `mapstructure:"allow"`
		Headers map[string]any `mapstructure:"headers"`
	}

	var result TestResult

	input := map[string]any{
		"subject": map[string]any{
			"username": "test",
			"roles":    []string{"idm_superuser"},
		},
		"path": "/test",
	}

	err = engine.QueryOne(context.TODO(), "data.cisidm.forward_auth", input, &result)
	require.NoError(t, err)

	expected := TestResult{
		Allow: true,
		Headers: map[string]any{
			"Foo": "Bar",
		},
	}

	assert.Equal(t, expected, result)
}
