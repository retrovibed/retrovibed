package cmdmedia

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/stretchr/testify/require"
)

func TestKnownGoldenEnv(t *testing.T) {
	t.Run("outputs RETROVIBED_ARCHIVE_START_DATE from deeppool", func(t *testing.T) {
		var (
			buf     bytes.Buffer
			include []string
			exclude []string
		)
		ctx, done := testx.Context(t)
		defer done()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/mm/env", r.URL.Path)
			include = r.URL.Query()["include_source"]
			exclude = r.URL.Query()["exclude_source"]
			require.NoError(t, json.NewEncoder(w).Encode(&deeppool.MediaMetadataEnvResponse{Start: "2023-03-01"}))
		}))
		defer srv.Close()

		cmd := goldenenv{IncludeSource: []string{"tmdb"}, ExcludeSource: []string{"musicbrainz"}}
		require.NoError(t, cmd.run(ctx, &buf, &http.Client{Transport: httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil)}))

		require.Equal(t, "RETROVIBED_ARCHIVE_START_DATE=2023-03-01\n", buf.String())
		require.Equal(t, []string{"tmdb"}, include)
		require.Equal(t, []string{"musicbrainz"}, exclude)
	})

	t.Run("propagates server errors", func(t *testing.T) {
		var buf bytes.Buffer
		ctx, done := testx.Context(t)
		defer done()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		cmd := goldenenv{IncludeSource: []string{"tmdb"}}
		require.Error(t, cmd.run(ctx, &buf, &http.Client{Transport: httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil)}))
		require.Empty(t, buf.String())
	})
}
