package fsx

import (
	"errors"
	"os"
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
