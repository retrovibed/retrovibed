package library_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestWatchHistoryJSONLEncode(t *testing.T) {
	t.Run("streams matching rows as jsonl", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		profileID := uuidx.WithSuffix(1)
		mediaID := uuidx.WithSuffix(2)
		otherProfileID := uuidx.WithSuffix(3)

		var a library.WatchHistory
		require.NoError(t, library.WatchHistoryInsertWithDefaults(ctx, q, library.WatchHistory{
			ID:        uuidx.WithSuffix(4),
			ProfileID: profileID,
			MediaID:   mediaID,
			Duration:  5 * time.Second,
		}).Scan(&a))

		var b library.WatchHistory
		require.NoError(t, library.WatchHistoryInsertWithDefaults(ctx, q, library.WatchHistory{
			ID:        uuidx.WithSuffix(5),
			ProfileID: profileID,
			MediaID:   mediaID,
			Duration:  10 * time.Second,
		}).Scan(&b))

		// a different profile's row must not leak into the export for profileID.
		var other library.WatchHistory
		require.NoError(t, library.WatchHistoryInsertWithDefaults(ctx, q, library.WatchHistory{
			ID:        uuidx.WithSuffix(6),
			ProfileID: otherProfileID,
			MediaID:   mediaID,
			Duration:  1 * time.Second,
		}).Scan(&other))

		var buf bytes.Buffer
		require.NoError(t, library.WatchHistoryJSONLEncode(
			ctx, q,
			library.WatchHistorySearchBuilder().Where(library.WatchHistoryQueryProfileID(profileID)),
			&buf,
		))

		dec := jsonl.NewDecoder(&buf)
		got := map[string]*library.WatchHistoryRecord{}
		for {
			rec := &library.WatchHistoryRecord{}
			err := dec.Decode(rec)
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
			got[rec.Id] = rec
		}

		require.Len(t, got, 2)
		require.Contains(t, got, a.ID)
		require.Contains(t, got, b.ID)
		require.Equal(t, mediaID, got[a.ID].MediaId)
		require.EqualValues(t, 5000, got[a.ID].Duration)
		require.Equal(t, mediaID, got[b.ID].MediaId)
		require.EqualValues(t, 10000, got[b.ID].Duration)
	})

	t.Run("no matching rows encodes nothing", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var buf bytes.Buffer
		require.NoError(t, library.WatchHistoryJSONLEncode(
			ctx, q,
			library.WatchHistorySearchBuilder().Where(library.WatchHistoryQueryProfileID(uuidx.WithSuffix(99))),
			&buf,
		))

		require.Empty(t, buf.String())
	})

	t.Run("propagates a bad builder as an error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var buf bytes.Buffer
		badBuilder := library.WatchHistorySearchBuilder().Where(squirrel.Expr("not_a_real_column = ?", "x"))
		require.Error(t, library.WatchHistoryJSONLEncode(ctx, q, badBuilder, &buf))
	})
}
