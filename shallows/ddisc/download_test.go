package ddisc_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/autobind"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/james-lawrence/torrent/torrenttestx"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/bytesx"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

// servesTorrent starts an httptest server serving a single, freshly built
// .torrent (private per the argument) at "/x.torrent" and returns its URL.
func servesTorrent(t *testing.T, private bool) string {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "content"), []byte("payload"), 0600))

	info := testx.Must(metainfo.NewFromPath(filepath.Join(dir, "content")))(t)
	info.Private = new(private)
	mi := metainfo.MetaInfo{InfoBytes: testx.Must(metainfo.Encode(info))(t)}
	encoded := testx.Must(metainfo.Encode(mi))(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/x.torrent", httptestx.HandleIO(bytes.NewReader(encoded)))
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv.URL + "/x.torrent"
}

func TestDownloadDiscovered(t *testing.T) {
	t.Run("clears the default-private flag once a public torrent is actually resolved", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		importer := tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir()))

		uri := servesTorrent(t, false)
		d := ddisc.NewDiscoveredFromImport(&ddiscapi.Import{Uri: uri})
		require.True(t, d.Private, "an unresolved candidate must default to private")

		d, _, err := ddisc.DownloadDiscovered(ctx, q, importer, d, ddisc.AcquisitionStateDownloading)
		require.NoError(t, err)
		require.False(t, d.Private, "resolving a public torrent must clear the default")

		// the struct returned by DownloadDiscovered is only as trustworthy
		// as what actually landed in ddisc_media - peers are only ever
		// filtered against that row, never against this in-memory value.
		require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE id = ? AND private = false", d.ID))
	})

	t.Run("leaves an actually-private torrent private once resolved", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		importer := tracking.NewURIImport(q, http.DefaultClient, fsx.DirVirtual(t.TempDir()))

		uri := servesTorrent(t, true)
		d := ddisc.NewDiscoveredFromImport(&ddiscapi.Import{Uri: uri})
		require.True(t, d.Private)

		d, _, err := ddisc.DownloadDiscovered(ctx, q, importer, d, ddisc.AcquisitionStateDownloading)
		require.NoError(t, err)
		require.True(t, d.Private, "a BEP 27 private torrent must stay private after resolving")

		require.EqualValues(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM ddisc_media WHERE id = ? AND private = true", d.ID))
	})

	t.Run("the resolved torrent actually downloads real bytes from a peer", func(t *testing.T) {
		// this exercises the whole pipeline for real: the .torrent bytes
		// URIImport fetches over HTTP are the same bytes a torrent.Client
		// needs to open and download the torrent, written to the same
		// rootstore layout tracking.Resume reads from (see
		// tracking.NewResumableMetadata) - not just parsed and discarded.
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		seedir := t.TempDir()
		mi, _, err := torrenttest.Random(seedir, 128*bytesx.KiB)
		require.NoError(t, err)

		seeder := torrenttestx.Client(
			t,
			autobind.NewLoopback(autobind.EnableDHT(torrenttestx.QuickDHT(t))),
			torrent.NewMetadataCache(seedir),
			blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(seedir)),
		)
		defer seeder.Close()

		smd, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(seedir)))
		require.NoError(t, err)
		seederTorrent, _, err := seeder.Start(smd)
		require.NoError(t, err)
		require.NoError(t, torrent.Verify(ctx, seederTorrent))

		// serve the exact same metainfo a search plugin's tracker uri would.
		encoded := testx.Must(metainfo.Encode(metainfo.MetaInfo{InfoBytes: testx.Must(metainfo.Encode(mi))(t)}))(t)
		mux := http.NewServeMux()
		mux.HandleFunc("/x.torrent", httptestx.HandleIO(bytes.NewReader(encoded)))
		torrentSrv := httptest.NewServer(mux)
		defer torrentSrv.Close()

		root := fsx.DirVirtual(t.TempDir())
		leechdir := root.Path("torrent")
		require.NoError(t, fsx.MkDirs(0700, leechdir, root.Path("media")))

		importer := tracking.NewURIImport(q, http.DefaultClient, root)
		disc := ddisc.NewDiscoveredFromImport(&ddiscapi.Import{Uri: torrentSrv.URL + "/x.torrent"})

		disc, lmd, err := ddisc.DownloadDiscovered(ctx, q, importer, disc, ddisc.AcquisitionStateDownloading)
		require.NoError(t, err)
		require.False(t, disc.Private)

		tstore := blockcache.NewTorrentFromVirtualFS(fsx.DirVirtual(leechdir))
		leecher := torrenttestx.Client(
			t,
			autobind.NewLoopback(autobind.EnableDHT(torrenttestx.QuickDHT(t))),
			torrent.NewMetadataCache(leechdir),
			tstore,
		)
		defer leecher.Close()

		_, added, err := tracking.Resume(ctx, q, root, library.QueryCleanerNoop(), leecher, tstore, lmd, asyncx.NewWakeup(ctx), torrent.TuneClientPeer(seeder), torrent.TuneNewConns)
		require.NoError(t, err)
		require.True(t, added)

		var completed tracking.Metadata
		require.Eventually(t, func() bool {
			require.NoError(t, tracking.MetadataFindByID(ctx, q, lmd.ID).Scan(&completed))
			return completed.CompletedAt.Before(time.Now())
		}, 15*time.Second, 100*time.Millisecond)

		require.EqualValues(t, mi.TotalLength(), completed.Downloaded, "the real bytes must actually have been pulled from the seeder, not just the metadata")
	})
}
