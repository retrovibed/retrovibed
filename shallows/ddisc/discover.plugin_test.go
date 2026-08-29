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
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func (t fakePluginSeq) Search(ctx context.Context, mimetypes []string, query string, adult, public bool) iterx.Seq[*ddiscapi.Import] {
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

	plugins := fakePluginSeq{results: []*ddiscapi.Import{{Uri: magnet, Uritype: mimex.Magnet, Health: 5, Bytes: 1234, Mimetype: mimex.Video}}}
	seq := ddisc.PluginStrategy(plugins).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid, Query: "ubuntu", Mimetypes: []string{mimex.Video}})

	var got []ddisc.Discovered
	for v := range seq.Each(t.Context()) {
		got = append(got, v)
	}
	require.NoError(t, seq.Err())
	require.Len(t, got, 1)
	require.Equal(t, kid, got[0].KnownMediaID)
	require.EqualValues(t, 1234, got[0].Bytes)
	require.Equal(t, 0, sqltestx.Count(t, q, fmt.Sprintf("SELECT COUNT(*) FROM ddisc_media WHERE known_media_id = '%s'", kid)), "PluginStrategy itself must not persist - that's the caller's job")
}

func TestPluginStrategyNoopsWithoutTitle(t *testing.T) {
	q := sqltestx.Metadatabase(t)

	plugins := fakePluginSeq{results: []*ddiscapi.Import{{Uri: "magnet:?xt=urn:btih:1111111111111111111111111111111111111111", Uritype: mimex.Magnet}}}
	seq := ddisc.PluginStrategy(plugins).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: uuid.Must(uuid.NewV4()).String()})

	var count int
	for range seq.Each(t.Context()) {
		count++
	}
	require.NoError(t, seq.Err())
	require.Equal(t, 0, count)
	require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM library_known_media"))
}

func TestPluginStrategyResolvesNonMagnetURI(t *testing.T) {
	uri := "https://tracker.example/download/1000.abc"
	kid := uuid.Must(uuid.NewV4()).String()

	plugins := fakePluginSeq{results: []*ddiscapi.Import{{
		Uri:      uri,
		Uritype:  mimex.Bittorrent,
		Health:   7,
		Mimetype: mimex.Video,
		Title:    "Some Release",
	}}}
	seq := ddisc.PluginStrategy(plugins).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid, Query: "ubuntu", Mimetypes: []string{mimex.Video}})

	var got []ddisc.Discovered
	for v := range seq.Each(t.Context()) {
		got = append(got, v)
	}
	require.NoError(t, seq.Err())
	require.Len(t, got, 1, "a non-magnet uri (e.g. unit3d's download_link fallback) must still be usable, not dropped as unresolvable")
	require.Equal(t, uri, got[0].URI)
	require.Equal(t, md5x.FormatUUID(md5x.Digest(uri)), got[0].ID)
	require.Equal(t, "Some Release", got[0].Title)
}
