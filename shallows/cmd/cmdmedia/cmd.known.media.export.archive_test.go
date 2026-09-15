package cmdmedia

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/tarx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

// tarchiveexport.run resolves each archive path found by walking t.Pattern
// with a bare os.Open(path), where path is relative to t.Pattern rather than
// to the process's working directory. That only works when t.Pattern is "."
// (i.e. the working directory itself), so every subtest here chdirs into the
// archive directory and sets Pattern: "." to match.
func TestTarchiveexportRun(t *testing.T) {
	t.Run("empty directory produces no output", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()
		t.Chdir(dir)

		var out bytes.Buffer
		require.NoError(t, tarchiveexport{Pattern: "."}.run(ctx, &out))
		require.Equal(t, 0, out.Len())
	})

	t.Run("skips files without a .tar.gz suffix", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()
		t.Chdir(dir)
		require.NoError(t, os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("not an archive"), 0600))

		var out bytes.Buffer
		require.NoError(t, tarchiveexport{Pattern: "."}.run(ctx, &out))
		require.Equal(t, 0, out.Len())
	})

	t.Run("exports the single record in a single archive", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Alpha"
		known.OriginalTitle = "Alpha Original"
		known.Overview = "An overview"

		src := t.TempDir()
		var record bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&record).Encode(known))
		require.NoError(t, os.WriteFile(filepath.Join(src, "record"), record.Bytes(), 0600))

		dir := t.TempDir()
		archive, err := os.Create(filepath.Join(dir, "a.tar.gz"))
		require.NoError(t, err)
		require.NoError(t, tarx.Pack(archive, src))
		require.NoError(t, archive.Close())

		t.Chdir(dir)

		var out bytes.Buffer
		require.NoError(t, tarchiveexport{Pattern: "."}.run(ctx, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.Len(t, got, 1)
		require.Equal(t, known.UID, got[0].UID)
		require.Equal(t, "Alpha\nAlpha Original\nAn overview", got[0].AutoDescription)
	})

	t.Run("exports records from multiple archives", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		dir := t.TempDir()

		var a, b library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults))

		for name, known := range map[string]library.Known{"a.tar.gz": a, "b.tar.gz": b} {
			src := t.TempDir()
			var record bytes.Buffer
			require.NoError(t, jsonl.NewEncoder(&record).Encode(known))
			require.NoError(t, os.WriteFile(filepath.Join(src, "record"), record.Bytes(), 0600))

			archive, err := os.Create(filepath.Join(dir, name))
			require.NoError(t, err)
			require.NoError(t, tarx.Pack(archive, src))
			require.NoError(t, archive.Close())
		}

		t.Chdir(dir)

		var out bytes.Buffer
		require.NoError(t, tarchiveexport{Pattern: "."}.run(ctx, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.Len(t, got, 2)
		require.ElementsMatch(t, []string{a.UID, b.UID}, []string{got[0].UID, got[1].UID})
	})

	t.Run("exports every record within a single tar member", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var first, second library.Known
		require.NoError(t, testx.Fake(&first, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&second, library.KnownOptionTestDefaults))

		src := t.TempDir()
		var record bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&record).Encode(first))
		require.NoError(t, jsonl.NewEncoder(&record).Encode(second))
		require.NoError(t, os.WriteFile(filepath.Join(src, "record"), record.Bytes(), 0600))

		dir := t.TempDir()
		archive, err := os.Create(filepath.Join(dir, "a.tar.gz"))
		require.NoError(t, err)
		require.NoError(t, tarx.Pack(archive, src))
		require.NoError(t, archive.Close())

		t.Chdir(dir)

		var out bytes.Buffer
		require.NoError(t, tarchiveexport{Pattern: "."}.run(ctx, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.Len(t, got, 2)
		require.ElementsMatch(t, []string{first.UID, second.UID}, []string{got[0].UID, got[1].UID})
	})

	t.Run("returns an error but still exports valid archives when one archive is corrupt", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))

		src := t.TempDir()
		var record bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&record).Encode(known))
		require.NoError(t, os.WriteFile(filepath.Join(src, "record"), record.Bytes(), 0600))

		dir := t.TempDir()
		archive, err := os.Create(filepath.Join(dir, "good.tar.gz"))
		require.NoError(t, err)
		require.NoError(t, tarx.Pack(archive, src))
		require.NoError(t, archive.Close())

		require.NoError(t, os.WriteFile(filepath.Join(dir, "bad.tar.gz"), []byte("not a gzip file"), 0600))

		t.Chdir(dir)

		var out bytes.Buffer
		require.Error(t, tarchiveexport{Pattern: "."}.run(ctx, &out))

		var got []library.Known
		for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(&out)).Each(ctx) {
			got = append(got, v)
		}
		require.Len(t, got, 1)
		require.Equal(t, known.UID, got[0].UID)
	})
}
