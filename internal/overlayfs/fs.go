package overlayfs

import (
	"errors"
	"io/fs"
)

type FS struct {
	fsList []fs.FS
}

func NewFS(fileSystems ...fs.FS) *FS {
	return &FS{
		fsList: fileSystems,
	}
}

func (root *FS) Open(path string) (fs.File, error) {
	for _, child := range root.fsList {
		file, err := child.Open(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}

			return nil, err
		}

		return file, nil
	}

	return nil, fs.ErrNotExist
}

// Compile time check.
var _ fs.FS = (*FS)(nil)
