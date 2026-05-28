package daemons_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestTorrentMetadataIdentify(t *testing.T) {
	t.Run("assigns known media id when description matches", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "The Grand Budapest Hotel"
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		lmd := tracking.NewMetadata(
			langx.Autoptr(int160.Random()),
			tracking.MetadataOptionDescription("The Grand Budapest Hotel"),
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		require.NoError(t, daemons.IdentifyTorrentMedia(ctx, q, library.QueryCleanerNoop()))

		require.NoError(t, tracking.MetadataFindByID(ctx, q, lmd.ID).Scan(&lmd))
		require.Equal(t, known.UID, lmd.KnownMediaID)
	})

	t.Run("no match leaves known media id unchanged", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		lmd := tracking.NewMetadata(
			langx.Autoptr(int160.Random()),
			tracking.MetadataOptionDescription("xyzzy completely unknown title 99999"),
		)
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		require.NoError(t, daemons.IdentifyTorrentMedia(ctx, q, library.QueryCleanerNoop()))

		require.NoError(t, tracking.MetadataFindByID(ctx, q, lmd.ID).Scan(&lmd))
		require.Equal(t, uuid.Max.String(), lmd.KnownMediaID)
	})

	t.Run("empty metadata table returns no error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, daemons.IdentifyTorrentMedia(ctx, q, library.QueryCleanerNoop()))
	})
}
