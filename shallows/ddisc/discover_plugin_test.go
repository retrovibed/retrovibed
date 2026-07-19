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
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func (t fakePluginSeq) Search(ctx context.Context, mimetypes []string, query string, adult bool) iterx.Seq[*ddiscapi.Import] {
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
	seq := ddisc.PluginStrategy(q, plugins).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid, Title: "ubuntu", Mimetypes: []string{mimex.Video}})

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
	seq := ddisc.PluginStrategy(q, plugins).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: uuid.Must(uuid.NewV4()).String()})

	var count int
	for range seq.Each(t.Context()) {
		count++
	}
	require.NoError(t, seq.Err())
	require.Equal(t, 0, count)
	require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM library_known_media"))
}

func TestPluginStrategyRecordsKnownMediaTOFU(t *testing.T) {
	q := sqltestx.Metadatabase(t)

	id := int160.Random()
	magnet := fmt.Sprintf("magnet:?xt=urn:btih:%s", id.String())
	kid := uuid.Must(uuid.NewV4()).String()

	plugins := fakePluginSeq{results: []*ddiscapi.Import{{
		Uri:          magnet,
		Uritype:      mimex.Magnet,
		Health:       5,
		Mimetype:     mimex.Video,
		KnownMediaId: kid,
		Title:        "Ubuntu Documentary",
		Overview:     "a documentary about ubuntu",
		Popularity:   4.2,
		PosterPath:   "/ubuntu.jpg",
		Source:       "unit3d",
	}}}
	seq := ddisc.PluginStrategy(q, plugins).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid, Title: "ubuntu", Mimetypes: []string{mimex.Video}})

	for range seq.Each(t.Context()) {
	}
	require.NoError(t, seq.Err())

	var known library.Known
	require.NoError(t, library.KnownFindByID(t.Context(), q, kid).Scan(&known))
	require.Equal(t, "Ubuntu Documentary", known.Title)
	require.Equal(t, "a documentary about ubuntu", known.Overview)
	require.Equal(t, 4.2, known.Popularity)
	require.Equal(t, "/ubuntu.jpg", known.PosterPath)
	require.Equal(t, "unit3d", known.Source)
}

func TestPluginStrategyKnownMediaTOFUDoesNotOverwrite(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	kid := uuid.Must(uuid.NewV4()).String()

	first := fakePluginSeq{results: []*ddiscapi.Import{{
		Uri:          fmt.Sprintf("magnet:?xt=urn:btih:%s", int160.Random()),
		Uritype:      mimex.Magnet,
		Mimetype:     mimex.Video,
		KnownMediaId: kid,
		Title:        "original title",
		Source:       "unit3d",
	}}}
	seq := ddisc.PluginStrategy(q, first).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid, Title: "q", Mimetypes: []string{mimex.Video}})
	for range seq.Each(t.Context()) {
	}
	require.NoError(t, seq.Err())

	second := fakePluginSeq{results: []*ddiscapi.Import{{
		Uri:          fmt.Sprintf("magnet:?xt=urn:btih:%s", int160.Random()),
		Uritype:      mimex.Magnet,
		Mimetype:     mimex.Video,
		KnownMediaId: kid,
		Title:        "different title",
		Source:       "leetx",
	}}}
	seq2 := ddisc.PluginStrategy(q, second).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid, Title: "q", Mimetypes: []string{mimex.Video}})
	for range seq2.Each(t.Context()) {
	}
	require.NoError(t, seq2.Err())

	var known library.Known
	require.NoError(t, library.KnownFindByID(t.Context(), q, kid).Scan(&known))
	require.Equal(t, "original title", known.Title)
	require.Equal(t, "unit3d", known.Source)
}

func TestPluginStrategyRecordsKnownMediaWithoutTitle(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	kid := uuid.Must(uuid.NewV4()).String()

	plugins := fakePluginSeq{results: []*ddiscapi.Import{{
		Uri:          fmt.Sprintf("magnet:?xt=urn:btih:%s", int160.Random()),
		Uritype:      mimex.Magnet,
		Mimetype:     mimex.Video,
		KnownMediaId: kid,
		Source:       "unit3d",
	}}}
	seq := ddisc.PluginStrategy(q, plugins).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid, Title: "q", Mimetypes: []string{mimex.Video}})
	for range seq.Each(t.Context()) {
	}
	require.NoError(t, seq.Err())

	var known library.Known
	require.NoError(t, library.KnownFindByID(t.Context(), q, kid).Scan(&known))
	require.Equal(t, kid, known.UID)
	require.Equal(t, "unit3d", known.Source)
}

func TestPluginStrategySkipsSentinelKnownMediaID(t *testing.T) {
	for _, kid := range []string{"", uuid.Nil.String(), uuid.Max.String()} {
		t.Run(fmt.Sprintf("known_media_id=%q", kid), func(t *testing.T) {
			q := sqltestx.Metadatabase(t)

			plugins := fakePluginSeq{results: []*ddiscapi.Import{{
				Uri:          fmt.Sprintf("magnet:?xt=urn:btih:%s", int160.Random()),
				Uritype:      mimex.Magnet,
				Mimetype:     mimex.Video,
				KnownMediaId: kid,
				Title:        "some title",
				Source:       "unit3d",
			}}}
			seq := ddisc.PluginStrategy(q, plugins).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: uuid.Must(uuid.NewV4()).String(), Title: "q", Mimetypes: []string{mimex.Video}})
			for range seq.Each(t.Context()) {
			}
			require.NoError(t, seq.Err())
			require.Equal(t, 0, sqltestx.Count(t, q, "SELECT COUNT(*) FROM library_known_media"))
		})
	}
}

func TestPluginStrategyResolvesNonMagnetURI(t *testing.T) {
	q := sqltestx.Metadatabase(t)

	uri := "https://tracker.example/download/1000.abc"
	kid := uuid.Must(uuid.NewV4()).String()

	plugins := fakePluginSeq{results: []*ddiscapi.Import{{
		Uri:      uri,
		Uritype:  mimex.Bittorrent,
		Health:   7,
		Mimetype: mimex.Video,
		Title:    "Some Release",
	}}}
	seq := ddisc.PluginStrategy(q, plugins).Discover(t.Context(), ddisc.DiscoverRequest{KnownMediaID: kid, Title: "ubuntu", Mimetypes: []string{mimex.Video}})

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
