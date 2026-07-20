package ddisc_test

import (
	"fmt"
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
	d := ddisc.NewDiscovered(&id, ddisc.DiscoveredOptionKnownMedia(kid), ddisc.DiscoveredOptionMimetype(mimex.Binary), ddisc.DiscoveredOptionAutoMagnet)
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

// The empty-string case is what a title-only discovery request actually
// carries (no known-media-id given) - SyncStrategies includes LocalStrategy
// for any KnownMediaID != uuid.Nil.String(), which "" satisfies, so this is
// the case that reaches LocalStrategy in that flow. Without the early
// return in localSeq.Each, DiscoveredQueryKnownMediaID("") resolves to a
// no-op predicate and, as the sole argument to .Where(...), rendered a
// dangling "WHERE" with nothing after it - which DuckDB's parser rejected
// with "Parser Error: syntax error at end of input".
func TestLocalStrategyReturnsEmpty(t *testing.T) {
	for _, kid := range []string{"", uuid.Nil.String(), uuid.Must(uuid.NewV4()).String()} {
		t.Run(fmt.Sprintf("known_media_id=%q", kid), func(t *testing.T) {
			q := sqltestx.Metadatabase(t)

			seq := ddisc.LocalStrategy(q).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid})
			var count int
			for range seq.Each(t.Context()) {
				count++
			}
			require.NoError(t, seq.Err())
			require.Equal(t, 0, count)
		})
	}
}
