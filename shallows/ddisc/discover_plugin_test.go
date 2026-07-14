package ddisc_test

import (
	"context"
	"fmt"
	"iter"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func (t fakePluginSeq) Search(ctx context.Context, category, query string) iterx.Seq[*ddiscapi.Import] {
	return fakePluginSeq{results: t.results}
}

type fakePluginSeq struct {
	results []*ddiscapi.Import
}

func (t fakePluginSeq) Each(ctx context.Context) iter.Seq[*ddiscapi.Import] {
	return iterx.From(t.results...)
}

func (t fakePluginSeq) Err() error { return nil }

func TestPluginStrategyYieldsUnpersisted(t *testing.T) {
	q := sqltestx.Metadatabase(t)

	id := int160.Random()
	magnet := fmt.Sprintf("magnet:?xt=urn:btih:%s", id.String())
	kid := uuid.Must(uuid.NewV4()).String()

	plugins := fakePluginSeq{results: []*ddiscapi.Import{{Magnet: magnet, Health: 5, Mimetype: mimex.Video}}}
	seq := ddisc.PluginStrategy(plugins).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid, Title: "ubuntu", Category: mimex.Video})

	var got []ddisc.Discovered
	for v := range seq.Each(t.Context()) {
		got = append(got, v)
	}
	require.NoError(t, seq.Err())
	require.Len(t, got, 1)
	require.Equal(t, kid, got[0].KnownMediaID)
	require.Equal(t, 0, sqltestx.Count(t, q, fmt.Sprintf("SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = '%s'", kid)), "PluginStrategy itself must not persist - that's Discover's job")
}

func TestPluginStrategyNoopsWithoutTitle(t *testing.T) {
	plugins := fakePluginSeq{results: []*ddiscapi.Import{{Magnet: "magnet:?xt=urn:btih:1111111111111111111111111111111111111111"}}}
	seq := ddisc.PluginStrategy(plugins).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: uuid.Must(uuid.NewV4()).String()})

	var count int
	for range seq.Each(t.Context()) {
		count++
	}
	require.NoError(t, seq.Err())
	require.Equal(t, 0, count)
}
