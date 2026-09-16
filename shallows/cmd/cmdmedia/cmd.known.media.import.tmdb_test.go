package cmdmedia

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func newTmdbTestClient(t *testing.T, srv *httptest.Server) *tmdb.Client {
	t.Helper()
	c, err := tmdb.Init("test-key")
	require.NoError(t, err)
	c.SetCustomBaseURL(srv.URL)
	return c
}

func TestTmdbImportImgpath(t *testing.T) {
	t.Run("blank path stays blank", func(t *testing.T) {
		tm := tmdbimport{URL: "https://image.tmdb.org/t/p/original"}
		require.Equal(t, "", tm.imgpath(""))
	})

	t.Run("non-blank path is prefixed with the configured base url", func(t *testing.T) {
		tm := tmdbimport{URL: "https://image.tmdb.org/t/p/original"}
		require.Equal(t, "https://image.tmdb.org/t/p/original/poster.jpg", tm.imgpath("/poster.jpg"))
	})
}

func TestTmdbImportSeries(t *testing.T) {
	t.Run("advances past a zero-result date instead of exceeding tmdb's page limit", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			if requests > 500 {
				// mirrors the real TMDB error body that triggered the production incident.
				w.WriteHeader(http.StatusInternalServerError)
				_ = errorsx.Zero(fmt.Fprint(w, `{"status_code":22,"status_message":"Invalid page: Pages start at 1 and max at 500. They are expected to be an integer.","success":false}`))
				return
			}
			errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":0,"total_pages":0,"results":[]}`))
		}))
		defer srv.Close()

		day := time.Date(1946, 5, 15, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{StartAt: day, EndAt: day, Attempts: 1}

		for range tm.series(ctx, newTmdbTestClient(t, srv)) {
			t.Fatal("expected no results for a date with zero matches")
		}

		require.NoError(t, tm.cause)
		require.Equal(t, 1, requests, "a zero-result date must advance to the next day after a single page request, not keep incrementing the page")
	})

	t.Run("paginates through multiple pages before advancing the date", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var requests atomic.Int64
		routes := mux.NewRouter()
		routes.HandleFunc("/discover/tv", func(w http.ResponseWriter, r *http.Request) {
			requests.Add(1)
			page := r.URL.Query().Get("page")
			switch page {
			case "1":
				_ = errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":2,"total_pages":2,"results":[{"id":1,"name":"Show One"}]}`))
			case "2":
				_ = errorsx.Zero(fmt.Fprint(w, `{"page":2,"total_results":2,"total_pages":2,"results":[{"id":2,"name":"Show Two"}]}`))
			default:
				t.Fatalf("unexpected page requested: %s", page)
			}
		})
		routes.HandleFunc("/tv/{id}", func(w http.ResponseWriter, r *http.Request) {
			requests.Add(1)
			// series fetches each show's episodes concurrently via a pool of
			// workers, which fetches show details before any season/episode
			// lookups - a show with no seasons ends that fetch here.
			_ = errorsx.Zero(fmt.Fprint(w, `{"id":`+mux.Vars(r)["id"]+`,"seasons":[]}`))
		})
		srv := httptest.NewServer(routes)
		defer srv.Close()

		day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{StartAt: day, EndAt: day, Attempts: 1}

		var titles []string
		for v := range tm.series(ctx, newTmdbTestClient(t, srv)) {
			titles = append(titles, v.Title)
		}

		require.NoError(t, tm.cause)
		require.ElementsMatch(t, []string{"Show One", "Show Two"}, titles, "shows may arrive in any order once episode fetches run concurrently")
		require.Equal(t, int64(4), requests.Load(), "2 discover pages plus 1 tv-details fetch per discovered show")
	})

	t.Run("requests season details in the show's original language", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var seasonLanguage string
		routes := mux.NewRouter()
		routes.HandleFunc("/discover/tv", func(w http.ResponseWriter, r *http.Request) {
			_ = errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":1,"total_pages":1,"results":[{"id":1,"name":"Show One","original_language":"ja"}]}`))
		})
		routes.HandleFunc("/tv/{id}", func(w http.ResponseWriter, r *http.Request) {
			_ = errorsx.Zero(fmt.Fprint(w, `{"id":`+mux.Vars(r)["id"]+`,"seasons":[{"season_number":1}]}`))
		})
		routes.HandleFunc("/tv/{id}/season/{season}", func(w http.ResponseWriter, r *http.Request) {
			seasonLanguage = r.URL.Query().Get("language")
			_ = errorsx.Zero(fmt.Fprint(w, `{"episodes":[]}`))
		})
		srv := httptest.NewServer(routes)
		defer srv.Close()

		day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{StartAt: day, EndAt: day, Attempts: 1}

		for range tm.series(ctx, newTmdbTestClient(t, srv)) {
		}

		require.NoError(t, tm.cause)
		require.Equal(t, "ja", seasonLanguage, "season details should be requested in the show's original_language so TMDB can return an episode title it lacks in English")
	})

	t.Run("maps tv show fields onto the known record", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":1,"total_pages":1,"results":[{
				"id": 42,
				"name": "Mystery Theater",
				"original_name": "Mystery Theater Original",
				"original_language": "en",
				"overview": "A spooky anthology series.",
				"first_air_date": "1946-05-15",
				"poster_path": "/poster.jpg",
				"backdrop_path": "/backdrop.jpg",
				"popularity": 12.5,
				"adult": true
			}]}`))
		}))
		defer srv.Close()

		day := time.Date(1946, 5, 15, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{Source: "tmdb", URL: "https://image.tmdb.org/t/p/original", StartAt: day, EndAt: day, Attempts: 1}

		var results []library.Known
		for v := range tm.series(ctx, newTmdbTestClient(t, srv)) {
			results = append(results, v)
		}

		require.NoError(t, tm.cause)
		require.Len(t, results, 1)
		got := results[0]
		require.Equal(t, "42", got.ID)
		require.Equal(t, "tmdb", got.Source)
		require.Equal(t, "Mystery Theater", got.Title)
		require.Equal(t, "Mystery Theater Original", got.OriginalTitle)
		require.Equal(t, "en", got.OriginalLanguage)
		require.Equal(t, "A spooky anthology series.", got.Overview)
		require.Equal(t, "https://image.tmdb.org/t/p/original/poster.jpg", got.PosterPath)
		require.Equal(t, "https://image.tmdb.org/t/p/original/backdrop.jpg", got.BackdropPath)
		require.Equal(t, 12.5, got.Popularity)
		require.True(t, got.Adult)
		require.Equal(t, mimex.Video, got.Mimetype)
		require.Equal(t, day, got.Released)
		require.NotEmpty(t, got.UID)
		require.NotEmpty(t, got.Md5)
		// parent_uid is a NOT NULL UUID column; a show is never anyone's
		// child, so this must be the nil UUID, not Go's zero-value "".
		require.Equal(t, uuid.Nil.String(), got.ParentUID)
	})

	t.Run("a movie, a show, and an episode sharing the same tmdb numeric id never collide into the same UID", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		// TMDB movie ids, tv show ids, and tv episode ids are independent
		// numeric sequences - real-world example: tmdb movie 4271 is "Life
		// Is a Long Quiet River" while tmdb tv show 4271 is Farscape. Using
		// the same id (4271) for all three here reproduces that collision.
		routes := mux.NewRouter()
		routes.HandleFunc("/discover/movie", func(w http.ResponseWriter, r *http.Request) {
			_ = errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":1,"total_pages":1,"results":[{"id":4271,"title":"Life Is a Long Quiet River"}]}`))
		})
		routes.HandleFunc("/discover/tv", func(w http.ResponseWriter, r *http.Request) {
			_ = errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":1,"total_pages":1,"results":[{"id":4271,"name":"Farscape"}]}`))
		})
		routes.HandleFunc("/tv/4271", func(w http.ResponseWriter, r *http.Request) {
			_ = errorsx.Zero(fmt.Fprint(w, `{"id":4271,"seasons":[{"season_number":1}]}`))
		})
		routes.HandleFunc("/tv/4271/season/1", func(w http.ResponseWriter, r *http.Request) {
			_ = errorsx.Zero(fmt.Fprint(w, `{"episodes":[{"id":4271,"season_number":1,"episode_number":1,"name":"Mind the Baby"}]}`))
		})
		srv := httptest.NewServer(routes)
		defer srv.Close()

		day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{Source: "tmdb", StartAt: day, EndAt: day, Attempts: 1}
		client := newTmdbTestClient(t, srv)

		var movieUID, showUID, episodeUID string
		for v := range tm.movies(ctx, client) {
			movieUID = v.UID
		}
		require.NoError(t, tm.cause)
		require.NotEmpty(t, movieUID)

		for v := range tm.series(ctx, client) {
			if v.ParentUID == uuid.Nil.String() {
				showUID = v.UID
			} else {
				episodeUID = v.UID
			}
		}
		require.NoError(t, tm.cause)
		require.NotEmpty(t, showUID)
		require.NotEmpty(t, episodeUID)

		require.NotEqual(t, movieUID, showUID, "a movie and a tv show sharing the same tmdb id must not collide")
		require.NotEqual(t, movieUID, episodeUID, "a movie and a tv episode sharing the same tmdb id must not collide")
		require.NotEqual(t, showUID, episodeUID, "a tv show and a tv episode sharing the same tmdb id must not collide")
	})

	t.Run("records the cause and stops once tmdb persistently errors", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			w.WriteHeader(http.StatusInternalServerError)
			errorsx.Zero(fmt.Fprint(w, `{"status_code":34,"status_message":"The resource you requested could not be found.","success":false}`))
		}))
		defer srv.Close()

		day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		// Attempts:0 still allows one retry after the initial failure, so the
		// backoff strategy waits once before giving up - keep this in mind if
		// the test feels slow, it's the documented retry behavior, not a hang.
		tm := &tmdbimport{StartAt: day, EndAt: day, Attempts: 0}

		for range tm.series(ctx, newTmdbTestClient(t, srv)) {
			t.Fatal("expected no results once tmdb errors persistently")
		}

		require.Error(t, tm.cause)
		require.Contains(t, tm.cause.Error(), "failed to discover series")
	})

	t.Run("stops issuing requests once the consumer breaks early", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var requests atomic.Int64
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests.Add(1)
			errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":2,"total_pages":2,"results":[{"id":1,"name":"Show One"},{"id":2,"name":"Show Two"}]}`))
		}))
		defer srv.Close()

		day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{StartAt: day, EndAt: day, Attempts: 1, Concurrency: 8}

		count := 0
		for range tm.series(ctx, newTmdbTestClient(t, srv)) {
			count++
			break
		}

		require.Equal(t, 1, count)
		// series() fetches episode details concurrently across a bounded pool
		// of workers, so some requests already in flight for shows near the
		// front of the backlog complete even after the consumer stops - this
		// bounds that prefetch instead of requiring exactly zero extra work.
		require.LessOrEqual(t, requests.Load(), int64(1+int(tm.Concurrency)), "should not issue unbounded additional requests once the consumer stops iterating")
	})

	t.Run("requests the next page without checking whether the previous page's results were empty", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		// Same junk-date pocket TMDB behavior as the movies discovery: a day
		// can report a large total_pages up front while the actual results
		// dry up well before that count.
		var discoverRequests atomic.Int64
		routes := mux.NewRouter()
		routes.HandleFunc("/discover/tv", func(w http.ResponseWriter, r *http.Request) {
			discoverRequests.Add(1)
			day := r.URL.Query().Get("first_air_date.gte")
			page := r.URL.Query().Get("page")
			switch {
			case day == "1700-02-09" && page == "1":
				errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":1,"total_pages":500,"results":[{"id":1,"name":"Show One"}]}`))
			case day == "1700-02-09" && page == "2":
				errorsx.Zero(fmt.Fprint(w, `{"page":2,"total_results":1,"total_pages":500,"results":[]}`))
			case day == "1700-02-10" && page == "1":
				errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":0,"total_pages":0,"results":[]}`))
			default:
				t.Fatalf("unexpected request: day=%s page=%s", day, page)
			}
		})
		routes.HandleFunc("/tv/{id}", func(w http.ResponseWriter, r *http.Request) {
			errorsx.Zero(fmt.Fprint(w, `{"id":`+mux.Vars(r)["id"]+`,"seasons":[]}`))
		})
		srv := httptest.NewServer(routes)
		defer srv.Close()

		start := time.Date(1700, 2, 9, 0, 0, 0, 0, time.UTC)
		end := time.Date(1700, 2, 10, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{StartAt: start, EndAt: end, Attempts: 1}

		var titles []string
		for v := range tm.series(ctx, newTmdbTestClient(t, srv)) {
			titles = append(titles, v.Title)
		}

		require.NoError(t, tm.cause)
		require.Equal(t, []string{"Show One"}, titles)
		require.Equal(t, int64(3), discoverRequests.Load(), "should advance to the next day after the first empty page instead of paginating all the way to a bogus total_pages")
	})
}

