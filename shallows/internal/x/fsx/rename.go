package fsx

import (
	"errors"
)

// will rename a file if possible, if the vfs doesnt support renaming
// it'll return ErrUnsupported
func Rename(vfs Virtual, oldpath, newpath string) error {
	type renameable interface {
		Rename(oldpath, newpath string) error
	}

	if tmp, ok := vfs.(renameable); ok {
		return tmp.Rename(oldpath, newpath)
	}

	return errors.ErrUnsupported
}
