package ddisc_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
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
}
