package library_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
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

		res, err := identifier.Identify(ctx, "", "Inception")
		require.NoError(t, err)
		require.Equal(t, known.UID, res.UID)
		require.Greater(t, res.Relevance, 0.0)
	})

	t.Run("returns the unknown media when nothing matches", func(t *testing.T) {
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

		res, err := identifier.Identify(ctx, "", "xyzzy")
		require.NoError(t, err)
		require.Equal(t, uuid.Nil.String(), res.UID)
		require.Zero(t, res.Relevance)
	})

	t.Run("returns the unknown media when the catalog is empty", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		identifier := library.NewKnownIdentifier(db, library.QueryCleanerNoop())
		identifier.Cutoff = library.KnownMatchCutoff
		identifier.Threshold = library.KnownMatchCutoff

		res, err := identifier.Identify(ctx, "", "Inception")
		require.NoError(t, err)
		require.Equal(t, uuid.Nil.String(), res.UID)
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

		res, err := identifier.Identify(ctx, "", "Inception")
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

		res, err := identifier.Identify(ctx, "", "Inception")
		require.NoError(t, err)
		require.Equal(t, uuid.Nil.String(), res.UID)
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

		res, err := identifier.Identify(ctx, "", "Inception")
		require.NoError(t, err)
		require.Equal(t, uuid.Nil.String(), res.UID)

		identifier.Explicit = true
		res, err = identifier.Identify(ctx, "", "Inception")
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

		res, err := identifier.Identify(ctx, "", "Inception")
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
		res, err := identifier.Identify(ctx, "", "Inception 1080p x264 GROUP")
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

		res, err := identifier.Identify(ctx, "", "Inception S01E02")
		require.NoError(t, err)
		require.Equal(t, expected.UID, res.UID)
	})

	t.Run("defaults the configuration", func(t *testing.T) {
		identifier := library.NewKnownIdentifier(nil, library.QueryCleanerNoop())

		require.Equal(t, float32(0.7), identifier.Cutoff)
		require.Equal(t, float32(0.7), identifier.Threshold)
		require.Equal(t, 0.85, identifier.MinRelevance)
		require.Equal(t, uint(8), identifier.Limit)
		require.False(t, identifier.Explicit)
	})

	t.Run("options override the defaults", func(t *testing.T) {
		identifier := library.NewKnownIdentifier(nil, library.QueryCleanerNoop(), func(t *library.KnownIdentifier) {
			t.Limit = 2
		})

		require.Equal(t, uint(2), identifier.Limit)
		require.Equal(t, float32(0.7), identifier.Cutoff)
	})

	t.Run("restricts candidates to the mimetype unless blank", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults, library.KnownOptionMimetype(mimex.Application)))
		known.Title = "Inception"
		library.KnownOptionAutoDescription(&known)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, known).Scan(&known))

		identifier := library.NewKnownIdentifier(db, library.QueryCleanerNoop())

		res, err := identifier.Identify(ctx, mimex.Application, "Inception")
		require.NoError(t, err)
		require.Equal(t, known.UID, res.UID)

		res, err = identifier.Identify(ctx, mimex.Audio, "Inception")
		require.NoError(t, err)
		require.Equal(t, uuid.Nil.String(), res.UID)

		res, err = identifier.Identify(ctx, "", "Inception")
		require.NoError(t, err)
		require.Equal(t, known.UID, res.UID)
	})

	t.Run("only scores candidates within the limit", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		// decoys share the standalone collation the search targets, so they are ordered ahead of the
		// best match despite being far weaker matches.
		for range 8 {
			var decoy library.Known
			require.NoError(t, testx.Fake(&decoy, library.KnownOptionTestDefaults, library.KnownOptionCollation(0)))
			decoy.Title = "Zzzzz"
			decoy.OriginalTitle = "Zzzzz"
			decoy.Overview = "a film that mentions Inception in passing"
			library.KnownOptionAutoDescription(&decoy)
			require.NoError(t, library.KnownInsertWithDefaults(ctx, db, decoy).Scan(&decoy))
		}

		var best library.Known
		require.NoError(t, testx.Fake(&best, library.KnownOptionTestDefaults, library.KnownOptionCollation(library.KnownCollationEpisode(1, 1))))
		best.Title = "Inception"
		library.KnownOptionAutoDescription(&best)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, best).Scan(&best))

		identifier := library.NewKnownIdentifier(db, library.QueryCleanerNoop())

		res, err := identifier.Identify(ctx, "", "Inception")
		require.NoError(t, err)
		require.Equal(t, uuid.Nil.String(), res.UID)

		identifier.Limit = 9
		res, err = identifier.Identify(ctx, "", "Inception")
		require.NoError(t, err)
		require.Equal(t, best.UID, res.UID)
	})

	t.Run("does not search when the cleaner blanks the input", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		cleaner := library.NewQueryCleanerFn(func(string) string { return "" })
		// the nil queryer fails any attempt to query the catalog.
		identifier := library.NewKnownIdentifier(nil, cleaner)

		res, err := identifier.Identify(ctx, "", "Inception")
		require.NoError(t, err)
		require.Equal(t, uuid.Nil.String(), res.UID)
	})

	t.Run("identifies a candidate ranked behind weaker candidates within the limit", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var decoy library.Known
		require.NoError(t, testx.Fake(&decoy, library.KnownOptionTestDefaults, library.KnownOptionCollation(0)))
		decoy.Title = "Zzzzz"
		decoy.OriginalTitle = "Zzzzz"
		decoy.Overview = "a film that mentions Inception in passing"
		library.KnownOptionAutoDescription(&decoy)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, decoy).Scan(&decoy))

		var best library.Known
		require.NoError(t, testx.Fake(&best, library.KnownOptionTestDefaults, library.KnownOptionCollation(library.KnownCollationEpisode(1, 1))))
		best.Title = "Inception"
		library.KnownOptionAutoDescription(&best)
		require.NoError(t, library.KnownInsertWithDefaults(ctx, db, best).Scan(&best))

		identifier := library.NewKnownIdentifier(db, library.QueryCleanerNoop())

		res, err := identifier.Identify(ctx, "", "Inception")
		require.NoError(t, err)
		require.Equal(t, best.UID, res.UID)
	})
}
