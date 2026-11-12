package library_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetadataUpdate(t *testing.T) {
	t.Run("should allow updating the archive_id", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		var tmp = library.Metadata{
			Description: "Example",
		}

		require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, uuidx.Nil().String(), tmp.ArchiveID)

		tmp = langx.Clone(tmp, library.MetadataOptionArchiveID(uuidx.Max().String()))
		require.NoError(t, library.MetadataUpdate(ctx, db, tmp.ID, tmp).Scan(&tmp))
		require.Equal(t, uuidx.Max().String(), tmp.ArchiveID)
	})

	t.Run("should allow updating the encryption_seed when not archived", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		var tmp = library.Metadata{
			Description: "Example",
		}

		require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		require.Equal(t, uuidx.Nil().String(), tmp.ArchiveID)

		seed := uuidx.NewV4().String()
		tmp = langx.Clone(tmp, library.MetadataOptionEncryptionSeed(seed))
		require.NoError(t, library.MetadataUpdate(ctx, db, tmp.ID, tmp).Scan(&tmp))
		require.Equal(t, seed, tmp.EncryptionSeed)
	})

	t.Run("should ignore encryption_seed updates once archived", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		var tmp = library.Metadata{
			Description: "Example",
		}

		require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults, library.MetadataOptionArchiveID(uuidx.Max().String())))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		old := tmp

		tmp = langx.Clone(tmp, library.MetadataOptionEncryptionSeed(uuidx.NewV4().String()))
		require.NoError(t, library.MetadataUpdate(ctx, db, tmp.ID, tmp).Scan(&tmp))
		assert.Equal(t, old.EncryptionSeed, tmp.EncryptionSeed)
	})

	t.Run("should bump updated_at", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)
		var tmp = library.Metadata{
			Description: "Example",
		}

		require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))
		old := tmp

		tmp = langx.Clone(tmp, library.MetadataOptionKnownMediaID(uuidx.NewV4().String()))
		require.NoError(t, library.MetadataUpdate(ctx, db, tmp.ID, tmp).Scan(&tmp))
		require.True(t, tmp.UpdatedAt.After(old.UpdatedAt))
	})
}
