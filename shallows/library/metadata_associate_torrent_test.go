package library_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/stretchr/testify/require"
)

func TestMetadataAssociateTorrent(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()
	db := sqltestx.Metadatabase(t)
	var tmp = library.Metadata{
		Description: "The Misfit of Demon King Academy S01E11 1080p BD Dual Audio x265-AceAres.mkv",
	}

	require.NoError(t, testx.Fake(&tmp, library.MetadataOptionTestDefaults, library.MetadataOptionDescription("The Misfit of Demon King Academy S01E11 1080p BD Dual Audio x265-AceAres.mkv")))
	require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, tmp).Scan(&tmp))

	require.NoError(t, library.MetadataAssociateTorrent(ctx, sqlx.Debug(db), tmp.Description, uuid.Max.String()).Scan(&tmp))
	require.Equal(t, uuid.Max.String(), tmp.TorrentID)
}
