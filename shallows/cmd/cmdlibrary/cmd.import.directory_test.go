package cmdlibrary

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/shallows/httpauthtest"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/stretchr/testify/require"
)

func newImportDirectoryServer(t *testing.T, q *sql.DB) *http.Client {
	t.Helper()

	vfs := fsx.DirVirtual(t.TempDir())
	routes := mux.NewRouter()
	media.NewHTTPLibrary(
		q,
		asyncx.NewWakeup(t.Context()),
		asyncx.NewWakeup(t.Context()),
		vfs,
		nil,
		media.HTTPLibraryOptionJWTSecret(httpauthtest.UnsafeJWTSecretSource),
	).Bind(routes.PathPrefix("/m").Subrouter())

	srv := httptest.NewServer(routes)
	t.Cleanup(srv.Close)

	headers := http.Header{"Authorization": []string{httpauthtest.UnsafeTokenAuto(t)}}
	return &http.Client{
		Transport: httpx.NewHeadersTransport(headers, httpx.HTORoundTripper(
			httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil),
		)),
	}
}

func TestImportDirectory(t *testing.T) {
	decodeAll := func(t *testing.T, buf *bytes.Buffer) []*media.Media {
		t.Helper()
		var results []*media.Media
		dec := json.NewDecoder(buf)
		for dec.More() {
			var v media.Media
			require.NoError(t, dec.Decode(&v))
			results = append(results, &v)
		}
		return results
	}

	t.Run("uploads each file and writes result to encoder", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "a.mp4"), []byte("video content"), 0600))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "b.mp4"), []byte("more video"), 0600))

		var buf bytes.Buffer
		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: dir}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&buf), newImportDirectoryServer(t, q)))

		require.Equal(t, 2, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata"))(t))
		require.Len(t, decodeAll(t, &buf), 2)
	})

	t.Run("recurses into subdirectories", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "file.mp4"), []byte("content"), 0600))
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "subdir"), 0700))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "subdir", "nested.mp4"), []byte("nested"), 0600))

		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: dir}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&bytes.Buffer{}), newImportDirectoryServer(t, q)))

		require.Equal(t, 2, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata"))(t))
	})

	t.Run("empty directory produces no uploads", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: t.TempDir()}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&bytes.Buffer{}), newImportDirectoryServer(t, q)))

		require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata"))(t))
	})

	t.Run("filename matches basename", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		dir := filepath.Join(t.TempDir(), "media")
		require.NoError(t, os.MkdirAll(dir, 0700))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "my-video.mp4"), []byte("data"), 0600))

		var buf bytes.Buffer
		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: dir}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&buf), newImportDirectoryServer(t, q)))

		results := decodeAll(t, &buf)
		require.Len(t, results, 1)

		var md library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, results[0].Id).Scan(&md))
		require.Equal(t, " media my-video.mp4", md.Description)
	})

	t.Run("mime type detected from extension", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		dir := filepath.Join(t.TempDir(), "media")
		require.NoError(t, os.MkdirAll(dir, 0700))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "video.mp4"), []byte("video data"), 0600))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "unknown.xyz"), []byte("unknown data"), 0600))

		var buf bytes.Buffer
		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: dir}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&buf), newImportDirectoryServer(t, q)))

		results := decodeAll(t, &buf)
		require.Len(t, results, 2)

		mimetypes := map[string]string{}
		for _, m := range results {
			var md library.Metadata
			require.NoError(t, library.MetadataFindByID(ctx, q, m.Id).Scan(&md))
			mimetypes[md.Description] = md.Mimetype
		}

		require.Equal(t, "video/mp4", mimetypes[" media video.mp4"])
		require.NotEmpty(t, mimetypes[" media unknown.xyz"])
	})

	t.Run("mimetype flag overrides extension detection", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "video.mp4"), []byte("video data"), 0600))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "unknown.xyz"), []byte("unknown data"), 0600))

		var buf bytes.Buffer
		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Mimetype: "video/quicktime", Directory: dir}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&buf), newImportDirectoryServer(t, q)))

		results := decodeAll(t, &buf)
		require.Len(t, results, 2)

		for _, m := range results {
			var md library.Metadata
			require.NoError(t, library.MetadataFindByID(ctx, q, m.Id).Scan(&md))
			require.Equal(t, "video/quicktime", md.Mimetype)
		}
	})

	t.Run("output contains uploaded media IDs", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		dir := filepath.Join(t.TempDir(), "media")
		require.NoError(t, os.MkdirAll(dir, 0700))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "clip.mp4"), []byte("content"), 0600))

		var buf bytes.Buffer
		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: dir}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&buf), newImportDirectoryServer(t, q)))

		results := decodeAll(t, &buf)
		require.Len(t, results, 1)

		var md library.Metadata
		require.NoError(t, library.MetadataFindByID(ctx, q, results[0].Id).Scan(&md))
		require.Equal(t, " media clip.mp4", md.Description)
	})
}
