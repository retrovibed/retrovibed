package cmdmedia

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/meta"
	"github.com/stretchr/testify/require"
)

func TestDeeppoolImportKnownFromPublished(t *testing.T) {
	m := deeppoolimport{Source: "deeppool"}

	t.Run("maps title and overview", func(t *testing.T) {
		pc := &meta.PublishedContent{
			Id:          uuid.Must(uuid.NewV7()).String(),
			Title:       "My Movie",
			Description: "A great film",
			Mimetype:    mimex.Video,
		}
		known := m.knownFromPublished(pc)
		require.Equal(t, "My Movie", known.Title)
		require.Equal(t, "A great film", known.Overview)
		require.Equal(t, mimex.Video, known.Mimetype)
		require.Equal(t, "deeppool", known.Source)
		require.Equal(t, pc.Id, known.ID)
	})

	t.Run("prefers known_media_id as uid when present", func(t *testing.T) {
		knownMediaID := uuid.Must(uuid.NewV7()).String()
		pc := &meta.PublishedContent{
			Id:           uuid.Must(uuid.NewV7()).String(),
			KnownMediaId: knownMediaID,
		}
		known := m.knownFromPublished(pc)
		require.Equal(t, knownMediaID, known.UID)
	})

	t.Run("falls back to namespaced pc id when known_media_id absent", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7())
		pc := &meta.PublishedContent{Id: id.String()}
		known := m.knownFromPublished(pc)
		require.NotEmpty(t, known.UID)
		require.NotEqual(t, id.String(), known.UID, "uid should be namespaced, not the raw id")
	})

	t.Run("defaults mimetype to video when blank", func(t *testing.T) {
		pc := &meta.PublishedContent{Id: uuid.Must(uuid.NewV7()).String()}
		known := m.knownFromPublished(pc)
		require.Equal(t, mimex.Video, known.Mimetype)
	})

	t.Run("md5 is stable for same input", func(t *testing.T) {
		pc := &meta.PublishedContent{Id: "01234567-89ab-7def-0123-456789abcdef", Title: "Stable Film"}
		a := m.knownFromPublished(pc)
		b := m.knownFromPublished(pc)
		require.Equal(t, a.Md5, b.Md5)
		require.Equal(t, a.Md5Lower, b.Md5Lower)
	})
}

func TestDeeppoolImportRun(t *testing.T) {
	makeItem := func(id, title string) *meta.PublishedContent {
		return &meta.PublishedContent{Id: id, Title: title, Mimetype: mimex.Video}
	}

	syncResponse := func(items ...*meta.PublishedContent) deeppool.PublishedContentSyncResponse {
		return deeppool.PublishedContentSyncResponse{Items: items}
	}

	newHTTPClient := func(t *testing.T, srv *httptest.Server) *http.Client {
		t.Helper()
		c := &http.Client{}
		c.Transport = httpx.RewriteHostTransport(testx.Must(url.ParseRequestURI(srv.URL))(t), c.Transport)
		return c
	}

	decodeAll := func(t *testing.T, buf *bytes.Buffer) []library.Known {
		t.Helper()
		var results []library.Known
		dec := json.NewDecoder(buf)
		for dec.More() {
			var v library.Known
			require.NoError(t, dec.Decode(&v))
			results = append(results, v)
		}
		return results
	}

	t.Run("returns records from a single page", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		id1 := uuid.Must(uuid.NewV7()).String()
		id2 := uuid.Must(uuid.NewV7()).String()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/p/sync", r.URL.Path)
			if r.URL.Query().Get("sync") == uuid.Nil.String() {
				require.NoError(t, json.NewEncoder(w).Encode(syncResponse(
					makeItem(id1, "Film One"),
					makeItem(id2, "Film Two"),
				)))
			} else {
				require.NoError(t, json.NewEncoder(w).Encode(syncResponse()))
			}
		}))
		defer srv.Close()

		var buf bytes.Buffer
		m := deeppoolimport{Cursor: uuid.Nil.String(), Source: "deeppool"}
		require.NoError(t, m.run(ctx, json.NewEncoder(&buf), newHTTPClient(t, srv)))

		results := decodeAll(t, &buf)
		require.Len(t, results, 2)
		require.Equal(t, "Film One", results[0].Title)
		require.Equal(t, "Film Two", results[1].Title)
	})

	t.Run("paginates until empty page", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		callCount := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			if r.URL.Query().Get("sync") == uuid.Nil.String() {
				require.NoError(t, json.NewEncoder(w).Encode(syncResponse(
					makeItem("aaaaaaaa-0000-0000-0000-000000000001", "Film A"),
					makeItem("aaaaaaaa-0000-0000-0000-000000000002", "Film B"),
				)))
			} else {
				require.NoError(t, json.NewEncoder(w).Encode(syncResponse()))
			}
		}))
		defer srv.Close()

		var buf bytes.Buffer
		m := deeppoolimport{Cursor: uuid.Nil.String(), Source: "deeppool"}
		require.NoError(t, m.run(ctx, json.NewEncoder(&buf), newHTTPClient(t, srv)))
		require.Len(t, decodeAll(t, &buf), 2)
		require.Equal(t, 2, callCount)
	})

	t.Run("advances cursor to last item id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		page1LastID := "bbbbbbbb-0000-0000-0000-000000000002"
		var seenCursors []string

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cursor := r.URL.Query().Get("sync")
			seenCursors = append(seenCursors, cursor)
			if cursor == uuid.Nil.String() {
				require.NoError(t, json.NewEncoder(w).Encode(syncResponse(
					makeItem("bbbbbbbb-0000-0000-0000-000000000001", "Film B1"),
					makeItem(page1LastID, "Film B2"),
				)))
			} else {
				require.NoError(t, json.NewEncoder(w).Encode(syncResponse()))
			}
		}))
		defer srv.Close()

		m := deeppoolimport{Cursor: uuid.Nil.String(), Source: "deeppool"}
		require.NoError(t, m.run(ctx, json.NewEncoder(&bytes.Buffer{}), newHTTPClient(t, srv)))
		require.Equal(t, []string{uuid.Nil.String(), page1LastID}, seenCursors)
	})

	t.Run("returns error on server failure", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		m := deeppoolimport{Cursor: uuid.Nil.String(), Source: "deeppool"}
		require.Error(t, m.run(ctx, json.NewEncoder(&bytes.Buffer{}), newHTTPClient(t, srv)))
	})
}
