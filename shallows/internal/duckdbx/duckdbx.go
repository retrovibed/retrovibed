//go:build !wasip1

package duckdbx

import (
	"github.com/duckdb/duckdb-go/v2"
)

// checks if the error is a unique constraint violation.
func ErrUniqueConstraintViolation(err error) error {
	switch cause := err.(type) {
	case *duckdb.Error:
		switch cause.Type {
		case duckdb.ErrorTypeConstraint:
			return err
		default:
			return nil
		}
	default:
	}

	return nil
}
