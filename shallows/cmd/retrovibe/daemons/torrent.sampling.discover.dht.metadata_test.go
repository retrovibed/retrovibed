package daemons_test

import (
	"net/netip"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestDiscoverDHTMetadata(t *testing.T) {
	t.Run("downloads metadata for a reachable peer and records it as discovered, removing the unknown hash", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		seederdir := t.TempDir()
		seederinfo, _, err := torrenttest.Random(seederdir, 16*1024)
		require.NoError(t, err)

		seeder := torrenttestx.QuickClient(t)
		defer seeder.Close()

		seedermd, err := torrent.NewFromInfo(seederinfo, torrent.OptionStorage(storage.NewFile(seederdir)))
		require.NoError(t, err)
		infohash := seedermd.ID

		_, _, err = seeder.Start(seedermd)
		require.NoError(t, err)

		addr := torrentx.ClientAddress(seeder)
		require.True(t, addr.IsValid(), "expected seeder to have a dialable listen address")

		unk := tracking.NewUnknownHash(
			infohash,
			tracking.OptionUnknownHashPeer(int160.Random(), addr),
			func(uh *tracking.UnknownHash) { uh.NextCheck = time.Now().Add(-time.Minute) },
		)
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, unk).Scan(&unk))

		tclient := torrenttestx.QuickClient(t)
		defer tclient.Close()

		notblocked := func(k []byte) bool { return false }

		go func() {
			_ = daemons.DiscoverDHTMetadata(ctx, 1, q, tclient, notblocked)
		}()

		expectedID := torrentx.HashUID(&infohash)

		require.Eventually(t, func() bool {
			var disc ddisc.Discovered
			return ddisc.DiscoveredFindByID(ctx, q, expectedID).Scan(&disc) == nil
		}, 20*time.Second, 200*time.Millisecond, "expected metadata to be discovered")

		var disc ddisc.Discovered
		require.NoError(t, ddisc.DiscoveredFindByID(ctx, q, expectedID).Scan(&disc))
		require.Equal(t, uuid.Max.String(), disc.KnownMediaID, "expected discovered record to be indexed")

		sql, args, err := tracking.UnknownSearchBuilder().RemoveColumns().Columns("COUNT(*)").Where(squirrel.Eq{"id": expectedID}).ToSql()
		require.NoError(t, err)
		require.Eventually(t, func() bool {
			return sqltestx.Count(t, q, sql, args...) == 0
		}, 10*time.Second, 100*time.Millisecond, "expected unknown hash record to be removed")
	})

	t.Run("respects the blocked filter and records the metadata as unindexed", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		seederdir := t.TempDir()
		seederinfo, _, err := torrenttest.Random(seederdir, 16*1024)
		require.NoError(t, err)

		seeder := torrenttestx.QuickClient(t)
		defer seeder.Close()

		seedermd, err := torrent.NewFromInfo(seederinfo, torrent.OptionStorage(storage.NewFile(seederdir)))
		require.NoError(t, err)
		infohash := seedermd.ID

		_, _, err = seeder.Start(seedermd)
		require.NoError(t, err)

		addr := torrentx.ClientAddress(seeder)
		require.True(t, addr.IsValid(), "expected seeder to have a dialable listen address")

		unk := tracking.NewUnknownHash(
			infohash,
			tracking.OptionUnknownHashPeer(int160.Random(), addr),
			func(uh *tracking.UnknownHash) { uh.NextCheck = time.Now().Add(-time.Minute) },
		)
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, unk).Scan(&unk))

		tclient := torrenttestx.QuickClient(t)
		defer tclient.Close()

		blockall := func(k []byte) bool { return true }

		go func() {
			_ = daemons.DiscoverDHTMetadata(ctx, 1, q, tclient, blockall)
		}()

		expectedID := torrentx.HashUID(&infohash)

		require.Eventually(t, func() bool {
			var disc ddisc.Discovered
			return ddisc.DiscoveredFindByID(ctx, q, expectedID).Scan(&disc) == nil
		}, 20*time.Second, 200*time.Millisecond, "expected metadata to be discovered")

		var disc ddisc.Discovered
		require.NoError(t, ddisc.DiscoveredFindByID(ctx, q, expectedID).Scan(&disc))
		require.Equal(t, uuid.Nil.String(), disc.KnownMediaID, "expected blocked discovered record to remain unindexed")
	})

	t.Run("sanitizes a corrupted (invalid utf8) torrent name instead of failing the insert", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		seederdir := t.TempDir()
		seederinfo, _, err := torrenttest.Random(seederdir, 16*bytesx.KiB)
		require.NoError(t, err)
		seederinfo.Name = string([]byte{0xff, 0xfe, 0x00, 0x80})

		seeder := torrenttestx.QuickClient(t)
		defer seeder.Close()

		seedermd, err := torrent.NewFromInfo(seederinfo, torrent.OptionStorage(storage.NewFile(seederdir)))
		require.NoError(t, err)
		infohash := seedermd.ID

		_, _, err = seeder.Start(seedermd)
		require.NoError(t, err)

		addr := torrentx.ClientAddress(seeder)
		require.True(t, addr.IsValid(), "expected seeder to have a dialable listen address")

		unk := tracking.NewUnknownHash(
			infohash,
			tracking.OptionUnknownHashPeer(int160.Random(), addr),
			func(uh *tracking.UnknownHash) { uh.NextCheck = time.Now().Add(-time.Minute) },
		)
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, unk).Scan(&unk))

		tclient := torrenttestx.QuickClient(t)
		defer tclient.Close()

		notblocked := func(k []byte) bool { return false }

		go func() {
			_ = daemons.DiscoverDHTMetadata(ctx, 1, q, tclient, notblocked)
		}()

		expectedID := torrentx.HashUID(&infohash)

		require.Eventually(t, func() bool {
			var disc ddisc.Discovered
			return ddisc.DiscoveredFindByID(ctx, q, expectedID).Scan(&disc) == nil
		}, 20*time.Second, 200*time.Millisecond, "expected metadata to be discovered despite the corrupted name")

		var disc ddisc.Discovered
		require.NoError(t, ddisc.DiscoveredFindByID(ctx, q, expectedID).Scan(&disc))
		require.True(t, utf8.ValidString(disc.Title), "expected corrupted title to be sanitized to valid utf8")
		require.Equal(t, uuid.Nil.String(), disc.SyncUID, "expected corrupted record to be excluded from sync")

		sql, args, err := tracking.UnknownSearchBuilder().RemoveColumns().Columns("COUNT(*)").Where(squirrel.Eq{"id": expectedID}).ToSql()
		require.NoError(t, err)
		require.Eventually(t, func() bool {
			return sqltestx.Count(t, q, sql, args...) == 0
		}, 10*time.Second, 100*time.Millisecond, "expected unknown hash record to be removed")
	})

	t.Run("cools down an unknown hash whose peer never responds instead of discarding it", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		infohash := int160.Random()
		originalNextCheck := time.Now().Add(-time.Minute)

		unk := tracking.NewUnknownHash(
			infohash,
			// zero port (with a valid, but unspecified, IP to satisfy the
			// NOT NULL column) makes the daemon fall back to torrent.TuneNoop,
			// so it never finds a peer and the per-attempt timeout fires.
			tracking.OptionUnknownHashPeer(int160.Random(), netip.AddrPort{}),
			func(uh *tracking.UnknownHash) { uh.NextCheck = originalNextCheck },
		)
		require.NoError(t, tracking.UnknownHashInsertWithDefaults(ctx, q, unk).Scan(&unk))

		tclient := torrenttestx.QuickClient(t)
		defer tclient.Close()

		notblocked := func(k []byte) bool { return false }

		go func() {
			_ = daemons.DiscoverDHTMetadata(ctx, 1, q, tclient, notblocked)
		}()

		expectedID := torrentx.HashUID(&infohash)

		countsql, countargs, err := tracking.UnknownSearchBuilder().RemoveColumns().Columns("COUNT(*)").Where(squirrel.Eq{"id": expectedID}).ToSql()
		require.NoError(t, err)
		attemptssql, attemptsargs, err := tracking.UnknownSearchBuilder().RemoveColumns().Columns("attempts").Where(squirrel.Eq{"id": expectedID}).ToSql()
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			return sqltestx.Count(t, q, countsql, countargs...) == 1 && sqltestx.Count(t, q, attemptssql, attemptsargs...) >= 1
		}, 30*time.Second, 200*time.Millisecond, "expected unknown hash to be cooled down after a peer timeout")

		discsql, discargs, err := ddisc.DiscoveredSearchBuilder().RemoveColumns().Columns("COUNT(*)").Where(squirrel.Eq{"id": expectedID}).ToSql()
		require.NoError(t, err)
		require.Equal(t, 0, sqltestx.Count(t, q, discsql, discargs...), "expected no discovered record for an undownloadable hash")
	})
}
