package cmdmedia

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/dashotv/tvdb"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

// tvdb.Client has no exported way to point it at a custom server URL or HTTP
// client, so requests are redirected by rewriting the process-wide default
// transport for the duration of the test.
func newTvdbTestClient(t *testing.T, srv *httptest.Server) *tvdb.Client {
	t.Helper()
	prev := http.DefaultTransport
	http.DefaultTransport = httpx.RewriteHostTransport(testx.Must(url.Parse(srv.URL))(t), prev)
	t.Cleanup(func() { http.DefaultTransport = prev })
	return tvdb.New("test-key", "test-token")
}

func TestTvdbImportImgpath(t *testing.T) {
	t.Run("blank path stays blank", func(t *testing.T) {
		tm := tvdbimport{URL: "https://thetvdb.com"}
		require.Equal(t, "", tm.imgpath(""))
	})

	t.Run("non-blank path is prefixed with the configured base url", func(t *testing.T) {
		tm := tvdbimport{URL: "https://thetvdb.com"}
		require.Equal(t, "https://thetvdb.com/banners/posters/121361-1.jpg", tm.imgpath("/banners/posters/121361-1.jpg"))
	})
}

func TestTvdbImportRecords(t *testing.T) {
	t.Run("maps series fields onto the known record", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v4/series", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			errorsx.Zero(fmt.Fprint(w, `{"data":[{
				"id": 121361,
				"name": "Game of Thrones",
				"originalLanguage": "eng",
				"overview": "Seven noble families fight for control.",
				"image": "/banners/posters/121361-1.jpg",
				"firstAired": "2011-04-17"
			}], "links": {}}`))
		}))
		defer srv.Close()

		m := tvdbimport{Source: "tvdb", Lang: "eng", URL: "https://thetvdb.com"}
		c := newTvdbTestClient(t, srv)

		var results []library.Known
		for v := range m.records(c) {
			results = append(results, v)
		}

		require.NoError(t, m.cause)
		require.Len(t, results, 1)
		got := results[0]
		require.Equal(t, "tvdb", got.Source)
		require.Equal(t, "121361", got.ID)
		require.Equal(t, "eng", got.OriginalLanguage)
		require.Equal(t, "Game of Thrones", got.OriginalTitle)
		require.Equal(t, "Game of Thrones", got.Title)
		require.Equal(t, "Seven noble families fight for control.", got.Overview)
		require.Equal(t, "https://thetvdb.com/banners/posters/121361-1.jpg", got.PosterPath)
		require.Equal(t, mimex.Video, got.Mimetype)
		require.Equal(t, time.Date(2011, 4, 17, 0, 0, 0, 0, time.UTC), got.Released)
		require.NotEmpty(t, got.UID)
		require.NotEmpty(t, got.Md5)
	})

	t.Run("skips records with no poster image", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			errorsx.Zero(fmt.Fprint(w, `{"data":[
				{"id": 1, "name": "No Poster", "originalLanguage": "eng", "overview": "missing image"},
				{"id": 2, "name": "Has Poster", "originalLanguage": "eng", "overview": "has image", "image": "/poster.jpg"}
			], "links": {}}`))
		}))
		defer srv.Close()

		m := tvdbimport{Source: "tvdb", Lang: "eng", URL: "https://thetvdb.com"}
		c := newTvdbTestClient(t, srv)

		var titles []string
		for v := range m.records(c) {
			titles = append(titles, v.Title)
		}

		require.NoError(t, m.cause)
		require.Equal(t, []string{"Has Poster"}, titles)
	})

	t.Run("prefers the alias name in the configured language", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			errorsx.Zero(fmt.Fprint(w, `{"data":[{
				"id": 5, "name": "Game of Thrones", "originalLanguage": "eng",
				"overview": "base overview", "image": "/poster.jpg",
				"aliases": [{"language": "eng", "name": "Thrones (EN Alias)"}]
			}], "links": {}}`))
		}))
		defer srv.Close()

		m := tvdbimport{Source: "tvdb", Lang: "eng", URL: "https://thetvdb.com"}
		c := newTvdbTestClient(t, srv)

		var titles []string
		for v := range m.records(c) {
			titles = append(titles, v.Title)
		}

		require.NoError(t, m.cause)
		require.Equal(t, []string{"Thrones (EN Alias)"}, titles)
	})

	t.Run("fetches the translation when original language differs and is available", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/v4/series":
				errorsx.Zero(fmt.Fprint(w, `{"data":[{
					"id": 555, "name": "Attack on Titan JP", "originalLanguage": "jpn",
					"overview": "JP overview", "image": "/poster.jpg",
					"overviewTranslations": ["eng"],
					"aliases": [{"language": "eng", "name": "Attack on Titan EN Alias"}]
				}], "links": {}}`))
			case "/v4/series/555/translations/eng":
				errorsx.Zero(fmt.Fprint(w, `{"data":{"name":"Attack on Titan","overview":"English overview from tvdb translation"}}`))
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}))
		defer srv.Close()

		m := tvdbimport{Source: "tvdb", Lang: "eng", URL: "https://thetvdb.com"}
		c := newTvdbTestClient(t, srv)

		var results []library.Known
		for v := range m.records(c) {
			results = append(results, v)
		}

		require.NoError(t, m.cause)
		require.Len(t, results, 1)
		require.Equal(t, "Attack on Titan", results[0].Title)
		require.Equal(t, "English overview from tvdb translation", results[0].Overview)
	})

	t.Run("falls back to the alias derived translation when the translation fetch fails", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/v4/series":
				w.Header().Set("Content-Type", "application/json")
				errorsx.Zero(fmt.Fprint(w, `{"data":[{
					"id": 556, "name": "Show JP", "originalLanguage": "jpn",
					"overview": "JP overview", "image": "/poster.jpg",
					"overviewTranslations": ["eng"],
					"aliases": [{"language": "eng", "name": "Show EN Alias"}]
				}], "links": {}}`))
			case "/v4/series/556/translations/eng":
				w.WriteHeader(http.StatusInternalServerError)
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
		}))
		defer srv.Close()

		m := tvdbimport{Source: "tvdb", Lang: "eng", URL: "https://thetvdb.com"}
		c := newTvdbTestClient(t, srv)

		var results []library.Known
		for v := range m.records(c) {
			results = append(results, v)
		}

		require.NoError(t, m.cause)
		require.Len(t, results, 1)
		require.Equal(t, "Show EN Alias", results[0].Title)
		require.Equal(t, "JP overview", results[0].Overview)
	})

	t.Run("released falls back to the epoch when firstAired is blank", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			errorsx.Zero(fmt.Fprint(w, `{"data":[{"id": 9, "name": "No Air Date", "originalLanguage": "eng", "overview": "x", "image": "/poster.jpg"}], "links": {}}`))
		}))
		defer srv.Close()

		m := tvdbimport{Source: "tvdb", Lang: "eng", URL: "https://thetvdb.com"}
		c := newTvdbTestClient(t, srv)

		var results []library.Known
		for v := range m.records(c) {
			results = append(results, v)
		}

		require.NoError(t, m.cause)
		require.Len(t, results, 1)
		// the epoch fallback in records() derives its date from time.Unix(0, 0)
		// in local time, so the expected day shifts with the machine's timezone.
		wantEpoch, err := time.Parse(time.DateOnly, time.Unix(0, 0).Format(time.DateOnly))
		require.NoError(t, err)
		require.Equal(t, wantEpoch, results[0].Released)
	})

	t.Run("paginates until links.next is blank", func(t *testing.T) {
		var seenPages []string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			seenPages = append(seenPages, r.URL.Query().Get("page"))
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Query().Get("page") {
			case "0":
				errorsx.Zero(fmt.Fprint(w, `{"data":[
					{"id": 1, "name": "A", "originalLanguage": "eng", "overview": "a", "image": "/a.jpg"},
					{"id": 2, "name": "B", "originalLanguage": "eng", "overview": "b", "image": "/b.jpg"}
				], "links": {"next": "1"}}`))
			case "1":
				errorsx.Zero(fmt.Fprint(w, `{"data":[
					{"id": 3, "name": "C", "originalLanguage": "eng", "overview": "c", "image": "/c.jpg"}
				], "links": {}}`))
			default:
				t.Fatalf("unexpected page requested: %s", r.URL.Query().Get("page"))
			}
		}))
		defer srv.Close()

		m := tvdbimport{Source: "tvdb", Lang: "eng", URL: "https://thetvdb.com"}
		c := newTvdbTestClient(t, srv)

		require.Equal(t, uint64(3), testx.SeqCount(m.records(c)))
		require.NoError(t, m.cause)
		require.Equal(t, []string{"0", "1"}, seenPages)
	})

	t.Run("respects the maxpage limit even when more pages are available", func(t *testing.T) {
		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			w.Header().Set("Content-Type", "application/json")
			errorsx.Zero(fmt.Fprint(w, `{"data":[{"id": 1, "name": "A", "originalLanguage": "eng", "overview": "a", "image": "/a.jpg"}], "links": {"next": "1"}}`))
		}))
		defer srv.Close()

		m := tvdbimport{Source: "tvdb", Lang: "eng", URL: "https://thetvdb.com", Limit: 1}
		c := newTvdbTestClient(t, srv)

		require.Equal(t, uint64(1), testx.SeqCount(m.records(c)))
		require.NoError(t, m.cause)
		require.Equal(t, 1, requests)
	})

	t.Run("sets the cause and stops once tvdb errors", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		m := tvdbimport{Source: "tvdb", Lang: "eng", URL: "https://thetvdb.com"}
		c := newTvdbTestClient(t, srv)

		require.Equal(t, uint64(0), testx.SeqCount(m.records(c)))
		require.Error(t, m.cause)
		require.Contains(t, m.cause.Error(), "failed to discover records")
	})

	t.Run("stops requesting further pages once the consumer breaks early", func(t *testing.T) {
		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			w.Header().Set("Content-Type", "application/json")
			errorsx.Zero(fmt.Fprint(w, `{"data":[
				{"id": 1, "name": "A", "originalLanguage": "eng", "overview": "a", "image": "/a.jpg"},
				{"id": 2, "name": "B", "originalLanguage": "eng", "overview": "b", "image": "/b.jpg"}
			], "links": {"next": "1"}}`))
		}))
		defer srv.Close()

		m := tvdbimport{Source: "tvdb", Lang: "eng", URL: "https://thetvdb.com"}
		c := newTvdbTestClient(t, srv)

		count := 0
		for range m.records(c) {
			count++
			break
		}

		require.Equal(t, 1, count)
		require.Equal(t, 1, requests, "should not request additional pages once the consumer stops iterating")
	})

	t.Run("md5 and uid are stable for the same input", func(t *testing.T) {
		newClient := func() (*tvdb.Client, func()) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				errorsx.Zero(fmt.Fprint(w, `{"data":[{"id": 121361, "name": "Game of Thrones", "originalLanguage": "eng", "overview": "x", "image": "/poster.jpg", "firstAired": "2011-04-17"}], "links": {}}`))
			}))
			return newTvdbTestClient(t, srv), srv.Close
		}

		m := tvdbimport{Source: "tvdb", Lang: "eng", URL: "https://thetvdb.com"}

		c1, close1 := newClient()
		defer close1()
		var a library.Known
		for v := range m.records(c1) {
			a = v
		}

		c2, close2 := newClient()
		defer close2()
		var b library.Known
		for v := range m.records(c2) {
			b = v
		}

		require.Equal(t, a.Md5, b.Md5)
		require.Equal(t, a.Md5Lower, b.Md5Lower)
		require.Equal(t, a.UID, b.UID)
		require.Equal(t, library.KnownImportedUintID("tvdb", uint64(121361)), a.UID)
	})
}
