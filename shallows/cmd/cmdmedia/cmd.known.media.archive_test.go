package cmdmedia

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/tarx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownarchiveRun(t *testing.T) {
	t.Run("creates an archive containing the record from stdin", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		root := t.TempDir()

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Alpha"
		known.OriginalTitle = "Alpha Original"
		known.Overview = "An overview"

		var in bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&in).Encode(known))

		require.NoError(t, knownarchive{Directory: root, Pattern: "bucket"}.run(ctx, &in))

		archdir := filepath.Join(root, "bucket")
		entries, err := os.ReadDir(archdir)
		require.NoError(t, err)

		var archives []string
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".tar.gz") {
				archives = append(archives, filepath.Join(archdir, e.Name()))
			}
		}
		require.Len(t, archives, 1)

		archive, err := os.Open(archives[0])
		require.NoError(t, err)
		defer archive.Close()

		scratch := t.TempDir()
		require.NoError(t, tarx.Unpack(scratch, archive))

		files, err := os.ReadDir(scratch)
		require.NoError(t, err)

		var got []library.Known
		for _, f := range files {
			data, err := os.ReadFile(filepath.Join(scratch, f.Name()))
			require.NoError(t, err)
			for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(bytes.NewReader(data))).Each(ctx) {
				got = append(got, v)
			}
		}

		require.Len(t, got, 1)
		require.Equal(t, known.UID, got[0].UID)
		require.Equal(t, "Alpha\nAlpha Original\nAn overview", got[0].AutoDescription)
	})

	t.Run("creates archives for multiple stdin records", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		root := t.TempDir()

		var a, b, c library.Known
		require.NoError(t, testx.Fake(&a, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&b, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&c, library.KnownOptionTestDefaults))

		var in bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&in).Encode(a))
		require.NoError(t, jsonl.NewEncoder(&in).Encode(b))
		require.NoError(t, jsonl.NewEncoder(&in).Encode(c))

		require.NoError(t, knownarchive{Directory: root, Pattern: "bucket"}.run(ctx, &in))

		archdir := filepath.Join(root, "bucket")
		entries, err := os.ReadDir(archdir)
		require.NoError(t, err)

		var got []library.Known
		for _, e := range entries {
			if !strings.HasSuffix(e.Name(), ".tar.gz") {
				continue
			}

			archive, err := os.Open(filepath.Join(archdir, e.Name()))
			require.NoError(t, err)

			scratch := t.TempDir()
			require.NoError(t, tarx.Unpack(scratch, archive))
			require.NoError(t, archive.Close())

			files, err := os.ReadDir(scratch)
			require.NoError(t, err)

			for _, f := range files {
				data, err := os.ReadFile(filepath.Join(scratch, f.Name()))
				require.NoError(t, err)
				for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(bytes.NewReader(data))).Each(ctx) {
					got = append(got, v)
				}
			}
		}

		gotUIDs := make([]string, len(got))
		for i, v := range got {
			gotUIDs[i] = v.UID
		}
		require.ElementsMatch(t, []string{a.UID, b.UID, c.UID}, gotUIDs)
	})

	t.Run("a second invocation folds new records in with previously archived ones", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		root := t.TempDir()

		var first, second library.Known
		require.NoError(t, testx.Fake(&first, library.KnownOptionTestDefaults))
		require.NoError(t, testx.Fake(&second, library.KnownOptionTestDefaults))

		var in1 bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&in1).Encode(first))
		require.NoError(t, knownarchive{Directory: root, Pattern: "bucket"}.run(ctx, &in1))

		var in2 bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&in2).Encode(second))
		require.NoError(t, knownarchive{Directory: root, Pattern: "bucket"}.run(ctx, &in2))

		archdir := filepath.Join(root, "bucket")
		entries, err := os.ReadDir(archdir)
		require.NoError(t, err)

		var got []library.Known
		for _, e := range entries {
			if !strings.HasSuffix(e.Name(), ".tar.gz") {
				continue
			}

			archive, err := os.Open(filepath.Join(archdir, e.Name()))
			require.NoError(t, err)

			scratch := t.TempDir()
			require.NoError(t, tarx.Unpack(scratch, archive))
			require.NoError(t, archive.Close())

			files, err := os.ReadDir(scratch)
			require.NoError(t, err)

			for _, f := range files {
				data, err := os.ReadFile(filepath.Join(scratch, f.Name()))
				require.NoError(t, err)
				for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(bytes.NewReader(data))).Each(ctx) {
					got = append(got, v)
				}
			}
		}

		gotUIDs := make([]string, len(got))
		for i, v := range got {
			gotUIDs[i] = v.UID
		}
		require.ElementsMatch(t, []string{first.UID, second.UID}, gotUIDs)
	})

	t.Run("DryRun does not persist the record into any archive", func(t *testing.T) {
		// The bucket directory for the record still gets created (and thus
		// still gets archived, empty) regardless of DryRun - only the
		// fsx.AppendTo write is actually skipped - so this checks that no
		// record data makes it into an archive rather than that no archive
		// file exists.
		ctx, done := testx.Context(t)
		defer done()

		root := t.TempDir()

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))

		var in bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&in).Encode(known))

		require.NoError(t, knownarchive{Directory: root, Pattern: "bucket", DryRun: true}.run(ctx, &in))

		archdir := filepath.Join(root, "bucket")
		entries, err := os.ReadDir(archdir)
		require.NoError(t, err)

		var got []library.Known
		for _, e := range entries {
			if !strings.HasSuffix(e.Name(), ".tar.gz") {
				continue
			}

			archive, err := os.Open(filepath.Join(archdir, e.Name()))
			require.NoError(t, err)

			scratch := t.TempDir()
			require.NoError(t, tarx.Unpack(scratch, archive))
			require.NoError(t, archive.Close())

			files, err := os.ReadDir(scratch)
			require.NoError(t, err)

			for _, f := range files {
				data, err := os.ReadFile(filepath.Join(scratch, f.Name()))
				require.NoError(t, err)
				for v := range jsonl.Iter[library.Known](jsonl.NewDecoder(bytes.NewReader(data))).Each(ctx) {
					got = append(got, v)
				}
			}
		}

		require.Empty(t, got)
	})

	t.Run("empty stdin produces no archives and no error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		root := t.TempDir()

		require.NoError(t, knownarchive{Directory: root, Pattern: "bucket"}.run(ctx, bytes.NewReader(nil)))

		archdir := filepath.Join(root, "bucket")
		entries, err := os.ReadDir(archdir)
		require.NoError(t, err)
		require.Empty(t, entries)
	})

	t.Run("wildcard pattern creates a fresh directory under Directory", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		root := t.TempDir()

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))

		var in bytes.Buffer
		require.NoError(t, jsonl.NewEncoder(&in).Encode(known))

		require.NoError(t, knownarchive{Directory: root, Pattern: "archive.*.d"}.run(ctx, &in))

		matches, err := filepath.Glob(filepath.Join(root, "archive.*.d"))
		require.NoError(t, err)
		require.Len(t, matches, 1)

		entries, err := os.ReadDir(matches[0])
		require.NoError(t, err)

		var archives int
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".tar.gz") {
				archives++
			}
		}
		require.Equal(t, 1, archives)
	})
}
