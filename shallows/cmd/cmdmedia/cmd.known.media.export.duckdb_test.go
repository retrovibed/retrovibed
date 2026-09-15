package cmdmedia

import (
	"bytes"
	"sort"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestDuckdbexportRun(t *testing.T) {
	t.Run("empty database returns no error and no output", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var out bytes.Buffer
		require.NoError(t, duckdbexport{Offset: uuid.Nil.String()}.run(ctx, db, &out))
		require.Equal(t, 0, out.Len())
	})

	t.Run("exports all records with no filters", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults))

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a))
		require.NoError(t, knownimport{}.run(ctx, db, &input))
		input.Reset()
		require.NoError(t, jsonl.NewEncoder(&input).Encode(b))
		require.NoError(t, knownimport{}.run(ctx, db, &input))

		var out bytes.Buffer
		require.NoError(t, duckdbexport{Offset: uuid.Nil.String()}.run(ctx, db, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.ElementsMatch(t, []string{a.UID, b.UID}, []string{got[0].UID, got[1].UID})
	})

	t.Run("orders results by ascending UID", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults))
		a.UID = "b0000000-0000-0000-0000-000000000000"
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults))
		b.UID = "a0000000-0000-0000-0000-000000000000"

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a))
		require.NoError(t, knownimport{}.run(ctx, db, &input))
		input.Reset()
		require.NoError(t, jsonl.NewEncoder(&input).Encode(b))
		require.NoError(t, knownimport{}.run(ctx, db, &input))

		var out bytes.Buffer
		require.NoError(t, duckdbexport{Offset: uuid.Nil.String()}.run(ctx, db, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.Len(t, got, 2)
		require.Equal(t, b.UID, got[0].UID)
		require.Equal(t, a.UID, got[1].UID)
		require.True(t, sort.StringsAreSorted([]string{got[0].UID, got[1].UID}))
	})

	t.Run("Offset excludes records at or before the offset", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults))
		a.UID = "a0000000-0000-0000-0000-000000000000"
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults))
		b.UID = "c0000000-0000-0000-0000-000000000000"

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a))
		require.NoError(t, knownimport{}.run(ctx, db, &input))
		input.Reset()
		require.NoError(t, jsonl.NewEncoder(&input).Encode(b))
		require.NoError(t, knownimport{}.run(ctx, db, &input))

		var out bytes.Buffer
		require.NoError(t, duckdbexport{Offset: "b0000000-0000-0000-0000-000000000000"}.run(ctx, db, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.Len(t, got, 1)
		require.Equal(t, b.UID, got[0].UID)
	})

	t.Run("ID restricts export to the specified known ids", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults))

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a))
		require.NoError(t, knownimport{}.run(ctx, db, &input))
		input.Reset()
		require.NoError(t, jsonl.NewEncoder(&input).Encode(b))
		require.NoError(t, knownimport{}.run(ctx, db, &input))

		var out bytes.Buffer
		require.NoError(t, duckdbexport{Offset: uuid.Nil.String(), ID: []string{a.UID}}.run(ctx, db, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.Len(t, got, 1)
		require.Equal(t, a.UID, got[0].UID)
	})

	t.Run("Language restricts export to matching original_language", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults))
		a.OriginalLanguage = "fr"
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults))
		b.OriginalLanguage = "de"

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a))
		require.NoError(t, knownimport{}.run(ctx, db, &input))
		input.Reset()
		require.NoError(t, jsonl.NewEncoder(&input).Encode(b))
		require.NoError(t, knownimport{}.run(ctx, db, &input))

		var out bytes.Buffer
		require.NoError(t, duckdbexport{Offset: uuid.Nil.String(), Language: "fr"}.run(ctx, db, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.Len(t, got, 1)
		require.Equal(t, a.UID, got[0].UID)
	})

	t.Run("Mimetype restricts export to matching mimetype", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults, library.KnownOptionMimetype("video/mp4")))
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults, library.KnownOptionMimetype("audio/mpeg")))

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a))
		require.NoError(t, knownimport{}.run(ctx, db, &input))
		input.Reset()
		require.NoError(t, jsonl.NewEncoder(&input).Encode(b))
		require.NoError(t, knownimport{}.run(ctx, db, &input))

		var out bytes.Buffer
		require.NoError(t, duckdbexport{Offset: uuid.Nil.String(), Mimetype: "video/mp4"}.run(ctx, db, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.Len(t, got, 1)
		require.Equal(t, a.UID, got[0].UID)
	})

	t.Run("Source restricts export to matching source", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults, library.KnownOptionSource("tmdb")))
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults, library.KnownOptionSource("tvdb")))

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a))
		require.NoError(t, knownimport{}.run(ctx, db, &input))
		input.Reset()
		require.NoError(t, jsonl.NewEncoder(&input).Encode(b))
		require.NoError(t, knownimport{}.run(ctx, db, &input))

		var out bytes.Buffer
		require.NoError(t, duckdbexport{Offset: uuid.Nil.String(), Source: []string{"tmdb"}}.run(ctx, db, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.Len(t, got, 1)
		require.Equal(t, a.UID, got[0].UID)
	})

	t.Run("Explicit defaults to excluding adult content", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults))
		a.Adult = false
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults))
		b.Adult = true

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a))
		require.NoError(t, knownimport{}.run(ctx, db, &input))
		input.Reset()
		require.NoError(t, jsonl.NewEncoder(&input).Encode(b))
		require.NoError(t, knownimport{}.run(ctx, db, &input))

		var out bytes.Buffer
		require.NoError(t, duckdbexport{Offset: uuid.Nil.String(), Explicit: false}.run(ctx, db, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.Len(t, got, 1)
		require.Equal(t, a.UID, got[0].UID)
	})

	t.Run("Explicit=true includes adult content", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults))
		a.Adult = false
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults))
		b.Adult = true

		var input bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&input).Encode(a))
		require.NoError(t, knownimport{}.run(ctx, db, &input))
		input.Reset()
		require.NoError(t, jsonl.NewEncoder(&input).Encode(b))
		require.NoError(t, knownimport{}.run(ctx, db, &input))

		var out bytes.Buffer
		require.NoError(t, duckdbexport{Offset: uuid.Nil.String(), Explicit: true}.run(ctx, db, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.ElementsMatch(t, []string{a.UID, b.UID}, []string{got[0].UID, got[1].UID})
	})

	t.Run("Limit caps the number of exported records", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		for range 3 {
			var known library.Known
			require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))

			var input bytes.Buffer
			require.NoError(t, jsonl.NewEncoder(&input).Encode(known))
			require.NoError(t, knownimport{}.run(ctx, db, &input))
		}

		var out bytes.Buffer
		require.NoError(t, duckdbexport{Offset: uuid.Nil.String(), Limit: 2}.run(ctx, db, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.Len(t, got, 2)
	})
}
