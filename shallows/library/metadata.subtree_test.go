package library_test

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

// a directory row below parent. folders carry a generated id because they have no content
// to hash.
func testfolder(t *testing.T, ctx context.Context, q sqlx.Queryer, parent, description string) library.Metadata {
	var md library.Metadata
	require.NoError(t, testx.Fake(
		&md,
		library.MetadataOptionTestDefaults,
		library.MetadataOptionTestID(uuid.Must(uuid.NewV7()).String()),
		library.MetadataOptionMimetype(mimex.Directory),
		library.MetadataOptionDescription(description),
		library.MetadataOptionParentID(parent),
		library.MetadataOptionBytes(0),
	))
	require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))
	return md
}

// a content row below parent.
func testfile(t *testing.T, ctx context.Context, q sqlx.Queryer, parent, description string) library.Metadata {
	var md library.Metadata
	require.NoError(t, testx.Fake(
		&md,
		library.MetadataOptionTestDefaults,
		library.MetadataOptionTestRandomID,
		library.MetadataOptionMimetype(mimex.Binary),
		library.MetadataOptionDescription(description),
		library.MetadataOptionParentID(parent),
	))
	require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, md).Scan(&md))
	return md
}

func testids(t *testing.T, s sqlx.Scanner[library.Metadata]) []string {
	var results []library.Metadata
	require.NoError(t, sqlx.ScanInto(s, &results))
	return slicesx.MapTransform(func(md library.Metadata) string { return md.ID }, results...)
}

func TestMetadataSubtreeByID(t *testing.T) {
	t.Run("returns the folder and every descendant", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testfolder(t, ctx, db, uuid.Nil.String(), "top")
		nested := testfolder(t, ctx, db, top.ID, "nested")
		deep := testfile(t, ctx, db, nested.ID, "deep")
		outside := testfile(t, ctx, db, uuid.Nil.String(), "outside")

		found := testids(t, library.MetadataSubtreeByID(ctx, db, top.ID))
		require.ElementsMatch(t, []string{top.ID, nested.ID, deep.ID}, found)
		require.NotContains(t, found, outside.ID)
	})

	t.Run("a leaf is its own subtree", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		lone := testfile(t, ctx, db, uuid.Nil.String(), "lone")

		require.Equal(t, []string{lone.ID}, testids(t, library.MetadataSubtreeByID(ctx, db, lone.ID)))
	})

	t.Run("an unknown id returns nothing", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		require.Empty(t, testids(t, library.MetadataSubtreeByID(ctx, db, uuid.Must(uuid.NewV7()).String())))
	})
}

func TestMetadataAncestorsByID(t *testing.T) {
	t.Run("returns the row and its ancestors root first", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testfolder(t, ctx, db, uuid.Nil.String(), "top")
		nested := testfolder(t, ctx, db, top.ID, "nested")
		deep := testfile(t, ctx, db, nested.ID, "deep")

		require.Equal(t, []string{top.ID, nested.ID, deep.ID}, testids(t, library.MetadataAncestorsByID(ctx, db, deep.ID)))
	})

	t.Run("a row at the root has only itself", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testfolder(t, ctx, db, uuid.Nil.String(), "top")

		require.Equal(t, []string{top.ID}, testids(t, library.MetadataAncestorsByID(ctx, db, top.ID)))
	})
}

func TestMetadataTombstoneSubtreeByID(t *testing.T) {
	t.Run("tombstones the folder and everything below it", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testfolder(t, ctx, db, uuid.Nil.String(), "top")
		nested := testfolder(t, ctx, db, top.ID, "nested")
		deep := testfile(t, ctx, db, nested.ID, "deep")
		outside := testfile(t, ctx, db, uuid.Nil.String(), "outside")

		tombstoned := testids(t, library.MetadataTombstoneSubtreeByID(ctx, db, top.ID))
		require.ElementsMatch(t, []string{top.ID, nested.ID, deep.ID}, tombstoned)

		// nothing below the deleted folder may survive as an orphan: NewTombstonedCleanup
		// hard deletes the folder row, and a child left behind would be listed by nothing.
		surviving := testids(t, library.MetadataSearch(ctx, db, library.MetadataSearchBuilder().Where(library.MetadataQueryNotTombstoned())))
		require.Equal(t, []string{outside.ID}, surviving)
	})

	t.Run("a leaf tombstones alone", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		keep := testfile(t, ctx, db, uuid.Nil.String(), "keep")
		drop := testfile(t, ctx, db, uuid.Nil.String(), "drop")

		require.Equal(t, []string{drop.ID}, testids(t, library.MetadataTombstoneSubtreeByID(ctx, db, drop.ID)))

		surviving := testids(t, library.MetadataSearch(ctx, db, library.MetadataSearchBuilder().Where(library.MetadataQueryNotTombstoned())))
		require.Equal(t, []string{keep.ID}, surviving)
	})
}
