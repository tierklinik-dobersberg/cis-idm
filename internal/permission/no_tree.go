package permission

import "slices"

type NoTree struct{}

func (NoTree) Resolve(perms []string) ([]string, error) { 
	set := make([]string, 0, len(perms))
	copy(set, perms)

	slices.Sort(set)
	set = slices.Compact(set)

	return set, nil
}
