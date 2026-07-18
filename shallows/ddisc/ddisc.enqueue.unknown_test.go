package ddisc_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestEnqueueUnknownWrapFound(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	tmpdir := t.TempDir()
	q := sqltestx.Metadatabase(t)

	kid := uuid.Must(uuid.NewV4()).String()

	id := int160.Random()
	info, _, err := torrenttest.Random(tmpdir, 128*bytesx.KiB)
	require.NoError(t, err)

	d := ddisc.NewDiscovered(
		&id,
		ddisc.DiscoveredOptionKnownMedia(kid),
		ddisc.DiscoveredOptionMimetype(mimex.Binary),
		ddisc.DiscoveredOptionFromTorrentInfo(info),
		ddisc.DiscoveredOptionAutoMagnet,
	)
	require.NoError(t, ddisc.DiscoveredInsertWithDefaults(ctx, q, d).Scan(&d))

	seq := ddisc.EnqueueUnknown().Wrap(q, kid, ddisc.FindMedia(q, kid))

	count := 0
	for range seq.Each(ctx) {
		count++
	}
	require.NoError(t, seq.Err())
	require.Equal(t, 1, count)
	require.EqualValues(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_search_queue WHERE known_media_id = ?", kid))
}

func TestEnqueueUnknownWrapNotFoundEnqueues(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	kid := uuid.Must(uuid.NewV4()).String()

	seq := ddisc.EnqueueUnknown().Wrap(q, kid, ddisc.FindMedia(q, kid))

	count := 0
	for range seq.Each(ctx) {
		count++
	}
	require.NoError(t, seq.Err())
	require.Equal(t, 0, count)
	require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_search_queue WHERE known_media_id = ?", kid))

	// re-wrapping the same still-missing kid is idempotent, not a duplicate row.
	seq2 := ddisc.EnqueueUnknown().Wrap(q, kid, ddisc.FindMedia(q, kid))
	for range seq2.Each(ctx) {
	}
	require.NoError(t, seq2.Err())
	require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_search_queue WHERE known_media_id = ?", kid))
}

func TestEnqueueUnknownWrapRateLimited(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	kid := uuid.Must(uuid.NewV4()).String()

	denied := rate.NewLimiter(rate.Limit(0), 0)
	seq := ddisc.EnqueueUnknown(ddisc.UnknownEnqueueOptionLimiter(denied)).Wrap(q, kid, ddisc.FindMedia(q, kid))

	for range seq.Each(ctx) {
	}
	require.NoError(t, seq.Err())
	require.EqualValues(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_search_queue WHERE known_media_id = ?", kid))
}
