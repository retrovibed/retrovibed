package daemons

import (
	"context"
	"fmt"
	"iter"
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
}

func (t fakeSearchPlugins) Search(ctx context.Context, mimetypes []string, query string) iterx.Seq[*ddiscapi.Import] {
	return fakeResultSeq(t)
}

func TestSearchQueueBackgroundRunPersistsFoundResults(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	var known library.Known
	require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
	require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

	require.NoError(t, ddisc.SearchQueueEnqueue(ctx, q, ddisc.SearchQueue{KnownMediaID: known.UID}).Scan(&ddisc.SearchQueue{}))

	id := int160.Random()
	magnet := fmt.Sprintf("magnet:?xt=urn:btih:%s", id.String())

	plugins := fakeSearchPlugins{results: []*ddiscapi.Import{{Magnet: magnet, Health: 10, Mimetype: mimex.Video}}}
	require.NoError(t, SearchQueueBackgroundRun(ctx, q, plugins))

	require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = ?", known.UID))
	require.EqualValues(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_search_queue WHERE known_media_id = ?", known.UID))
}

func TestSearchQueueBackgroundRunCooldownOnEmptyResults(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	var known library.Known
	require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
	require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

	require.NoError(t, ddisc.SearchQueueEnqueue(ctx, q, ddisc.SearchQueue{KnownMediaID: known.UID}).Scan(&ddisc.SearchQueue{}))

	require.NoError(t, SearchQueueBackgroundRun(ctx, q, fakeSearchPlugins{}))

	require.EqualValues(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = ?", known.UID))
	require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_search_queue WHERE known_media_id = ? AND attempts = 1", known.UID))
}
