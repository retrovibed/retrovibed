package daemons

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
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
}

func (t fakeSearchPlugins) Search(ctx context.Context, mimetypes []string, query string, adult, public bool) iterx.Seq[*ddiscapi.Import] {
	return fakeResultSeq(t)
}

func TestSearchQueueBackgroundRun(t *testing.T) {
	t.Run("should persist found results with the real resolved infohash", func(t *testing.T) {
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
		importer := tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir()))
		require.NoError(t, SearchQueueBackgroundRun(ctx, q, importer, plugins, ddisc.UnimplementedStrategy{}, library.QueryCleanerNoop()))

		require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = ?", known.UID))
		require.EqualValues(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_search_queue WHERE known_media_id = ?", known.UID))
		require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = ? AND infohash = ?", known.UID, id.Bytes()), "the real infohash resolved from the magnet should be persisted, not a placeholder")
	})

	t.Run("should cooldown when no results are found", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		require.NoError(t, ddisc.SearchQueueEnqueue(ctx, q, ddisc.SearchQueue{KnownMediaID: known.UID}).Scan(&ddisc.SearchQueue{}))

		importer := tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir()))
		require.NoError(t, SearchQueueBackgroundRun(ctx, q, importer, fakeSearchPlugins{}, ddisc.UnimplementedStrategy{}, library.QueryCleanerNoop()))

		require.EqualValues(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = ?", known.UID))
		require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_search_queue WHERE known_media_id = ? AND attempts = 1", known.UID))
	})
}