func TestTmdbImportMovies(t *testing.T) {
	t.Run("advances past a zero-result date instead of exceeding tmdb's page limit", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			if requests > 500 {
				w.WriteHeader(http.StatusInternalServerError)
				errorsx.Zero(fmt.Fprint(w, `{"status_code":22,"status_message":"Invalid page: Pages start at 1 and max at 500. They are expected to be an integer.","success":false}`))
				return
			}
			errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":0,"total_pages":0,"results":[]}`))
		}))
		defer srv.Close()

		day := time.Date(1946, 5, 15, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{StartAt: day, EndAt: day, Attempts: 1}

		for range tm.movies(ctx, newTmdbTestClient(t, srv)) {
			t.Fatal("expected no results for a date with zero matches")
		}

		require.NoError(t, tm.cause)
		require.Equal(t, 1, requests, "a zero-result date must advance to the next day after a single page request, not keep incrementing the page")
	})

	t.Run("paginates through multiple pages before advancing the date", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			page := r.URL.Query().Get("page")
			switch page {
			case "1":
				errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":2,"total_pages":2,"results":[{"id":1,"title":"Movie One"}]}`))
			case "2":
				errorsx.Zero(fmt.Fprint(w, `{"page":2,"total_results":2,"total_pages":2,"results":[{"id":2,"title":"Movie Two"}]}`))
			default:
				t.Fatalf("unexpected page requested: %s", page)
			}
		}))
		defer srv.Close()

		day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{StartAt: day, EndAt: day, Attempts: 1}

		var titles []string
		for v := range tm.movies(ctx, newTmdbTestClient(t, srv)) {
			titles = append(titles, v.Title)
		}

		require.NoError(t, tm.cause)
		require.Equal(t, []string{"Movie One", "Movie Two"}, titles)
		require.Equal(t, 2, requests)
	})

	t.Run("maps movie fields onto the known record", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":1,"total_pages":1,"results":[{
				"id": 7,
				"title": "The Great Heist",
				"original_title": "The Great Heist Original",
				"original_language": "en",
				"overview": "A daring crew pulls off the impossible.",
				"release_date": "1946-05-15",
				"poster_path": "/poster.jpg",
				"backdrop_path": "/backdrop.jpg",
				"popularity": 8.25,
				"adult": false
			}]}`))
		}))
		defer srv.Close()

		day := time.Date(1946, 5, 15, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{Source: "tmdb", URL: "https://image.tmdb.org/t/p/original", StartAt: day, EndAt: day, Attempts: 1}

		var results []library.Known
		for v := range tm.movies(ctx, newTmdbTestClient(t, srv)) {
			results = append(results, v)
		}

		require.NoError(t, tm.cause)
		require.Len(t, results, 1)
		got := results[0]
		require.Equal(t, "7", got.ID)
		require.Equal(t, "tmdb", got.Source)
		require.Equal(t, "The Great Heist", got.Title)
		require.Equal(t, "The Great Heist Original", got.OriginalTitle)
		require.Equal(t, "en", got.OriginalLanguage)
		require.Equal(t, "A daring crew pulls off the impossible.", got.Overview)
		require.Equal(t, "https://image.tmdb.org/t/p/original/poster.jpg", got.PosterPath)
		require.Equal(t, "https://image.tmdb.org/t/p/original/backdrop.jpg", got.BackdropPath)
		require.Equal(t, 8.25, got.Popularity)
		require.False(t, got.Adult)
		require.Equal(t, mimex.Video, got.Mimetype)
		require.Equal(t, day, got.Released)
		require.NotEmpty(t, got.UID)
		require.NotEmpty(t, got.Md5)
		// parent_uid is a NOT NULL UUID column; a movie is never anyone's
		// child, so this must be the nil UUID, not Go's zero-value "".
		require.Equal(t, uuid.Nil.String(), got.ParentUID)
	})

	t.Run("records the cause and stops once tmdb persistently errors", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			errorsx.Zero(fmt.Fprint(w, `{"status_code":34,"status_message":"The resource you requested could not be found.","success":false}`))
		}))
		defer srv.Close()

		day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{StartAt: day, EndAt: day, Attempts: 0}

		for range tm.movies(ctx, newTmdbTestClient(t, srv)) {
			t.Fatal("expected no results once tmdb errors persistently")
		}

		require.Error(t, tm.cause)
		require.Contains(t, tm.cause.Error(), "failed to discover movies")
	})

	t.Run("requests the next page without checking whether the previous page's results were empty", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		// TMDB has junk movie entries sharing a single placeholder
		// release_date (e.g. 1700-xx-xx). Discovering one of those dates can
		// report a large total_pages up front, but the actual results dry up
		// well before that count - production saw this drive a single date
		// through 70+ page requests because nothing checks for an empty
		// results page before requesting the next one.
		var requests atomic.Int64
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests.Add(1)
			page := r.URL.Query().Get("page")
			switch page {
			case "1":
				errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":1,"total_pages":500,"results":[{"id":1,"title":"Movie One"}]}`))
			case "2":
				errorsx.Zero(fmt.Fprint(w, `{"page":2,"total_results":1,"total_pages":500,"results":[]}`))
			default:
				t.Fatalf("unexpected page requested: %s", page)
			}
		}))
		defer srv.Close()

		day := time.Date(1700, 2, 9, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{StartAt: day, EndAt: day, Attempts: 1}

		var titles []string
		for v := range tm.movies(ctx, newTmdbTestClient(t, srv)) {
			titles = append(titles, v.Title)
		}

		require.NoError(t, tm.cause)
		require.Equal(t, []string{"Movie One"}, titles)
		require.Equal(t, int64(2), requests.Load(), "should stop at the first empty page instead of paginating all the way to a bogus total_pages")
	})
}
