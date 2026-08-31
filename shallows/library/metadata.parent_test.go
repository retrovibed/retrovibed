package library_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestMetadataQueryParent(t *testing.T) {
	t.Run("scopes the listing to one folder", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testfolder(t, ctx, db, uuid.Nil.String(), "top")
		held := testfile(t, ctx, db, top.ID, "held")
		testfile(t, ctx, db, uuid.Nil.String(), "loose")

		found := testids(t, library.MetadataSearch(ctx, db, library.MetadataSearchBuilder().Where(library.MetadataQueryParent(top.ID))))
		require.Equal(t, []string{held.ID}, found)
	})

	t.Run("the zero uuid is the root listing", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testfolder(t, ctx, db, uuid.Nil.String(), "top")
		testfile(t, ctx, db, top.ID, "held")
		loose := testfile(t, ctx, db, uuid.Nil.String(), "loose")

		found := testids(t, library.MetadataSearch(ctx, db, library.MetadataSearchBuilder().Where(library.MetadataQueryParent(uuid.Nil.String()))))
		require.ElementsMatch(t, []string{top.ID, loose.ID}, found)
	})

	t.Run("an empty folder lists nothing", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testfolder(t, ctx, db, uuid.Nil.String(), "top")
		testfile(t, ctx, db, uuid.Nil.String(), "loose")

		require.Empty(t, testids(t, library.MetadataSearch(ctx, db, library.MetadataSearchBuilder().Where(library.MetadataQueryParent(top.ID)))))
	})
}

func TestMetadataQueryDirectory(t *testing.T) {
	t.Run("true returns only folders", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testfolder(t, ctx, db, uuid.Nil.String(), "top")
		testfile(t, ctx, db, uuid.Nil.String(), "loose")

		found := testids(t, library.MetadataSearch(ctx, db, library.MetadataSearchBuilder().Where(library.MetadataQueryDirectory(true))))
		require.Equal(t, []string{top.ID}, found)
	})

	t.Run("false returns everything but folders", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testfolder(t, ctx, db, uuid.Nil.String(), "top")
		held := testfile(t, ctx, db, top.ID, "held")
		loose := testfile(t, ctx, db, uuid.Nil.String(), "loose")

		found := testids(t, library.MetadataSearch(ctx, db, library.MetadataSearchBuilder().Where(library.MetadataQueryDirectory(false))))
		require.ElementsMatch(t, []string{held.ID, loose.ID}, found)
	})
}
