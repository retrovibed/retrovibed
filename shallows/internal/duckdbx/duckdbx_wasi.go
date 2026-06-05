//go:build wasip1

package duckdbx

func ErrUniqueConstraintViolation(err error) error {
	return nil
}
