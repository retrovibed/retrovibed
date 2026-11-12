package library_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
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

	require.NoError(t, library.MetadataAssociateTorrent(ctx, db, tmp.Description, uuidx.Max().String()).Scan(&tmp))
	require.Equal(t, uuidx.Max().String(), tmp.TorrentID)
}
