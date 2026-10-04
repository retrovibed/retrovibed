package cmdmedia

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownGoldenUpload(t *testing.T) {
	t.Run("uploads stdin jsonl in batches", func(t *testing.T) {
		var (
			input   bytes.Buffer
			batches [][]*deeppool.MediaMetadata
		)
		ctx, done := testx.Context(t)
		defer done()

		enc := json.NewEncoder(&input)
		for _, title := range []string{"one", "two", "three"} {
			require.NoError(t, enc.Encode(library.Known{
				UID:       uuid.Must(uuid.NewV4()).String(),
				Md5:       uuid.Must(uuid.NewV4()).String(),
				Source:    "tmdb",
				Title:     title,
				Collation: 7,
				Released:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				ParentUID: uuid.Nil.String(),
			}))
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req deeppool.MediaMetadataUploadRequest
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/mm/", r.URL.Path)
			require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
			batches = append(batches, req.Items)
			require.NoError(t, json.NewEncoder(w).Encode(&deeppool.MediaMetadataUploadResponse{Received: uint64(len(req.Items)), Upserted: uint64(len(req.Items))}))
		}))
		defer srv.Close()

		cmd := goldenupload{Batch: 2}
		require.NoError(t, cmd.run(ctx, &input, &http.Client{Transport: httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil)}))

		require.Len(t, batches, 2)
		require.Len(t, batches[0], 2)
		require.Len(t, batches[1], 1)
		require.Equal(t, "one", batches[0][0].Title)
		require.Equal(t, "three", batches[1][0].Title)
		require.Equal(t, "tmdb", batches[0][0].Source)
		require.EqualValues(t, 7, batches[0][0].Collation)
		require.Equal(t, "2020-01-01T00:00:00Z", batches[0][0].Released)
	})

	t.Run("propagates server errors", func(t *testing.T) {
		var input bytes.Buffer
		ctx, done := testx.Context(t)
		defer done()

		require.NoError(t, json.NewEncoder(&input).Encode(library.Known{UID: uuid.Must(uuid.NewV4()).String(), Md5: uuid.Must(uuid.NewV4()).String()}))

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		cmd := goldenupload{Batch: 2}
		require.Error(t, cmd.run(ctx, &input, &http.Client{Transport: httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil)}))
	})
}
