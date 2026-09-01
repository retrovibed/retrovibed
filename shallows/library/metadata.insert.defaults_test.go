package library_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestMetadataInsertWithDefaults(t *testing.T) {
	t.Run("upsert should allow updating the archive_id from nil", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		var tmp = library.Metadata{
			Description: "Example",
		}

		require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, uuid.Nil.String(), tmp.ArchiveID)

		tmp = langx.Clone(tmp, library.MetadataOptionArchivable(true))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, uuid.Max.String(), tmp.ArchiveID)
	})

	t.Run("upsert should not overwrite a real archive_id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		var tmp = library.Metadata{
			Description: "Example",
		}

		realID := uuid.Must(uuid.NewV4()).String()
		require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults, library.MetadataOptionArchiveID(realID)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, realID, tmp.ArchiveID)

		differentID := uuid.Must(uuid.NewV4()).String()
		tmp = langx.Clone(tmp, library.MetadataOptionArchiveID(differentID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, realID, tmp.ArchiveID)
	})

	t.Run("upsert should allow updating known_media_id from nil", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		var tmp = library.Metadata{
			Description: "Example",
		}

		require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, uuid.Nil.String(), tmp.KnownMediaID)

		realID := uuid.Must(uuid.NewV4()).String()
		tmp = langx.Clone(tmp, library.MetadataOptionKnownMediaID(realID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, realID, tmp.KnownMediaID)
	})

	t.Run("upsert should allow updating known_media_id from max", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		var tmp = library.Metadata{
			Description: "Example",
		}

		require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults, library.MetadataOptionKnownMediaID(uuid.Max.String())))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, uuid.Max.String(), tmp.KnownMediaID)

		realID := uuid.Must(uuid.NewV4()).String()
		tmp = langx.Clone(tmp, library.MetadataOptionKnownMediaID(realID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, realID, tmp.KnownMediaID)
	})

	t.Run("upsert should not overwrite a real known_media_id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		var tmp = library.Metadata{
			Description: "Example",
		}

		realID := uuid.Must(uuid.NewV4()).String()
		require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults, library.MetadataOptionKnownMediaID(realID)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, realID, tmp.KnownMediaID)

		differentID := uuid.Must(uuid.NewV4()).String()
		tmp = langx.Clone(tmp, library.MetadataOptionKnownMediaID(differentID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, realID, tmp.KnownMediaID)
	})

	t.Run("upsert should not move a row into a directory", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		var tmp library.Metadata

		require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, uuid.Nil.String(), tmp.DirectoryID)

		// directory_id is absent from the conflict clause, so organization is the
		// filesystem's to change and never a side effect of re-inserting content.
		tmp = langx.Clone(tmp, library.MetadataOptionDirectoryID(uuid.Must(uuid.NewV7()).String()))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, uuid.Nil.String(), tmp.DirectoryID)
	})

	t.Run("upsert should not overwrite a directory the user already chose", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		var tmp library.Metadata

		// id is the md5 of the content, so a torrent completion or a filesystem rescan
		// re-inserts a row the user has already filed. an unguarded assignment here drags
		// the whole library back into a flat pile.
		directory := uuid.Must(uuid.NewV7()).String()
		require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults, library.MetadataOptionTestRandomID, library.MetadataOptionDirectoryID(directory)))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, directory, tmp.DirectoryID)

		tmp = langx.Clone(tmp, library.MetadataOptionDirectoryID(uuid.Nil.String()))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, directory, tmp.DirectoryID)
	})
}
