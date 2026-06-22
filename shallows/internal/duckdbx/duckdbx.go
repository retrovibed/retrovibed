//go:build !wasip1

package duckdbx

import (
	"errors"

	"github.com/duckdb/duckdb-go/v2"
)

// checks if the error is a unique constraint violation.
func ErrUniqueConstraintViolation(err error) error {
	if cause, ok := errors.AsType[*duckdb.Error](err); ok {
		switch cause.Type {
		case duckdb.ErrorTypeConstraint:
			return err
		default:
			return nil
		}
	}

	return nil
}
