package duckdbx_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/duckdb/duckdb-go/v2"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/stretchr/testify/require"
)

func TestErrUniqueConstraintViolation(t *testing.T) {
	constraintErr := &duckdb.Error{Type: duckdb.ErrorTypeConstraint, Msg: "Constraint Error: CHECK constraint failed"}

	t.Run("nil error returns nil", func(t *testing.T) {
		require.NoError(t, duckdbx.ErrUniqueConstraintViolation(nil))
	})

	t.Run("direct duckdb constraint error is detected", func(t *testing.T) {
		require.Equal(t, constraintErr, duckdbx.ErrUniqueConstraintViolation(constraintErr))
	})

	t.Run("duckdb error of a different type is ignored", func(t *testing.T) {
		other := &duckdb.Error{Type: duckdb.ErrorTypeSyntax, Msg: "Syntax Error: ..."}
		require.NoError(t, duckdbx.ErrUniqueConstraintViolation(other))
	})

	t.Run("unrelated error is ignored", func(t *testing.T) {
		require.NoError(t, duckdbx.ErrUniqueConstraintViolation(errors.New("boom")))
	})

	t.Run("errorsx.Wrap wrapped constraint error is detected", func(t *testing.T) {
		wrapped := errorsx.Wrap(constraintErr, "unable to insert published content")
		require.Equal(t, wrapped, duckdbx.ErrUniqueConstraintViolation(wrapped))
	})

	t.Run("errorsx.Wrap wrapped unrelated error is ignored", func(t *testing.T) {
		wrapped := errorsx.Wrap(errors.New("boom"), "unable to insert published content")
		require.NoError(t, duckdbx.ErrUniqueConstraintViolation(wrapped))
	})

	t.Run("errors.Join containing the constraint error is detected", func(t *testing.T) {
		joined := errors.Join(errors.New("rows close failed"), constraintErr)
		require.Equal(t, joined, duckdbx.ErrUniqueConstraintViolation(joined))
	})

	t.Run("errors.Join without a constraint error is ignored", func(t *testing.T) {
		joined := errors.Join(errors.New("boom"), errors.New("bang"))
		require.NoError(t, duckdbx.ErrUniqueConstraintViolation(joined))
	})

	t.Run("errors.Join wrapped further by errorsx.Wrap is detected", func(t *testing.T) {
		joined := errors.Join(errors.New("rows close failed"), constraintErr)
		wrapped := errorsx.Wrap(joined, "unable to insert published content")
		require.Equal(t, wrapped, duckdbx.ErrUniqueConstraintViolation(wrapped))
	})

	t.Run("fmt.Errorf %w wrapped constraint error is detected", func(t *testing.T) {
		wrapped := fmt.Errorf("unable to insert published content: %w", constraintErr)
		require.Equal(t, wrapped, duckdbx.ErrUniqueConstraintViolation(wrapped))
	})

	t.Run("fmt.Errorf %w wrapped unrelated error is ignored", func(t *testing.T) {
		wrapped := fmt.Errorf("unable to insert published content: %w", errors.New("boom"))
		require.NoError(t, duckdbx.ErrUniqueConstraintViolation(wrapped))
	})

	t.Run("nested fmt.Errorf %w chain is detected", func(t *testing.T) {
		inner := fmt.Errorf("scan failed: %w", constraintErr)
		outer := fmt.Errorf("unable to insert published content: %w", inner)
		require.Equal(t, outer, duckdbx.ErrUniqueConstraintViolation(outer))
	})
}
