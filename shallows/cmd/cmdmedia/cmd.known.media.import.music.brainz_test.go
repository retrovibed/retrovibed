package cmdmedia

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/michiwend/gomusicbrainz"

	"github.com/retrovibed/retrovibed/internal/backoffx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestMBImportReleases(t *testing.T) {
	testDate := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	unlimited := rate.NewLimiter(rate.Inf, 1)
	immediate := backoffx.Constant(0)

	releaseEntryXML := func(id, title, date, lang string) string {
		return fmt.Sprintf(`<release id="%s"><title>%s</title><date>%s</date><text-representation><language>%s</language></text-representation><release-group id="%s"><title>%s</title></release-group></release>`, id, title, date, lang, id, title)
	}

	releaseXML := func(count, offset int, releases ...string) string {
		var b strings.Builder
		fmt.Fprintf(&b, `<?xml version="1.0"?><metadata><release-list count="%d" offset="%d">`, count, offset)
		for _, r := range releases {
			errorsx.Must(b.WriteString(r))
		}
		errorsx.Must(b.WriteString(`</release-list></metadata>`))
		return b.String()
	}

	// fullPage generates n release entries to simulate a full page (limit=100).
	fullPage := func(n int) string {
		var b strings.Builder
		for i := range n {
			errorsx.Must(b.WriteString(releaseEntryXML(
				fmt.Sprintf("8f6a4a2b-e29b-41d4-a716-%012d", i),
				fmt.Sprintf("Album %d", i),
				"2020-01-15",
				"",
			)))
		}
		return b.String()
	}

	t.Run("returns releases from a single page", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/ws/2/release", r.URL.Path)
			require.Equal(t, "date:2020-01-15", r.URL.Query().Get("query"))
			errorsx.Must(fmt.Fprint(w, releaseXML(2, 0,
				releaseEntryXML("8f6a4a2b-e29b-41d4-a716-446655440001", "Album One", "2020-01-15", "eng"),
				releaseEntryXML("8f6a4a2b-e29b-41d4-a716-446655440002", "Album Two", "2020-01-15", "fra"),
			)))
		}))
		defer srv.Close()

		c, err := gomusicbrainz.NewWS2Client(
			testx.Must(url.Parse(srv.URL))(t).String(),
			"TestApp",
			"0.0.1",
			"test@example.com",
		)
		require.NoError(t, err)

		m := mbimport{
			StartAt:  testDate,
			EndAt:    testDate,
			Source:   "musicbrainz",
			Attempts: 1,
		}

		type result struct {
			title, originalTitle, lang, posterPath string
		}
		var results []result
		for known := range m.releases(ctx, c, unlimited, immediate) {
			results = append(results, result{known.Title, known.OriginalTitle, known.OriginalLanguage, known.PosterPath})
		}

		require.NoError(t, m.cause)
		require.Equal(t, []result{
			{"Album One", "Album One", "en", "https://coverartarchive.org/release-group/8f6a4a2b-e29b-41d4-a716-446655440001/front-500"},
			{"Album Two", "Album Two", "fr", "https://coverartarchive.org/release-group/8f6a4a2b-e29b-41d4-a716-446655440002/front-500"},
		}, results)
	})

	t.Run("paginates across multiple pages", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		// Pagination continues when a full page (limit=100) is returned and count > offset.
		const totalCount = 101
		callCount := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			offset := r.URL.Query().Get("offset")
			if offset == "0" || offset == "" {
				errorsx.Must(fmt.Fprint(w, releaseXML(totalCount, 0, fullPage(100))))
			} else {
				errorsx.Must(fmt.Fprint(w, releaseXML(totalCount, 100,
					releaseEntryXML("8f6a4a2b-e29b-41d4-a716-000000000999", "Last Album", "2020-01-15", ""),
				)))
			}
		}))
		defer srv.Close()

		c, err := gomusicbrainz.NewWS2Client(
			testx.Must(url.Parse(srv.URL))(t).String(),
			"TestApp",
			"0.0.1",
			"test@example.com",
		)
		require.NoError(t, err)

		m := mbimport{
			StartAt:  testDate,
			EndAt:    testDate,
			Source:   "musicbrainz",
			Attempts: 1,
		}

		require.Equal(t, uint64(totalCount), testx.SeqCount(m.releases(ctx, c, unlimited, immediate)))
		require.NoError(t, m.cause)
		require.Equal(t, 2, callCount)
	})

	t.Run("stops on empty release list", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorsx.Must(fmt.Fprint(w, releaseXML(0, 0)))
		}))
		defer srv.Close()

		c, err := gomusicbrainz.NewWS2Client(
			testx.Must(url.Parse(srv.URL))(t).String(),
			"TestApp",
			"0.0.1",
			"test@example.com",
		)
		require.NoError(t, err)

		m := mbimport{
			StartAt:  testDate,
			EndAt:    testDate,
			Source:   "musicbrainz",
			Attempts: 1,
		}

		require.Equal(t, uint64(0), testx.SeqCount(m.releases(ctx, c, unlimited, immediate)))
		require.NoError(t, m.cause)
	})

	t.Run("iterates over multiple days", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		seenDates := map[string]bool{}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query().Get("query")
			seenDates[q] = true
			id := fmt.Sprintf("8f6a4a2b-e29b-41d4-a716-44665544000%d", len(seenDates))
			errorsx.Must(fmt.Fprint(w, releaseXML(1, 0, releaseEntryXML(id, q, "2020-01-15", ""))))
		}))
		defer srv.Close()

		c, err := gomusicbrainz.NewWS2Client(
			testx.Must(url.Parse(srv.URL))(t).String(),
			"TestApp",
			"0.0.1",
			"test@example.com",
		)
		require.NoError(t, err)

		m := mbimport{
			StartAt:  testDate,
			EndAt:    testDate.AddDate(0, 0, 2),
			Source:   "musicbrainz",
			Attempts: 1,
		}

		require.Equal(t, uint64(3), testx.SeqCount(m.releases(ctx, c, unlimited, immediate)))
		require.NoError(t, m.cause)
		require.True(t, seenDates["date:2020-01-15"])
		require.True(t, seenDates["date:2020-01-16"])
		require.True(t, seenDates["date:2020-01-17"])
	})

	t.Run("sets cause on server error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorsx.Must(fmt.Fprint(w, "not valid xml"))
		}))
		defer srv.Close()

		c, err := gomusicbrainz.NewWS2Client(
			testx.Must(url.Parse(srv.URL))(t).String(),
			"TestApp",
			"0.0.1",
			"test@example.com",
		)
		require.NoError(t, err)

		m := mbimport{
			StartAt:  testDate,
			EndAt:    testDate,
			Source:   "musicbrainz",
			Attempts: 0, // stop after the first failed attempt
		}

		require.Equal(t, uint64(0), testx.SeqCount(m.releases(ctx, c, unlimited, immediate)))
		require.Error(t, m.cause)
	})
}
