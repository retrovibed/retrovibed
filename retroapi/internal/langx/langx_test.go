package langx_test

import (
	"errors"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/internal/langx"
	"github.com/stretchr/testify/require"
)

func TestFirstNonZero(t *testing.T) {
	t.Run("returns first non-zero string", func(t *testing.T) {
		require.Equal(t, "a", langx.FirstNonZero("", "a", "b"))
	})

	t.Run("returns first non-zero int", func(t *testing.T) {
		require.Equal(t, 1, langx.FirstNonZero(0, 0, 1, 2))
	})

	t.Run("returns zero value when all are zero", func(t *testing.T) {
		require.Equal(t, "", langx.FirstNonZero("", "", ""))
	})

	t.Run("returns zero value for empty input", func(t *testing.T) {
		require.Equal(t, 0, langx.FirstNonZero[int]())
	})

	t.Run("returns first non-nil error", func(t *testing.T) {
		err1 := errors.New("first")
		err2 := errors.New("second")
		require.Equal(t, err1, langx.FirstNonZero(nil, err1, err2))
	})

	t.Run("returns nil when all errors are nil", func(t *testing.T) {
		require.Nil(t, langx.FirstNonZero[error](nil, nil))
	})

	t.Run("returns single non-zero value", func(t *testing.T) {
		require.Equal(t, "only", langx.FirstNonZero("only"))
	})
}
