package cmdmedia

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
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

		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			page := r.URL.Query().Get("page")
			switch page {
			case "1":
				_ = errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":2,"total_pages":2,"results":[{"id":1,"name":"Show One"}]}`))
			case "2":
				_ = errorsx.Zero(fmt.Fprint(w, `{"page":2,"total_results":2,"total_pages":2,"results":[{"id":2,"name":"Show Two"}]}`))
			default:
				t.Fatalf("unexpected page requested: %s", page)
			}
		}))
		defer srv.Close()

		day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{StartAt: day, EndAt: day, Attempts: 1}

		var titles []string
		for v := range tm.series(ctx, newTmdbTestClient(t, srv)) {
			titles = append(titles, v.Title)
		}

		require.NoError(t, tm.cause)
		require.Equal(t, []string{"Show One", "Show Two"}, titles)
		require.Equal(t, 2, requests)
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

		requests := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			errorsx.Zero(fmt.Fprint(w, `{"page":1,"total_results":2,"total_pages":2,"results":[{"id":1,"name":"Show One"},{"id":2,"name":"Show Two"}]}`))
		}))
		defer srv.Close()

		day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		tm := &tmdbimport{StartAt: day, EndAt: day, Attempts: 1}

		count := 0
		for range tm.series(ctx, newTmdbTestClient(t, srv)) {
			count++
			break
		}

		require.Equal(t, 1, count)
		require.Equal(t, 1, requests, "should not request additional pages once the consumer stops iterating")
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
}
