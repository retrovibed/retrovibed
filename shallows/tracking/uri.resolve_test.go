package tracking_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httptestx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestURIImportResolveHTTP(t *testing.T) {
	t.Run("should resolve a torrent served over http and cache it to disk", func(t *testing.T) {
		info := testx.Must(metainfo.NewFromPath(testx.Fixture()))(t)
		md := metainfo.MetaInfo{
			InfoBytes: testx.Must(metainfo.Encode(info))(t),
		}
		buf := testx.Must(metainfo.Encode(md))(t)

		mux := http.NewServeMux()
		mux.HandleFunc("/example.torrent", httptestx.HandleIO(bytes.NewReader(buf)))
		srv := httptest.NewServer(mux)
		defer srv.Close()

		vfs := fsx.DirVirtual(t.TempDir())
		uri := tracking.NewURIImport(nil, http.DefaultClient, vfs)

		meta, err := uri.Resolve(t.Context(), srv.URL+"/example.torrent")
		require.NoError(t, err)
		require.Equal(t, md.ID().Bytes(), meta.Infohash)
		require.Equal(t, mimex.Bittorrent, meta.Mimetype)
		require.EqualValues(t, info.TotalLength(), meta.Bytes)

		cached, err := os.ReadFile(filepath.Join(vfs.Path("torrent"), md.HashInfoBytes().String()+".torrent"))
		require.NoError(t, err)
		require.Equal(t, buf, cached)
	})

	t.Run("should follow a single redirect before resolving the torrent", func(t *testing.T) {
		info := testx.Must(metainfo.NewFromPath(testx.Fixture()))(t)
		md := metainfo.MetaInfo{
			InfoBytes: testx.Must(metainfo.Encode(info))(t),
		}
		buf := testx.Must(metainfo.Encode(md))(t)

		mux := http.NewServeMux()
		mux.HandleFunc("/redirect.torrent", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/example.torrent", http.StatusFound)
		})
		mux.HandleFunc("/example.torrent", httptestx.HandleIO(bytes.NewReader(buf)))
		srv := httptest.NewServer(mux)
		defer srv.Close()

		vfs := fsx.DirVirtual(t.TempDir())
		uri := tracking.NewURIImport(nil, http.DefaultClient, vfs)

		meta, err := uri.Resolve(t.Context(), srv.URL+"/redirect.torrent")
		require.NoError(t, err)
		require.Equal(t, md.ID().Bytes(), meta.Infohash)
		require.Equal(t, mimex.Bittorrent, meta.Mimetype)
		require.EqualValues(t, info.TotalLength(), meta.Bytes)

		cached, err := os.ReadFile(filepath.Join(vfs.Path("torrent"), md.HashInfoBytes().String()+".torrent"))
		require.NoError(t, err)
		require.Equal(t, buf, cached)
	})

	t.Run("should error when the endpoint responds with a failure status", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/missing.torrent", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		vfs := fsx.DirVirtual(t.TempDir())
		uri := tracking.NewURIImport(nil, http.DefaultClient, vfs)

		_, err := uri.Resolve(t.Context(), srv.URL+"/missing.torrent")
		require.Error(t, err)
	})

	t.Run("should error when the endpoint responds with a body that isn't a torrent", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/notatorrent", httptestx.HandleIO(bytes.NewReader([]byte("not a torrent"))))
		srv := httptest.NewServer(mux)
		defer srv.Close()

		vfs := fsx.DirVirtual(t.TempDir())
		uri := tracking.NewURIImport(nil, http.DefaultClient, vfs)

		_, err := uri.Resolve(t.Context(), srv.URL+"/notatorrent")
		require.Error(t, err)
	})
}
