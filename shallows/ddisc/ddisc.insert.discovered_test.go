package ddisc_test

import (
	"testing"
	"time"
	"unicode/utf8"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/stretchr/testify/require"
)

func TestDiscoveredInsertWithDefaults(t *testing.T) {
	t.Run("InsertWithDefaults", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()

		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
		require.NoError(t, err)

		d := ddisc.NewDiscovered(
			&id,
			ddisc.DiscoveredOptionIndex(true),
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionFromTorrentInfo(info),
		)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
	})

	t.Run("applies database defaults for unset columns", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()

		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
		require.NoError(t, err)

		d := ddisc.NewDiscovered(
			&id,
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionFromTorrentInfo(info),
		)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		require.WithinDuration(t, time.Now(), d.CreatedAt, time.Minute)
		require.WithinDuration(t, time.Now(), d.UpdatedAt, time.Minute)
		require.WithinDuration(t, time.Now(), d.NextCheckAt, time.Minute)
		require.True(t, d.TombstonedAt.Equal(timex.Inf()), "tombstoned_at should default to infinity")
		require.True(t, d.ReleasedAt.Equal(timex.Inf()), "released_at should default to infinity")
	})

	t.Run("round trips explicit column values", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()

		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
		require.NoError(t, err)

		d := ddisc.NewDiscovered(
			&id,
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionFromTorrentInfo(info),
		)
		d.Attempts = 3
		d.AudioBitrate = 320
		d.AudioDefaultLocale = "en-US"
		d.Collation = "alphabetical"
		d.Description = "a description"
		d.SubtitlesDefaultLocale = "es-ES"
		d.VideoResolution = "1920x1080"
		d.VideoRuntime = 90 * time.Minute

		expected := d
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

		require.Equal(t, expected.Attempts, d.Attempts)
		require.Equal(t, expected.AudioBitrate, d.AudioBitrate)
		require.Equal(t, expected.AudioDefaultLocale, d.AudioDefaultLocale)
		require.Equal(t, expected.Bytes, d.Bytes)
		require.Equal(t, expected.Collation, d.Collation)
		require.Equal(t, expected.Description, d.Description)
		require.Equal(t, expected.ID, d.ID)
		require.Equal(t, expected.Infohash, d.Infohash)
		require.Equal(t, expected.KnownMediaID, d.KnownMediaID)
		require.Equal(t, expected.Mimetype, d.Mimetype)
		require.Equal(t, expected.Partition, d.Partition)
		require.Equal(t, expected.SubtitlesDefaultLocale, d.SubtitlesDefaultLocale)
		require.Equal(t, expected.SyncUID, d.SyncUID)
		require.Equal(t, expected.Title, d.Title)
		require.Equal(t, expected.VideoResolution, d.VideoResolution)
		require.Equal(t, expected.VideoRuntime, d.VideoRuntime)
	})

	t.Run("on conflict only refreshes updated_at", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()

		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
		require.NoError(t, err)

		d := ddisc.NewDiscovered(
			&id,
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionFromTorrentInfo(info),
		)
		d.Title = "original title"
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
		firstUpdatedAt := d.UpdatedAt

		time.Sleep(10 * time.Millisecond)

		conflict := d
		conflict.Title = "changed title"
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, conflict).Scan(&conflict))

		require.Equal(t, "original title", conflict.Title, "conflict clause only refreshes updated_at, other columns are left untouched")
		require.True(t, conflict.UpdatedAt.After(firstUpdatedAt))
	})

	t.Run("rejects corrupted titles from syncing instead of failing the insert", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tmpdir := t.TempDir()

		q := sqltestx.Metadatabase(t)

		id := int160.Random()
		info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
		require.NoError(t, err)
		info.Name = string([]byte{0xff, 0xfe, 0x00, 0x80})

		d := ddisc.NewDiscovered(
			&id,
			ddisc.DiscoveredOptionMimetype(mimex.Binary),
			ddisc.DiscoveredOptionFromTorrentInfo(info),
			ddisc.DiscoveredOptionDetectCorrupted,
		)
		require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
		require.True(t, utf8.ValidString(d.Title))
		require.Equal(t, uuid.Nil.String(), d.SyncUID)
	})

	t.Run("round trips hidden_at infinity sentinels", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			ts   time.Time
		}{
			{"positive infinity", timex.Inf()},
			{"negative infinity", timex.NegInf()},
		} {
			t.Run(tc.name, func(t *testing.T) {
				ctx, done := testx.Context(t)
				defer done()

				tmpdir := t.TempDir()

				q := sqltestx.Metadatabase(t)

				id := int160.Random()
				info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
				require.NoError(t, err)

				d := ddisc.NewDiscovered(
					&id,
					ddisc.DiscoveredOptionMimetype(mimex.Binary),
					ddisc.DiscoveredOptionFromTorrentInfo(info),
				)
				d.HiddenAt = tc.ts
				require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))
				require.True(t, d.HiddenAt.Equal(tc.ts))
			})
		}
	})
}
