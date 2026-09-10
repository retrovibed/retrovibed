package daemons_test

import (
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepareDefaultFeeds(t *testing.T) {
	t.Run("generate default feeds", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "default.feeds.json")
		q := sqltestx.Metadatabase(t)
		require.False(t, fsx.Exists(path))
		require.EqualValues(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))

		require.NoError(t, daemons.PrepareDefaultFeeds(t.Context(), q, dir))

		assert.True(t, fsx.Exists(path))
		assert.EqualValues(t, 2, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))
	})
}
