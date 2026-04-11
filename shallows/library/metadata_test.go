package library_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/stretchr/testify/require"
)

func TestMetadataQueries(t *testing.T) {
	t.Run("upsert should allow updating the archive_id", func(t *testing.T) {
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

	t.Run("MetadataUpdate should allow updating the archive_id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		var tmp = library.Metadata{
			Description: "Example",
		}

		require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, uuid.Nil.String(), tmp.ArchiveID)

		tmp = langx.Clone(tmp, library.MetadataOptionArchiveID(uuid.Max.String()))
		require.NoError(t, library.MetadataUpdate(ctx, db, tmp.ID, tmp).Scan(&tmp))
		require.Equal(t, uuid.Max.String(), tmp.ArchiveID)
	})
}
