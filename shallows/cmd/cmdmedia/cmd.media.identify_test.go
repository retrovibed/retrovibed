package cmdmedia

import (
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdtestx"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestIdentifyRun(t *testing.T) {
	// isolate from any locally cached neural model so the auto cleaner falls back to noop.
	t.Setenv(env.MediaIdentificationDisabled, "true")

	t.Run("assigns known media id when description matches", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		metapath := filepath.Join(t.TempDir(), "meta.db")
		cachepath := filepath.Join(t.TempDir(), "cache.db")

		db, err := cmdopts.DatabaseCustom(ctx, metapath, cachepath)
		require.NoError(t, err)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "The Grand Budapest Hotel"
		library.KnownOptionAutoDescription(&known)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		var lmd library.Metadata
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("The Grand Budapest Hotel"), library.MetadataOptionKnownMediaID(uuid.Max.String())))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, lmd).Scan(&lmd))
		require.NoError(t, db.Close())

		genparser := cmdtestx.Genparser(identify{})
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "--database", metapath, "--cache-database", cachepath))

		db, err = cmdopts.DatabaseCustom(ctx, metapath, cachepath)
		require.NoError(t, err)
		defer db.Close()

		require.NoError(t, library.MetadataFindByID(ctx, db, lmd.ID).Scan(&lmd))
		require.Equal(t, known.UID, lmd.KnownMediaID)
	})

	t.Run("no match marks known media id as nil", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		metapath := filepath.Join(t.TempDir(), "meta.db")
		cachepath := filepath.Join(t.TempDir(), "cache.db")

		db, err := cmdopts.DatabaseCustom(ctx, metapath, cachepath)
		require.NoError(t, err)

		var lmd library.Metadata
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("xyzzy completely unknown title 99999"), library.MetadataOptionKnownMediaID(uuid.Max.String())))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, lmd).Scan(&lmd))
		require.NoError(t, db.Close())

		genparser := cmdtestx.Genparser(identify{})
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "--database", metapath, "--cache-database", cachepath))

		db, err = cmdopts.DatabaseCustom(ctx, metapath, cachepath)
		require.NoError(t, err)
		defer db.Close()

		require.NoError(t, library.MetadataFindByID(ctx, db, lmd.ID).Scan(&lmd))
		require.Equal(t, uuid.Nil.String(), lmd.KnownMediaID)
	})

	t.Run("already identified records are left unchanged", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		metapath := filepath.Join(t.TempDir(), "meta.db")
		cachepath := filepath.Join(t.TempDir(), "cache.db")

		db, err := cmdopts.DatabaseCustom(ctx, metapath, cachepath)
		require.NoError(t, err)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "The Grand Budapest Hotel"
		library.KnownOptionAutoDescription(&known)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		kid := uuid.Must(uuid.NewV4()).String()
		var lmd library.Metadata
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDescription("The Grand Budapest Hotel"), library.MetadataOptionKnownMediaID(kid)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, lmd).Scan(&lmd))
		require.NoError(t, db.Close())

		genparser := cmdtestx.Genparser(identify{})
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "--database", metapath, "--cache-database", cachepath))

		db, err = cmdopts.DatabaseCustom(ctx, metapath, cachepath)
		require.NoError(t, err)
		defer db.Close()

		require.NoError(t, library.MetadataFindByID(ctx, db, lmd.ID).Scan(&lmd))
		require.Equal(t, kid, lmd.KnownMediaID)
	})

	t.Run("empty database returns no error", func(t *testing.T) {
		metapath := filepath.Join(t.TempDir(), "meta.db")
		cachepath := filepath.Join(t.TempDir(), "cache.db")

		genparser := cmdtestx.Genparser(identify{})
		require.NoError(t, cmdtestx.Execute(t, genparser(t), "command", "--database", metapath, "--cache-database", cachepath))
	})
}
