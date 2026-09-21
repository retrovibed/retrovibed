package library_test

import (
	"context"
	"errors"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestKnownIdentifierIdentify(t *testing.T) {
	t.Run("identifies the matching entry with a positive relevance", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Inception"
		library.KnownOptionAutoDescription(&known)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		identifier := library.NewKnownIdentifier(db, library.QueryCleanerNoop())
		identifier.Cutoff = library.KnownMatchCutoff
		identifier.Threshold = library.KnownMatchCutoff

		res, err := identifier.Identify(ctx, "Inception")
		require.NoError(t, err)
		require.Equal(t, known.UID, res.UID)
		require.Greater(t, res.Relevance, 0.0)
	})

	t.Run("returns the zero value when nothing matches", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Inception"
		library.KnownOptionAutoDescription(&known)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		identifier := library.NewKnownIdentifier(db, library.QueryCleanerNoop())
		identifier.Cutoff = library.KnownMatchCutoff
		identifier.Threshold = library.KnownMatchCutoff

		res, err := identifier.Identify(ctx, "xyzzy")
		require.NoError(t, err)
		require.Empty(t, res.UID)
		require.Zero(t, res.Relevance)
	})

	t.Run("returns the zero value when the catalog is empty", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		identifier := library.NewKnownIdentifier(db, library.QueryCleanerNoop())
		identifier.Cutoff = library.KnownMatchCutoff
		identifier.Threshold = library.KnownMatchCutoff

		res, err := identifier.Identify(ctx, "Inception")
		require.NoError(t, err)
		require.Empty(t, res.UID)
	})

	t.Run("prefers the candidate with the highest relevance", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var exact library.Known
		require.NoError(t, testx.Fake(&exact, library.KnownOptionTestDefaults))
		exact.Title = "Inception"
		library.KnownOptionAutoDescription(&exact)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, exact).Scan(&exact))

		var partial library.Known
		require.NoError(t, testx.Fake(&partial, library.KnownOptionTestDefaults))
		partial.Title = "Inception Documentary Extras"
		library.KnownOptionAutoDescription(&partial)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, partial).Scan(&partial))

		identifier := library.NewKnownIdentifier(db, library.QueryCleanerNoop())
		identifier.Cutoff = library.KnownMatchCutoff
		identifier.Threshold = library.KnownMatchCutoff

		res, err := identifier.Identify(ctx, "Inception")
		require.NoError(t, err)
		require.Equal(t, exact.UID, res.UID)
	})

	t.Run("weak title matches below the cutoff are not identified", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Zzzzz"
		known.OriginalTitle = "Zzzzz"
		known.Overview = "a film that mentions Inception in passing"
		library.KnownOptionAutoDescription(&known)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		identifier := library.NewKnownIdentifier(db, library.QueryCleanerNoop())
		identifier.Cutoff = library.KnownMatchCutoff
		identifier.Threshold = library.KnownMatchCutoff

		res, err := identifier.Identify(ctx, "Inception")
		require.NoError(t, err)
		require.Empty(t, res.UID)
	})

	t.Run("excludes adult content unless explicit", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Inception"
		known.Adult = true
		library.KnownOptionAutoDescription(&known)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		identifier := library.NewKnownIdentifier(db, library.QueryCleanerNoop())
		identifier.Cutoff = library.KnownMatchCutoff
		identifier.Threshold = library.KnownMatchCutoff

		res, err := identifier.Identify(ctx, "Inception")
		require.NoError(t, err)
		require.Empty(t, res.UID)

		identifier.Explicit = true
		res, err = identifier.Identify(ctx, "Inception")
		require.NoError(t, err)
		require.Equal(t, known.UID, res.UID)
	})

	t.Run("falls back to the raw input when the cleaner fails", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Inception"
		library.KnownOptionAutoDescription(&known)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		cleaner := library.QueryCleanerFn(func(context.Context, string) (string, error) {
			return "", errors.New("cleaner unavailable")
		})
		identifier := library.NewKnownIdentifier(db, cleaner)
		identifier.Cutoff = library.KnownMatchCutoff
		identifier.Threshold = library.KnownMatchCutoff

		res, err := identifier.Identify(ctx, "Inception")
		require.NoError(t, err)
		require.Equal(t, known.UID, res.UID)
	})

	t.Run("searches with the cleaned query rather than the raw input", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "Inception"
		library.KnownOptionAutoDescription(&known)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		cleaner := library.NewQueryCleanerFn(func(string) string { return "Inception" })
		identifier := library.NewKnownIdentifier(db, cleaner)
		identifier.Cutoff = library.KnownMatchCutoff
		identifier.Threshold = library.KnownMatchCutoff

		// the raw input carries release noise that would never match the catalog on its own.
		res, err := identifier.Identify(ctx, "Inception 1080p x264 GROUP")
		require.NoError(t, err)
		require.Equal(t, known.UID, res.UID)
	})

	t.Run("breaks relevance ties using the requested episode", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var wrong library.Known
		require.NoError(t, testx.Fake(&wrong, library.KnownOptionTestDefaults, library.KnownOptionCollation(library.KnownCollationEpisode(1, 3))))
		wrong.Title = "Inception"
		library.KnownOptionAutoDescription(&wrong)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, wrong).Scan(&wrong))

		var expected library.Known
		require.NoError(t, testx.Fake(&expected, library.KnownOptionTestDefaults, library.KnownOptionCollation(library.KnownCollationEpisode(1, 2))))
		expected.Title = "Inception"
		library.KnownOptionAutoDescription(&expected)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, expected).Scan(&expected))

		cleaner := library.NewQueryCleanerFn(func(string) string { return "Inception\nS01E02" })
		identifier := library.NewKnownIdentifier(db, cleaner)
		identifier.Cutoff = library.KnownMatchCutoff
		identifier.Threshold = library.KnownMatchCutoff

		res, err := identifier.Identify(ctx, "Inception S01E02")
		require.NoError(t, err)
		require.Equal(t, expected.UID, res.UID)
	})
}
