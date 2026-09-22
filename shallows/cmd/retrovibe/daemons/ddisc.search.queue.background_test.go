package daemons

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

type fakeResultSeq struct {
	results []*ddiscapi.Import
}

func (t fakeResultSeq) Each(ctx context.Context) iter.Seq[*ddiscapi.Import] {
	return iterx.From(t.results...)
}

func (t fakeResultSeq) Err() error { return nil }

type fakeSearchPlugins struct {
	results []*ddiscapi.Import
	// public, when set, records the public argument of the last Search call.
	public *bool
}

func (t fakeSearchPlugins) Search(ctx context.Context, mimetypes []string, query string, adult, public bool) iterx.Seq[*ddiscapi.Import] {
	if t.public != nil {
		*t.public = public
	}
	return fakeResultSeq{results: t.results}
}

func TestSearchQueueBackgroundRun(t *testing.T) {
	t.Run("should persist found results without resolving their real infohash", func(t *testing.T) {
		// This test is important. we dont want to spam plugins with data fetches every time we do a search.
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		require.NoError(t, ddisc.SearchQueueEnqueue(ctx, q, ddisc.SearchQueue{KnownMediaID: known.UID}).Scan(&ddisc.SearchQueue{}))

		id := int160.Random()
		magnet := fmt.Sprintf("magnet:?xt=urn:btih:%s", id.String())

		plugins := fakeSearchPlugins{results: []*ddiscapi.Import{{Uri: magnet, Uritype: mimex.Magnet, Health: 10, Mimetype: mimex.Video, Title: known.Title}}}
		require.NoError(t, SearchQueueBackgroundRun(ctx, q, plugins, ddisc.UnimplementedStrategy{}, library.QueryCleanerNoop()))

		require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = ?", known.UID))
		require.EqualValues(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_search_queue WHERE known_media_id = ?", known.UID))
		// a magnet's real infohash is known without any fetch - parsed from the uri itself.
		require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = ? AND infohash = ?", known.UID, id.Bytes()))
	})

	t.Run("should cooldown when no results are found", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		require.NoError(t, ddisc.SearchQueueEnqueue(ctx, q, ddisc.SearchQueue{KnownMediaID: known.UID}).Scan(&ddisc.SearchQueue{}))

		require.NoError(t, SearchQueueBackgroundRun(ctx, q, fakeSearchPlugins{}, ddisc.UnimplementedStrategy{}, library.QueryCleanerNoop()))

		require.EqualValues(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = ?", known.UID))
		require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_search_queue WHERE known_media_id = ? AND attempts = 1", known.UID))
	})

	t.Run("should cooldown, not delete, the entry when every candidate is policy-rejected", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		require.NoError(t, ddisc.SearchQueueEnqueue(ctx, q, ddisc.SearchQueue{KnownMediaID: known.UID}).Scan(&ddisc.SearchQueue{}))

		// title matches known.Title (so it survives the title filter) but is
		// a CAM-class release, which ddisc.DefaultPolicy hard-rejects. It
		// must never be persisted, and the entry must cooldown same as a
		// clean not-found, not be treated as resolved.
		plugins := fakeSearchPlugins{results: []*ddiscapi.Import{
			{Uri: "https://tracker.example/rejected.torrent", Uritype: mimex.Bittorrent, Title: known.Title + " CAM"},
		}}
		require.NoError(t, SearchQueueBackgroundRun(ctx, q, plugins, ddisc.UnimplementedStrategy{}, library.QueryCleanerNoop()))

		require.EqualValues(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = ?", known.UID))
		require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_search_queue WHERE known_media_id = ? AND attempts = 1", known.UID))
	})

	t.Run("should only ask search plugins for public results", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		require.NoError(t, ddisc.SearchQueueEnqueue(ctx, q, ddisc.SearchQueue{KnownMediaID: known.UID}).Scan(&ddisc.SearchQueue{}))

		var public bool
		require.NoError(t, SearchQueueBackgroundRun(ctx, q, fakeSearchPlugins{public: &public}, ddisc.UnimplementedStrategy{}, library.QueryCleanerNoop()))
		require.True(t, public)
	})

	t.Run("should never fetch a candidate's torrent file during a drain", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		require.NoError(t, ddisc.SearchQueueEnqueue(ctx, q, ddisc.SearchQueue{KnownMediaID: known.UID}).Scan(&ddisc.SearchQueue{}))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.False(t, true, "the search queue drain must never fetch a candidate's torrent file")
		}))
		defer srv.Close()

		plugins := fakeSearchPlugins{results: []*ddiscapi.Import{
			{Uri: srv.URL + "/one.torrent", Uritype: mimex.Bittorrent, Health: 10, Mimetype: mimex.Video, Title: known.Title},
			{Uri: srv.URL + "/two.torrent", Uritype: mimex.Bittorrent, Health: 10, Mimetype: mimex.Video, Title: known.Title},
		}}
		require.NoError(t, SearchQueueBackgroundRun(ctx, q, plugins, ddisc.UnimplementedStrategy{}, library.QueryCleanerNoop()))

		// both candidates are still persisted - just with their placeholder
		// infohash, unresolved - and defaulted private until selected for download.
		require.EqualValues(t, 2, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = ? AND private = true", known.UID))
	})
}
