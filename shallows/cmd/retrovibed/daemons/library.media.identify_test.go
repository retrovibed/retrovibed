package daemons_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibed/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestLibraryMetadataIdentify(t *testing.T) {
	t.Run("assigns known media id when description matches", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "The Grand Budapest Hotel"
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		lmd := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Description:    "The Grand Budapest Hotel",
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Max.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		require.NoError(t, daemons.IdentifyLibraryMedia(ctx, q, library.QueryCleanerNoop()))

		require.NoError(t, library.MetadataFindByID(ctx, q, lmd.ID).Scan(&lmd))
		require.Equal(t, known.UID, lmd.KnownMediaID)
	})

	t.Run("no match marks known media id as nil", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		lmd := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Description:    "xyzzy completely unknown title 99999",
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Max.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		require.NoError(t, daemons.IdentifyLibraryMedia(ctx, q, library.QueryCleanerNoop()))

		require.NoError(t, library.MetadataFindByID(ctx, q, lmd.ID).Scan(&lmd))
		require.Equal(t, uuid.Nil.String(), lmd.KnownMediaID)
	})

	t.Run("blank cleaned description marks known media id as nil", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		lmd := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Description:    "some description",
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Max.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		require.NoError(t, daemons.IdentifyLibraryMedia(ctx, q, library.NewQueryCleanerFn(func(string) string { return "" })))

		require.NoError(t, library.MetadataFindByID(ctx, q, lmd.ID).Scan(&lmd))
		require.Equal(t, uuid.Nil.String(), lmd.KnownMediaID)
	})

	t.Run("tombstoned metadata is skipped", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		var known library.Known
		require.NoError(t, testx.Fake(&known, library.KnownOptionTestDefaults))
		known.Title = "The Grand Budapest Hotel"
		require.NoError(t, library.KnownInsertWithDefaults(ctx, q, known).Scan(&known))

		lmd := library.Metadata{
			ID:             uuid.Must(uuid.NewV7()).String(),
			Description:    "The Grand Budapest Hotel",
			Bytes:          1024,
			TorrentID:      uuid.Nil.String(),
			KnownMediaID:   uuid.Max.String(),
			ArchiveID:      uuid.Nil.String(),
			EncryptionSeed: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))
		require.NoError(t, library.MetadataTombstoneByID(ctx, q, lmd.ID).Scan(&lmd))

		require.NoError(t, daemons.IdentifyLibraryMedia(ctx, q, library.QueryCleanerNoop()))

		require.NoError(t, library.MetadataFindByID(ctx, q, lmd.ID).Scan(&lmd))
		require.Equal(t, uuid.Max.String(), lmd.KnownMediaID)
	})

	t.Run("empty metadata table returns no error", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, daemons.IdentifyLibraryMedia(ctx, q, library.QueryCleanerNoop()))
	})
}
