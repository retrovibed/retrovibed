package ddisc_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/stretchr/testify/require"
)

func TestLocate(t *testing.T) {
	t.Run("should derive a deterministic id from query and mimetype", func(t *testing.T) {
		l := ddisc.NewLocate("ubuntu", mimex.Video)
		require.Equal(t, "ubuntu", l.Query)
		require.Equal(t, mimex.Video, l.Mimetype)
		require.Equal(t, uuid.Nil.String(), l.KnownMediaID)
		require.NotEmpty(t, l.ID)

		same := ddisc.NewLocate("ubuntu", mimex.Video)
		require.Equal(t, l.ID, same.ID, "id should be deterministic from (query, mimetype)")
	})

	t.Run("should not duplicate a row when the same query is inserted twice", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		l := ddisc.NewLocate("ubuntu", mimex.Video)
		require.NoError(t, ddisc.LocateInsertWithDefaults(t.Context(), q, l).Scan(&l))
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_locate"))

		again := ddisc.NewLocate("ubuntu", mimex.Video)
		require.NoError(t, ddisc.LocateInsertWithDefaults(t.Context(), q, again).Scan(&again))
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_locate"), "re-requesting the same query+mimetype should not duplicate the row")
	})

	t.Run("should remain pending after being located, and drop out of pending once completed", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		l := ddisc.NewLocate("ubuntu", mimex.Video)
		require.NoError(t, ddisc.LocateInsertWithDefaults(t.Context(), q, l).Scan(&l))

		tid := uuid.Must(uuid.NewV4()).String()
		require.NoError(t, ddisc.LocateLocated(t.Context(), q, l.ID, tid).Scan(&l))
		require.Equal(t, tid, l.LocatedTorrentID)

		pending := ddisc.LocateSearchBuilder().Where(ddisc.LocateQueryPending())

		s := sqlx.Scan(ddisc.LocateSearch(t.Context(), q, pending))
		found := false
		for range s.Iter() {
			found = true
		}
		require.NoError(t, s.Err())
		require.True(t, found, "should still be pending after only being located")

		require.NoError(t, ddisc.LocateCompleted(t.Context(), q, l.ID).Scan(&l))

		s = sqlx.Scan(ddisc.LocateSearch(t.Context(), q, pending))
		found = false
		for range s.Iter() {
			found = true
		}
		require.NoError(t, s.Err())
		require.False(t, found, "should no longer be pending after completion")
	})

	t.Run("should reopen a completed locate when the same query is searched again", func(t *testing.T) {
		q := sqltestx.Metadatabase(t)

		l := ddisc.NewLocate("ubuntu", mimex.Video)
		require.NoError(t, ddisc.LocateInsertWithDefaults(t.Context(), q, l).Scan(&l))

		tid := uuid.Must(uuid.NewV4()).String()
		require.NoError(t, ddisc.LocateLocated(t.Context(), q, l.ID, tid).Scan(&l))
		require.NoError(t, ddisc.LocateCompleted(t.Context(), q, l.ID).Scan(&l))

		// simulates: the recommendation this locate pointed at was since
		// consumed and hard-deleted (library.RecommendationDeleteByID).
		// searching the identical query again must reopen the request
		// instead of silently handing back the stale, now-dangling
		// located_torrent_id.
		again := ddisc.NewLocate("ubuntu", mimex.Video)
		require.NoError(t, ddisc.LocateInsertWithDefaults(t.Context(), q, again).Scan(&again))
		require.Equal(t, l.ID, again.ID, "sanity: same deterministic id")

		pending := ddisc.LocateSearchBuilder().Where(ddisc.LocateQueryPending())
		s := sqlx.Scan(ddisc.LocateSearch(t.Context(), q, pending))
		found := false
		for range s.Iter() {
			found = true
		}
		require.NoError(t, s.Err())
		require.True(t, found, "re-searching a completed locate should reopen it for reprocessing")
		require.Equal(t, uuid.Max.String(), again.LocatedTorrentID, "stale located_torrent_id should be cleared, not handed back")
	})
}
