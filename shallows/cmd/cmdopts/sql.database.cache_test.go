package cmdopts

import (
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/stretchr/testify/require"
)

func TestDatabaseCache(t *testing.T) {
	t.Run("opens and migrates the path privately", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db, err := DatabaseCache(ctx, filepath.Join(t.TempDir(), "other.db"))
		require.NoError(t, err)
		defer db.Close()

		var n int
		require.NoError(t, db.QueryRowContext(ctx, "SELECT COUNT(*) FROM cache.library_known_media").Scan(&n))
	})
}
