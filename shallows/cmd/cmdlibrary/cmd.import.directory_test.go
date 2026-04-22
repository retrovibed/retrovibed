package cmdlibrary

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportDirectory(t *testing.T) {
	newClient := func(t *testing.T, srv *httptest.Server) *http.Client {
		t.Helper()
		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		return c
	}

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

		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "a.mp4"), []byte("video content"), 0600))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "b.mp4"), []byte("more video"), 0600))

		uploadCount := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/m/", r.URL.Path)
			uploadCount++
			id := uuid.Must(uuid.NewV4()).String()
			assert.NoError(t, json.NewEncoder(w).Encode(&media.MediaUploadResponse{
				Media: &media.Media{Id: id},
			}))
		}))
		defer srv.Close()

		var buf bytes.Buffer
		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: dir}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&buf), newClient(t, srv)))

		require.Equal(t, 2, uploadCount)
		results := decodeAll(t, &buf)
		require.Len(t, results, 2)
	})

	t.Run("skips subdirectories", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "file.mp4"), []byte("content"), 0600))
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "subdir"), 0700))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "subdir", "nested.mp4"), []byte("nested"), 0600))

		uploadCount := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uploadCount++
			assert.NoError(t, json.NewEncoder(w).Encode(&media.MediaUploadResponse{
				Media: &media.Media{Id: uuid.Must(uuid.NewV4()).String()},
			}))
		}))
		defer srv.Close()

		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: dir}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&bytes.Buffer{}), newClient(t, srv)))

		require.Equal(t, 1, uploadCount, "only the immediate file should be uploaded, not nested ones")
	})

	t.Run("empty directory produces no uploads", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()

		uploadCount := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uploadCount++
		}))
		defer srv.Close()

		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: dir}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&bytes.Buffer{}), newClient(t, srv)))

		require.Equal(t, 0, uploadCount)
	})

	t.Run("server error returns error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "file.mp4"), []byte("content"), 0600))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: dir}
		require.Error(t, cmd.run(ctx, jsonl.NewEncoder(&bytes.Buffer{}), newClient(t, srv)))
	})

	t.Run("multipart field name is content and filename matches basename", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "my-video.mp4"), []byte("data"), 0600))

		var receivedFilename, receivedField string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, r.ParseMultipartForm(1<<20))
			for _, fh := range r.MultipartForm.File["content"] {
				receivedField = "content"
				receivedFilename = fh.Filename
			}
			assert.NoError(t, json.NewEncoder(w).Encode(&media.MediaUploadResponse{
				Media: &media.Media{Id: uuid.Must(uuid.NewV4()).String()},
			}))
		}))
		defer srv.Close()

		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: dir}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&bytes.Buffer{}), newClient(t, srv)))

		require.Equal(t, "content", receivedField)
		require.Equal(t, "my-video.mp4", receivedFilename)
	})

	t.Run("mime type detected from extension", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "video.mp4"), []byte("data"), 0600))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "unknown.xyz"), []byte("data"), 0600))

		received := map[string]string{}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, r.ParseMultipartForm(1<<20))
			for _, fh := range r.MultipartForm.File["content"] {
				received[fh.Filename] = fh.Header.Get("Content-Type")
			}
			assert.NoError(t, json.NewEncoder(w).Encode(&media.MediaUploadResponse{
				Media: &media.Media{Id: uuid.Must(uuid.NewV4()).String()},
			}))
		}))
		defer srv.Close()

		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: dir}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&bytes.Buffer{}), newClient(t, srv)))

		require.Equal(t, "video/mp4", received["video.mp4"])
		require.NotEmpty(t, received["unknown.xyz"], "unknown extension should fall back to binary mimetype")
	})

	t.Run("output contains uploaded media IDs", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "clip.mp4"), []byte("content"), 0600))

		expectedID := uuid.Must(uuid.NewV4()).String()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.NoError(t, json.NewEncoder(w).Encode(&media.MediaUploadResponse{
				Media: &media.Media{Id: expectedID, Description: "clip.mp4"},
			}))
		}))
		defer srv.Close()

		var buf bytes.Buffer
		cmd := importDirectory{Endpoint: "localhost:9998", Concurrency: 1, Directory: dir}
		require.NoError(t, cmd.run(ctx, jsonl.NewEncoder(&buf), newClient(t, srv)))

		results := decodeAll(t, &buf)
		require.Len(t, results, 1)
		require.Equal(t, expectedID, results[0].Id)
		require.Equal(t, "clip.mp4", results[0].Description)
	})
}
