package fsx

import (
	"errors"
	"os"
	"path/filepath"
)

func ErrIsNotExist(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func IgnoreIsNotExist(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

func IgnoreIsExist(err error) error {
	if errors.Is(err, os.ErrExist) {
		return nil
	}

	return err
}

func AutoCached(path string, gen func() ([]byte, error)) (s []byte, err error) {
	if s, err = os.ReadFile(path); err == nil {
		return s, nil
	}

	if s, err = gen(); err != nil {
		return nil, err
	}

	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}

	if err = os.WriteFile(path, s, 0600); err != nil {
		return nil, err
	}

	return s, err
}

// IsRegularFile returns true IFF a non-directory file exists at the provided path.
func IsRegularFile(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}
