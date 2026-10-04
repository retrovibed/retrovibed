package cmdmedia

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownGoldenSync(t *testing.T) {
	t.Run("pages until an empty response and writes known jsonl", func(t *testing.T) {
		var (
			buf     bytes.Buffer
			cursors []string
			results []library.Known
		)
		ctx, done := testx.Context(t)
		defer done()

		md5 := uuid.Must(uuid.NewV4())
		s1 := uuid.Must(uuid.NewV7()).String()
		s2 := uuid.Must(uuid.NewV7()).String()
		pages := [][]*deeppool.MediaMetadata{
			{{Uid: uuid.Must(uuid.NewV4()).String(), SyncId: s1, Md5: md5.String(), Source: "tmdb", Title: "one", Released: "2020-01-01T00:00:00Z", TombstonedAt: "9999-12-31T23:59:59.999Z", ParentUid: uuid.Nil.String()}},
			{{Uid: uuid.Must(uuid.NewV4()).String(), SyncId: s2, Md5: uuid.Must(uuid.NewV4()).String(), Source: "musicbrainz", Title: "two", Collation: 3, Released: "2021-01-01T00:00:00Z", TombstonedAt: "9999-12-31T23:59:59.999Z", ParentUid: uuid.Nil.String()}},
			{},
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/mm/sync", r.URL.Path)
			cursors = append(cursors, r.URL.Query().Get("sync"))
			items := pages[len(cursors)-1]
			resp := deeppool.MediaMetadataSyncResponse{Next: &deeppool.MediaMetadataSyncRequest{}, Items: items}
			if len(items) > 0 {
				resp.Next.Sync = items[len(items)-1].SyncId
			}
			require.NoError(t, json.NewEncoder(w).Encode(&resp))
		}))
		defer srv.Close()

		cmd := goldensync{Sync: uuid.Nil.String()}
		require.NoError(t, cmd.run(ctx, json.NewEncoder(&buf), &http.Client{Transport: httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil)}))

		require.Equal(t, []string{uuid.Nil.String(), s1, s2}, cursors)

		dec := json.NewDecoder(&buf)
		for dec.More() {
			var v library.Known
			require.NoError(t, dec.Decode(&v))
			results = append(results, v)
		}

		require.Len(t, results, 2)
		require.Equal(t, "one", results[0].Title)
		require.Equal(t, "tmdb", results[0].Source)
		require.Equal(t, md5.String(), results[0].Md5)
		require.Equal(t, binary.LittleEndian.Uint64(uuidx.LowN(md5, 64)), results[0].Md5Lower)
		require.Equal(t, 2020, results[0].Released.Year())
		require.Equal(t, "two", results[1].Title)
		require.EqualValues(t, 3, results[1].Collation)
	})

	t.Run("propagates server errors", func(t *testing.T) {
		var buf bytes.Buffer
		ctx, done := testx.Context(t)
		defer done()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		cmd := goldensync{Sync: uuid.Nil.String()}
		require.Error(t, cmd.run(ctx, json.NewEncoder(&buf), &http.Client{Transport: httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil)}))
	})
}
