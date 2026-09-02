package library_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestMetadataQueryDirectoryID(t *testing.T) {
	t.Run("scopes the listing to one folder", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testdirectory(t, ctx, db, uuid.Nil.String(), "top")
		held := testfile(t, ctx, db, top.ID, "held")
		testfile(t, ctx, db, uuid.Nil.String(), "loose")

		found := testids(t, library.MetadataSearch(ctx, db, library.MetadataSearchBuilder().Where(library.MetadataQueryDirectoryID(top.ID))))
		require.Equal(t, []string{held.ID}, found)
	})

	t.Run("the zero uuid is the root listing", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testdirectory(t, ctx, db, uuid.Nil.String(), "top")
		testfile(t, ctx, db, top.ID, "held")
		loose := testfile(t, ctx, db, uuid.Nil.String(), "loose")

		found := testids(t, library.MetadataSearch(ctx, db, library.MetadataSearchBuilder().Where(library.MetadataQueryDirectoryID(uuid.Nil.String()))))
		require.ElementsMatch(t, []string{top.ID, loose.ID}, found)
	})

	t.Run("an empty folder lists nothing", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testdirectory(t, ctx, db, uuid.Nil.String(), "top")
		testfile(t, ctx, db, uuid.Nil.String(), "loose")

		require.Empty(t, testids(t, library.MetadataSearch(ctx, db, library.MetadataSearchBuilder().Where(library.MetadataQueryDirectoryID(top.ID)))))
	})
}

func TestMetadataQueryIsDirectory(t *testing.T) {
	t.Run("true returns only folders", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testdirectory(t, ctx, db, uuid.Nil.String(), "top")
		testfile(t, ctx, db, uuid.Nil.String(), "loose")

		found := testids(t, library.MetadataSearch(ctx, db, library.MetadataSearchBuilder().Where(library.MetadataQueryIsDirectory(true))))
		require.Equal(t, []string{top.ID}, found)
	})

	t.Run("false returns everything but folders", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testdirectory(t, ctx, db, uuid.Nil.String(), "top")
		held := testfile(t, ctx, db, top.ID, "held")
		loose := testfile(t, ctx, db, uuid.Nil.String(), "loose")

		found := testids(t, library.MetadataSearch(ctx, db, library.MetadataSearchBuilder().Where(library.MetadataQueryIsDirectory(false))))
		require.ElementsMatch(t, []string{held.ID, loose.ID}, found)
	})
}
