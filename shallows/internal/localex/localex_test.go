package localex_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/localex"
	"github.com/stretchr/testify/require"
)

func TestFirstDefined(t *testing.T) {
	t.Run("returns first non-blank locale", func(t *testing.T) {
		require.Equal(t, "en", localex.FirstDefined("", "en", "fr"))
	})

	t.Run("skips undetermined locale", func(t *testing.T) {
		require.Equal(t, "en", localex.FirstDefined("und", "en"))
	})

	t.Run("returns undetermined when all blank", func(t *testing.T) {
		require.Equal(t, "und", localex.FirstDefined("", "", ""))
	})

	t.Run("returns undetermined for empty input", func(t *testing.T) {
		require.Equal(t, "und", localex.FirstDefined())
	})

	t.Run("returns undetermined when all undetermined", func(t *testing.T) {
		require.Equal(t, "und", localex.FirstDefined("und", "und"))
	})

	t.Run("returns single defined locale", func(t *testing.T) {
		require.Equal(t, "fr", localex.FirstDefined("fr"))
	})
}
