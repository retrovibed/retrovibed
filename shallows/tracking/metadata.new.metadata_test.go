package tracking

import (
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/stretchr/testify/require"
)

func TestNewMetadata(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	q := sqltestx.Metadatabase(t)

	lmd := NewMetadata(
		langx.Autoptr(int160.Random()),
		MetadataOptionDescription("Hello World"),
	)

	require.NoError(t, MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))
	require.Equal(t, "Hello World", lmd.Description)
	require.Equal(t, timex.Inf(), lmd.ImportedAt)
	require.Equal(t, timex.Inf(), lmd.HiddenAt)
}
