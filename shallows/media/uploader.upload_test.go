package media_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/stretchr/testify/require"
)

func TestUploaderUpload(t *testing.T) {
	t.Run("uploads into the library root", func(t *testing.T) {
		var (
			md library.Metadata
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		mediastore := fsx.DirVirtual(t.TempDir())

		routes := mux.NewRouter()
		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			asyncx.NewWakeup(t.Context()),
			mediastore,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/m").Subrouter())

		srv := httptest.NewServer(routes)
		t.Cleanup(srv.Close)

		c := &http.Client{
			Transport: httpx.NewHeadersTransport(http.Header{"Authorization": []string{httpauthtest.UnsafeTokenAuto(t)}}, httpx.HTORoundTripper(
				httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil),
			)),
		}

		path := filepath.Join(t.TempDir(), "example.mkv")
		require.NoError(t, os.WriteFile(path, []byte("example"), 0600))

		id, err := media.NewUploader(c, "https://localhost:9998").Upload(ctx, path, uuid.Nil.String())
		require.NoError(t, err)

		require.NoError(t, library.MetadataFindByID(ctx, q, id).Scan(&md))
		require.Equal(t, "example.mkv", md.Description)
		require.Equal(t, uint64(len("example")), md.Bytes)
		require.Equal(t, uuid.Nil.String(), md.DirectoryID)
		require.DirExists(t, mediastore.Path(id)) // uploads are stored as a block cache directory.
	})

	t.Run("files the upload into the given library directory", func(t *testing.T) {
		var (
			md library.Metadata
		)

		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		mediastore := fsx.DirVirtual(t.TempDir())
		libdir := uuid.Must(uuid.NewV4()).String()

		routes := mux.NewRouter()
		media.NewHTTPLibrary(
			q,
			asyncx.NewWakeup(t.Context()),
			asyncx.NewWakeup(t.Context()),
			mediastore,
			nil,
			media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
		).Bind(routes.PathPrefix("/m").Subrouter())

		srv := httptest.NewServer(routes)
		t.Cleanup(srv.Close)

		c := &http.Client{
			Transport: httpx.NewHeadersTransport(http.Header{"Authorization": []string{httpauthtest.UnsafeTokenAuto(t)}}, httpx.HTORoundTripper(
				httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil),
			)),
		}

		path := filepath.Join(t.TempDir(), "example.mkv")
		require.NoError(t, os.WriteFile(path, []byte("example"), 0600))

		id, err := media.NewUploader(c, "https://localhost:9998").Upload(ctx, path, libdir)
		require.NoError(t, err)

		require.NoError(t, library.MetadataFindByID(ctx, q, id).Scan(&md))
		require.Equal(t, libdir, md.DirectoryID)
	})

	t.Run("missing file errors", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		_, err := media.NewUploader(http.DefaultClient, "https://localhost:9998").Upload(ctx, filepath.Join(t.TempDir(), "missing.mkv"), uuid.Nil.String())
		require.Error(t, err)
	})
}
