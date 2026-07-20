package ddisc_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownStrategyCatalogMatch(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	var known library.Known
	require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Video)))
	require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

	seq := ddisc.KnownStrategy(q).Discover(ctx, ddisc.DiscoverRequest{Query: known.Title, Mimetypes: ddisc.Category(known.Mimetype)})
	var got []ddisc.Discovered
	for v := range seq.Each(ctx) {
		got = append(got, v)
	}
	require.NoError(t, seq.Err())
	require.Len(t, got, 1)
	require.Equal(t, known.UID, got[0].KnownMediaID)
	require.Equal(t, ddisc.PolicyRejectionCatalogOnly, got[0].PolicyRejection)
	require.Len(t, got[0].Infohash, 20)
	require.NotEmpty(t, got[0].URI)
}

func TestKnownStrategyMissReturnsEmpty(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	seq := ddisc.KnownStrategy(q).Discover(ctx, ddisc.DiscoverRequest{Query: "zzzznonexistenttitlezzzz9182734"})
	var count int
	for range seq.Each(ctx) {
		count++
	}
	require.NoError(t, seq.Err())
	require.Equal(t, 0, count)
}

func TestKnownStrategyEmptyTitleReturnsEmpty(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	seq := ddisc.KnownStrategy(q).Discover(ctx, ddisc.DiscoverRequest{})
	var count int
	for range seq.Each(ctx) {
		count++
	}
	require.NoError(t, seq.Err())
	require.Equal(t, 0, count)
}
