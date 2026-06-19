package daemons_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestBEP0051SamplerSnapshot(t *testing.T) {
	t.Run("returns zero hashes when nothing qualifies", func(t *testing.T) {
		_, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cachepath := filepath.Join(t.TempDir(), "sample.cache")

		s := daemons.NewSampler(q, time.Hour, cachepath)
		ttl, total, sample := s.Snapshot(128)

		require.Equal(t, uint(time.Hour/time.Second), ttl)
		require.Zero(t, total)
		require.Empty(t, sample)
	})

	t.Run("returns the infohashes of public seeding metadata", func(t *testing.T) {
		_, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cachepath := filepath.Join(t.TempDir(), "sample.cache")

		one := tracking.NewMetadata(langx.Autoptr(int160.Random()))
		one.Private = false
		one.Seeding = true
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, one).Scan(&one))

		two := tracking.NewMetadata(langx.Autoptr(int160.Random()))
		two.Private = false
		two.Seeding = true
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, two).Scan(&two))

		// noise that must never appear in the sample.
		priv := tracking.NewMetadata(langx.Autoptr(int160.Random()))
		priv.Private = true
		priv.Seeding = true
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, priv).Scan(&priv))

		notSeeding := tracking.NewMetadata(langx.Autoptr(int160.Random()))
		notSeeding.Private = false
		notSeeding.Seeding = false
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, notSeeding).Scan(&notSeeding))

		s := daemons.NewSampler(q, time.Hour, cachepath)
		ttl, total, sample := s.Snapshot(128)

		require.Equal(t, uint(time.Hour/time.Second), ttl)
		require.EqualValues(t, 2, total)
		require.Len(t, sample, 2*20)

		expected := map[string]struct{}{
			string(one.Infohash): {},
			string(two.Infohash): {},
		}
		for i := 0; i < len(sample)/20; i++ {
			_, ok := expected[string(sample[i*20:(i+1)*20])]
			require.True(t, ok, "unexpected infohash in sample")
		}
	})

	t.Run("caps the sample at max but reports the full total", func(t *testing.T) {
		_, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cachepath := filepath.Join(t.TempDir(), "sample.cache")

		for i := 0; i < 5; i++ {
			m := tracking.NewMetadata(langx.Autoptr(int160.Random()))
			m.Private = false
			m.Seeding = true
			require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, m).Scan(&m))
		}

		s := daemons.NewSampler(q, time.Hour, cachepath)
		_, total, sample := s.Snapshot(2)

		require.EqualValues(t, 5, total)
		require.Len(t, sample, 2*20)
	})

	t.Run("serves a cached sample until the ttl elapses", func(t *testing.T) {
		_, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		cachepath := filepath.Join(t.TempDir(), "sample.cache")

		first := tracking.NewMetadata(langx.Autoptr(int160.Random()))
		first.Private = false
		first.Seeding = true
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, first).Scan(&first))

		ttl := 100 * time.Millisecond
		s := daemons.NewSampler(q, ttl, cachepath)

		_, total, _ := s.Snapshot(128)
		require.EqualValues(t, 1, total)

		// a new qualifying row appears, but the cached sample should still
		// be served until the ttl elapses.
		second := tracking.NewMetadata(langx.Autoptr(int160.Random()))
		second.Private = false
		second.Seeding = true
		require.NoError(t, tracking.MetadataInsertWithDefaults(t.Context(), q, second).Scan(&second))

		_, total, _ = s.Snapshot(128)
		require.EqualValues(t, 1, total, "expected cached sample to be served within the ttl window")

		require.Eventually(t, func() bool {
			_, total, _ := s.Snapshot(128)
			return total == 2
		}, time.Second, 10*time.Millisecond, "expected sample to refresh once the ttl elapsed")
	})
}
