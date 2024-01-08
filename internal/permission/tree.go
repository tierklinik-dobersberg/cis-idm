package permission

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

/*
```yaml
permissions:

	users:
		write:
			- owner
		read:

	pets:
		- write
		- read

```

Output:

	tree := Tree{
		"users": {
			"write": {
				"owner": {}
			},
			"read": {}
		},
		"pets": {
			"write": {},
			"read": {},
		}
	}

Identifiers:

users:write
users:write:owner
users:read
pets:write
pets:read

Resolving:
*/
type Tree map[string]Tree

// IsFinalLeave returns true if this tree node is the final leave.
// That means that there are not further sub-permissions for this node.
func (t Tree) IsFinalLeave() bool {
	return len(t) == 0
}

func (t Tree) MergeFrom(another Tree) Tree {
	for key := range another {
		if _, ok := t[key]; !ok {
			t[key] = another[key]
		} else {
			t[key] = t[key].MergeFrom(another[key])
		}
	}

	return t
}

func (t Tree) Insert(path string) {
	split := strings.Split(path, ":")

	t.insert(split)
}

func (t Tree) InsertPath(path []string) {
	t.insert(path)
}

func (t Tree) insert(path []string) {
	if len(path) == 0 {
		return
	}

	tree, ok := t[path[0]]
	if !ok {
		tree = Tree{}

		t[path[0]] = tree
	}

	if len(path) > 1 {
		tree.insert(path[1:])
	}
}

func (t Tree) Resolve(perms []string) ([]string, error) {
	var set = make([]string, 0, len(perms))

	for _, p := range perms {
		split := strings.Split(p, ":")
		res := t.collect("", split)
		set = append(set, res...)
	}

	// FIXME(ppacher): remove any prefixes if the full sub-set is already included.
	// 				   i.e. remove "idm:users" from ["idm:users", "idm:users:read", "idm:users:write"]
	slices.Sort(set)
	set = slices.Compact(set)

	return set, nil
}

func (t Tree) collect(prefix string, p []string) []string {
	if len(p) == 0 {
		result := make([]string, 0, len(t))

		for k, child := range t {
			next_prefix := prefix + ":" + k

			if child.IsFinalLeave() {
				result = append(result, strings.TrimPrefix(next_prefix, ":"))
			} else {
				result = append(result, child.collect(next_prefix, nil)...)
			}
		}

		return result
	}

	children, ok := t[p[0]]
	if !ok {
		return []string{strings.TrimPrefix(prefix+":"+strings.Join(p, ":"), ":")}
	}

	return children.collect(prefix+":"+p[0], p[1:])
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface and adds support
// to parse a Tree node from an object (Tree) or string list.
func (tn *Tree) UnmarshalJSON(blob []byte) error {
	dec := json.NewDecoder(bytes.NewReader(blob))

	token, err := dec.Token()
	if err != nil {
		return err
	}

	delim, ok := token.(json.Delim)
	if !ok {
		return fmt.Errorf("unexpected token %T", token)
	}

	switch delim {
	case '{': // TreeNode
		var m map[string]Tree
		if err := json.Unmarshal(blob, &m); err != nil {
			return err
		}

		*tn = Tree(m)
	case '[': // Array
		var m []string
		if err := json.Unmarshal(blob, &m); err != nil {
			return err
		}

		*tn = make(Tree)
		for _, perm := range m {
			(*tn)[perm] = Tree{}
		}

	default:
		return fmt.Errorf("unexpected token %q", token)
	}

	return nil
}
