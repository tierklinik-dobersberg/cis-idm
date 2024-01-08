package permission_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tierklinik-dobersberg/cis-idm/internal/permission"
)

func Test_ParseTree(t *testing.T) {
	tree := getTree(t, `idm:
  users:
    - read
    - write
    - invite
  roles:
    read: {}
    write:
      - create
      - update
`)

	tree.Insert("idm:roles:write:delete")

	assert.Equal(t, permission.Tree{
		"idm": {
			"users": {
				"read":   {},
				"write":  {},
				"invite": {},
			},
			"roles": {
				"read": {},
				"write": {
					"create": {},
					"update": {},
					"delete": {},
				},
			},
		},
	}, tree)
}

func Test_MergeTree(t *testing.T) {
	users := permission.Tree{
		"users": {
			"read":  {},
			"write": {},
		},
	}

	inviter := permission.Tree{
		"users": {
			"write": {
				"invite": {},
			},
		},
	}

	roles := permission.Tree{
		"roles": {
			"read": {},
			"write": {
				"create": {},
				"update": {},
				"delete": {},
			},
		},
	}

	result := permission.Tree{}.
		MergeFrom(users).
		MergeFrom(inviter).
		MergeFrom(roles)

	assert.Equal(t, permission.Tree{
		"users": {
			"read": {},
			"write": {
				"invite": {},
			},
		},
		"roles": {
			"read": {},
			"write": {
				"create": {},
				"update": {},
				"delete": {},
			},
		},
	}, result)
}

func Test_Resolve(t *testing.T) {
	tree := permission.Tree{
		"idm": {
			"users": {
				"read":   {},
				"write":  {},
				"invite": {},
			},
			"roles": {
				"read": {},
				"write": {
					"create": {},
					"update": {},
					"delete": {},
				},
			},
		},
	}

	expectedResult := []string{
		"does-not-exist:foo:bar",
		"idm:users:read",
		"idm:users:write",
		"idm:users:invite",
		"idm:roles:write:create",
		"idm:roles:write:update",
		"idm:roles:write:delete",
		"idm:roles:read",
	}
	slices.Sort(expectedResult)

	set, err := tree.Resolve([]string{"idm", "does-not-exist:foo:bar"})
	require.NoError(t, err)

	assert.Equal(t, expectedResult, set)
}

func getTree(t *testing.T, perm string) permission.Tree {
	t.Helper()

	blob, err := yaml.YAMLToJSON([]byte(perm))
	if err != nil {
		fmt.Println(perm)
		require.NoError(t, err)
	}

	var tree permission.Tree
	if err := json.Unmarshal(blob, &tree); err != nil {
		require.NoError(t, err)
	}

	return tree
}
