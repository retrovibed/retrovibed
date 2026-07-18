package daemons

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestPublishDiscoveredMediaOne(t *testing.T) {
	t.Run("publishes a new discovered record for known media on a public torrent", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var tmd tracking.Metadata
		require.NoError(t, testx.Fake(&tmd, tracking.MetadataOptionTestDefaults))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, tmd).Scan(&tmd))

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		var lmd library.Metadata
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults, library.MetadataOptionKnownMediaID(known.UID), library.MetadataOptionTorrentID(tmd.ID)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		require.NoError(t, publishDiscoveredMediaOne(ctx, q, lmd))

		id := int160.FromBytes(tmd.Infohash)
		expected := ddisc.NewDiscoveredFromKnown(id, "", known)

		var disc ddisc.Discovered
		require.NoError(t, ddisc.DiscoveredFindByID(ctx, q, expected.ID).Scan(&disc))
		require.Equal(t, tmd.Infohash, disc.Infohash)
		require.Equal(t, known.UID, disc.KnownMediaID)
		require.Equal(t, known.Title, disc.Title)
	})

	t.Run("skips records that were already published", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var tmd tracking.Metadata
		require.NoError(t, testx.Fake(&tmd, tracking.MetadataOptionTestDefaults))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, tmd).Scan(&tmd))

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		var lmd library.Metadata
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults, library.MetadataOptionKnownMediaID(known.UID), library.MetadataOptionTorrentID(tmd.ID)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		id := int160.FromBytes(tmd.Infohash)
		existing := ddisc.NewDiscoveredFromKnown(id, "", known)
		existing.Title = "pre-existing title"
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, existing).Scan(&existing))

		require.NoError(t, publishDiscoveredMediaOne(ctx, q, lmd))

		var disc ddisc.Discovered
		require.NoError(t, ddisc.DiscoveredFindByID(ctx, q, existing.ID).Scan(&disc))
		require.Equal(t, "pre-existing title", disc.Title, "an already published record should not be overwritten")
	})

	t.Run("publishes private torrents locally but marks them private", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var tmd tracking.Metadata
		require.NoError(t, testx.Fake(&tmd, tracking.MetadataOptionTestDefaults))
		tmd.Private = true
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, tmd).Scan(&tmd))

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		var lmd library.Metadata
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults, library.MetadataOptionKnownMediaID(known.UID), library.MetadataOptionTorrentID(tmd.ID)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		require.NoError(t, publishDiscoveredMediaOne(ctx, q, lmd))

		id := int160.FromBytes(tmd.Infohash)
		expected := ddisc.NewDiscoveredFromKnown(id, "", known)

		var disc ddisc.Discovered
		require.NoError(t, ddisc.DiscoveredFindByID(ctx, q, expected.ID).Scan(&disc), "private torrents should still be recorded locally")
		require.True(t, disc.Private, "private torrents must be marked private so they're excluded from peer sync")
	})

	t.Run("returns an error when the torrent metadata cannot be found", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var lmd library.Metadata
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults, library.MetadataOptionTorrentID(uuid.Must(uuid.NewV4()).String())))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		require.Error(t, publishDiscoveredMediaOne(ctx, q, lmd))
	})

	t.Run("returns an error when the known media cannot be found", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var tmd tracking.Metadata
		require.NoError(t, testx.Fake(&tmd, tracking.MetadataOptionTestDefaults))
		require.NoError(t, tracking.MetadataInsertWithDefaults(ctx, q, tmd).Scan(&tmd))

		var lmd library.Metadata
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults, library.MetadataOptionKnownMediaID(uuid.Must(uuid.NewV4()).String()), library.MetadataOptionTorrentID(tmd.ID)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		require.Error(t, publishDiscoveredMediaOne(ctx, q, lmd))
	})
}
