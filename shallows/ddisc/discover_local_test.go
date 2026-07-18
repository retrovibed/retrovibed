package ddisc_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestLocalStrategyExactMatch(t *testing.T) {
	q := sqltestx.Metadatabase(t)

	kid := uuid.Must(uuid.NewV4()).String()
	id := int160.Random()
	d := ddisc.NewDiscovered(&id, "", ddisc.DiscoveredOptionKnownMedia(kid), ddisc.DiscoveredOptionMimetype(mimex.Binary))
	require.NoError(t, ddisc.DiscoveredInsertWithDefaults(t.Context(), q, d).Scan(&d))

	seq := ddisc.LocalStrategy(q).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid})
	var got []ddisc.Discovered
	for v := range seq.Each(t.Context()) {
		got = append(got, v)
	}
	require.NoError(t, seq.Err())
	require.Len(t, got, 1)
	require.Equal(t, kid, got[0].KnownMediaID)
}

func TestLocalStrategyMissReturnsEmpty(t *testing.T) {
	q := sqltestx.Metadatabase(t)

	seq := ddisc.LocalStrategy(q).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: uuid.Must(uuid.NewV4()).String()})
	var count int
	for range seq.Each(t.Context()) {
		count++
	}
	require.NoError(t, seq.Err())
	require.Equal(t, 0, count)
}
