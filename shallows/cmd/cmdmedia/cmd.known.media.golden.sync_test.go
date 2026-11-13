package cmdmedia

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
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

		md5 := uuidx.NewV4()
		s1 := uuidx.NewV7().String()
		s2 := uuidx.NewV7().String()
		pages := [][]*deeppool.MediaMetadata{
			{{Uid: uuidx.NewV4().String(), SyncId: s1, Md5: md5.String(), Source: "tmdb", Title: "one", Released: "2020-01-01T00:00:00Z", TombstonedAt: "9999-12-31T23:59:59.999Z", ParentUid: uuidx.Nil().String()}},
			{{Uid: uuidx.NewV4().String(), SyncId: s2, Md5: uuidx.NewV4().String(), Source: "musicbrainz", Title: "two", Collation: 3, Released: "2021-01-01T00:00:00Z", TombstonedAt: "9999-12-31T23:59:59.999Z", ParentUid: uuidx.Nil().String()}},
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

		cmd := goldensync{Sync: uuidx.Nil().String()}
		require.NoError(t, cmd.run(ctx, json.NewEncoder(&buf), &http.Client{Transport: httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil)}))

		require.Equal(t, []string{uuidx.Nil().String(), s1, s2}, cursors)

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
		require.Equal(t, s1, results[0].SyncID)
		require.Equal(t, "two", results[1].Title)
		require.Equal(t, s2, results[1].SyncID)
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

		cmd := goldensync{Sync: uuidx.Nil().String()}
		require.Error(t, cmd.run(ctx, json.NewEncoder(&buf), &http.Client{Transport: httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), nil)}))
	})

	t.Run("cursor defaults to the nil uuid for an empty cache", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		sync, err := goldensync{}.cursor(ctx, db)
		require.NoError(t, err)
		require.Equal(t, uuidx.Nil().String(), sync)
	})

	t.Run("cursor resumes from the latest imported sync id", func(t *testing.T) {
		var a, b, c library.Known
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		s1 := uuidx.NewV7().String()
		s2 := uuidx.NewV7().String()
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults, library.KnownOptionSyncID(s1)))
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults, library.KnownOptionSyncID(s2)))
		require.NoError(t, testx.Fake(&c, library.KnownOptionTestDefaults))

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a, b, c))
		require.NoError(t, knownimport{}.run(ctx, db, &input))

		sync, err := goldensync{}.cursor(ctx, db)
		require.NoError(t, err)
		require.Equal(t, s2, sync)
	})

	t.Run("reimporting from a non golden source does not reset the sync id", func(t *testing.T) {
		var a library.Known
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		s1 := uuidx.NewV7().String()
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults, library.KnownOptionSyncID(s1)))

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a))
		require.NoError(t, knownimport{}.run(ctx, db, &input))

		a.SyncID = ""
		input.Reset()
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a))
		require.NoError(t, knownimport{}.run(ctx, db, &input))

		sync, err := goldensync{}.cursor(ctx, db)
		require.NoError(t, err)
		require.Equal(t, s1, sync)
	})
}
