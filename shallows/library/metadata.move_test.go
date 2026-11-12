package library_test

import (
	"database/sql"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestMetadataMoveByID(t *testing.T) {
	t.Run("files a row under a folder", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		dst := testdirectory(t, ctx, db, uuidx.Nil().String(), "dst")
		md := testfile(t, ctx, db, uuidx.Nil().String(), "loose")

		require.NoError(t, library.MetadataMoveByID(ctx, db, md.ID, dst.ID).Scan(&md))
		require.Equal(t, dst.ID, md.DirectoryID)
	})

	t.Run("returns a row to the root", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		src := testdirectory(t, ctx, db, uuidx.Nil().String(), "src")
		md := testfile(t, ctx, db, src.ID, "filed")

		require.NoError(t, library.MetadataMoveByID(ctx, db, md.ID, uuidx.Nil().String()).Scan(&md))
		require.Equal(t, uuidx.Nil().String(), md.DirectoryID)
	})

	t.Run("moves a folder and its contents travel with it", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		dst := testdirectory(t, ctx, db, uuidx.Nil().String(), "dst")
		src := testdirectory(t, ctx, db, uuidx.Nil().String(), "src")
		held := testfile(t, ctx, db, src.ID, "held")

		require.NoError(t, library.MetadataMoveByID(ctx, db, src.ID, dst.ID).Scan(&src))
		require.Equal(t, dst.ID, src.DirectoryID)

		// the child is addressed by its parent, not by a path, so it needs no rewrite.
		require.Equal(t, []string{dst.ID, src.ID, held.ID}, testids(t, library.MetadataAncestorsByID(ctx, db, held.ID)))
	})

	t.Run("rejects a parent inside the row's own subtree", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testdirectory(t, ctx, db, uuidx.Nil().String(), "top")
		nested := testdirectory(t, ctx, db, top.ID, "nested")

		// left to run, this builds a directory_id cycle and every recursive descent then
		// spins until the process is killed.
		err := library.MetadataMoveByID(ctx, db, top.ID, nested.ID).Scan(&top)
		require.ErrorIs(t, err, sql.ErrNoRows)

		require.NoError(t, library.MetadataFindByID(ctx, db, top.ID).Scan(&top))
		require.Equal(t, uuidx.Nil().String(), top.DirectoryID)
	})

	t.Run("rejects a row into itself", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testdirectory(t, ctx, db, uuidx.Nil().String(), "top")

		require.ErrorIs(t, library.MetadataMoveByID(ctx, db, top.ID, top.ID).Scan(&top), sql.ErrNoRows)

		require.NoError(t, library.MetadataFindByID(ctx, db, top.ID).Scan(&top))
		require.Equal(t, uuidx.Nil().String(), top.DirectoryID)
	})

	t.Run("rejects a parent deep inside the row's own subtree", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		top := testdirectory(t, ctx, db, uuidx.Nil().String(), "top")
		nested := testdirectory(t, ctx, db, top.ID, "nested")
		deeper := testdirectory(t, ctx, db, nested.ID, "deeper")

		require.ErrorIs(t, library.MetadataMoveByID(ctx, db, top.ID, deeper.ID).Scan(&top), sql.ErrNoRows)
	})

	t.Run("an unknown id matches nothing", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		dst := testdirectory(t, ctx, db, uuidx.Nil().String(), "dst")

		var md library.Metadata
		require.ErrorIs(t, library.MetadataMoveByID(ctx, db, uuidx.NewV7().String(), dst.ID).Scan(&md), sql.ErrNoRows)
	})
}
