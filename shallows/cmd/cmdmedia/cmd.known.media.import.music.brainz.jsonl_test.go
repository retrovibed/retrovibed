package cmdmedia

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/testx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestMBJSONLImportReleases(t *testing.T) {
	m := &mbjsonlimport{Source: "musicbrainz"}

	encode := func(v any) string {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	releaseJSON := func(id, title string, rg mbJSONReleaseGroup, lang string) string {
		return encode(mbJSONRelease{
			ID:                 id,
			Title:              title,
			TextRepresentation: mbJSONTextRepr{Language: lang},
			ReleaseGroup:       rg,
		})
	}

	collect := func(t *testing.T, lines ...string) []library.Known {
		t.Helper()
		ctx, done := testx.Context(t)
		defer done()
		r := strings.NewReader(strings.Join(lines, "\n"))
		var results []library.Known
		for v := range m.releases(ctx, r) {
			results = append(results, v)
		}
		return results
	}

	t.Run("maps title and language from release group", func(t *testing.T) {
		line := releaseJSON("rel-id", "Release Title", mbJSONReleaseGroup{
			ID:    "rg-8f6a4a2b-e29b-41d4-a716-446655440001",
			Title: "Release Group Title",
		}, "eng")
		results := collect(t, line)
		require.Len(t, results, 1)
		require.Equal(t, "Release Group Title", results[0].Title)
		require.Equal(t, "Release Group Title", results[0].OriginalTitle)
		require.Equal(t, "en", results[0].OriginalLanguage)
		require.Equal(t, mimex.Audio, results[0].Mimetype)
		require.Equal(t, "musicbrainz", results[0].Source)
	})

	t.Run("falls back to release title when release group title is absent", func(t *testing.T) {
		line := releaseJSON("rel-id", "Release Title", mbJSONReleaseGroup{ID: "rg-id"}, "")
		results := collect(t, line)
		require.Len(t, results, 1)
		require.Equal(t, "Release Title", results[0].Title)
	})

	t.Run("falls back to release ID when release group ID is absent", func(t *testing.T) {
		line := releaseJSON("8f6a4a2b-e29b-41d4-a716-446655440002", "Album", mbJSONReleaseGroup{}, "")
		results := collect(t, line)
		require.Len(t, results, 1)
		require.Equal(t, "8f6a4a2b-e29b-41d4-a716-446655440002", results[0].ID)
	})

	t.Run("uses release group ID for deduplication key", func(t *testing.T) {
		line := releaseJSON("rel-id", "Album", mbJSONReleaseGroup{
			ID:    "8f6a4a2b-e29b-41d4-a716-446655440003",
			Title: "Album",
		}, "")
		results := collect(t, line)
		require.Equal(t, "8f6a4a2b-e29b-41d4-a716-446655440003", results[0].ID)
	})

	t.Run("cover art URL uses release group ID", func(t *testing.T) {
		line := releaseJSON("rel-id", "Album", mbJSONReleaseGroup{
			ID:    "8f6a4a2b-e29b-41d4-a716-446655440004",
			Title: "Album",
		}, "")
		results := collect(t, line)
		require.Equal(t,
			"https://coverartarchive.org/release-group/8f6a4a2b-e29b-41d4-a716-446655440004/front-500",
			results[0].PosterPath,
		)
	})

	t.Run("parses full date", func(t *testing.T) {
		line := releaseJSON("rel-id", "Album", mbJSONReleaseGroup{
			ID:               "rg-id",
			FirstReleaseDate: "2020-06-15",
		}, "")
		results := collect(t, line)
		require.Equal(t, time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC), results[0].Released)
	})

	t.Run("parses year-month date", func(t *testing.T) {
		line := releaseJSON("rel-id", "Album", mbJSONReleaseGroup{
			ID:               "rg-id",
			FirstReleaseDate: "2020-06",
		}, "")
		results := collect(t, line)
		require.Equal(t, time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC), results[0].Released)
	})

	t.Run("parses year-only date", func(t *testing.T) {
		line := releaseJSON("rel-id", "Album", mbJSONReleaseGroup{
			ID:               "rg-id",
			FirstReleaseDate: "2020",
		}, "")
		results := collect(t, line)
		require.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), results[0].Released)
	})

	t.Run("zero time for absent date", func(t *testing.T) {
		line := releaseJSON("rel-id", "Album", mbJSONReleaseGroup{ID: "rg-id"}, "")
		results := collect(t, line)
		require.True(t, results[0].Released.IsZero())
	})

	t.Run("md5 is stable for same input", func(t *testing.T) {
		line := releaseJSON("8f6a4a2b-e29b-41d4-a716-446655440005", "Stable Album", mbJSONReleaseGroup{
			ID: "rg-8f6a4a2b-e29b-41d4-a716-446655440005",
		}, "eng")
		a := collect(t, line)
		b := collect(t, line)
		require.Equal(t, a[0].Md5, b[0].Md5)
		require.Equal(t, a[0].Md5Lower, b[0].Md5Lower)
	})

	t.Run("yields multiple records", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		lines := strings.Join([]string{
			releaseJSON("8f6a4a2b-e29b-41d4-a716-446655440006", "Album One", mbJSONReleaseGroup{
				ID: "rg-8f6a4a2b-e29b-41d4-a716-446655440006",
			}, "eng"),
			releaseJSON("8f6a4a2b-e29b-41d4-a716-446655440007", "Album Two", mbJSONReleaseGroup{
				ID: "rg-8f6a4a2b-e29b-41d4-a716-446655440007",
			}, "fra"),
		}, "\n")
		require.Equal(t, uint64(2), testx.SeqCount(m.releases(ctx, strings.NewReader(lines))))
		require.NoError(t, m.cause)
	})

	t.Run("empty input yields no records", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		require.Equal(t, uint64(0), testx.SeqCount(m.releases(ctx, strings.NewReader(""))))
		require.NoError(t, m.cause)
	})

	t.Run("sets cause on malformed json", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		require.Equal(t, uint64(0), testx.SeqCount(m.releases(ctx, strings.NewReader("not valid json"))))
		require.Error(t, m.cause)
	})
}

func TestMBJSONLImportRun(t *testing.T) {
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

	t.Run("encodes records to output", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		input := strings.NewReader(`{"id":"8f6a4a2b-e29b-41d4-a716-446655440010","title":"Run Album","text-representation":{"language":"eng"},"release-group":{"id":"8f6a4a2b-e29b-41d4-a716-446655440010","title":"Run Album","first-release-date":"2021-03-01"}}`)
		var buf bytes.Buffer
		m := mbjsonlimport{Source: "musicbrainz"}
		require.NoError(t, m.run(ctx, input, jsonl.NewEncoder(&buf)))

		results := decodeAll(t, &buf)
		require.Len(t, results, 1)
		require.Equal(t, "Run Album", results[0].Title)
		require.Equal(t, "en", results[0].OriginalLanguage)
		require.Equal(t, mimex.Audio, results[0].Mimetype)
	})

	t.Run("returns error on malformed input", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		m := mbjsonlimport{Source: "musicbrainz"}
		require.Error(t, m.run(ctx, strings.NewReader("not json"), jsonl.NewEncoder(&bytes.Buffer{})))
	})
}
