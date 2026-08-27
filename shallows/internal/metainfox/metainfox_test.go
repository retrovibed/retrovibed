package metainfox_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/bytesx"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/metainfox"
	"github.com/stretchr/testify/require"
)

func TestPrinterPrint(t *testing.T) {
	t.Run("prints the details of a torrent", func(t *testing.T) {
		info, _, err := torrenttest.Random(t.TempDir(), 16*bytesx.KiB, metainfo.OptionDisplayName("example.bin"))
		require.NoError(t, err)

		md := metainfo.MetaInfo{
			InfoBytes:    testx.Must(metainfo.Encode(info))(t),
			Announce:     "https://tracker.example/announce",
			Comment:      "printer test torrent",
			CreatedBy:    "metainfox_test",
			CreationDate: time.Now().Unix(),
			Encoding:     "UTF-8",
		}

		var buf bytes.Buffer
		require.NoError(t, metainfox.NewPrinter(&md).Print(&buf))

		out := buf.String()
		require.Contains(t, out, "Name:")
		require.Contains(t, out, info.Name)
		require.Contains(t, out, "InfoHash:")
		require.Contains(t, out, md.HashInfoBytes().String())
		require.Contains(t, out, "Comment:")
		require.Contains(t, out, md.Comment)
		require.Contains(t, out, "Created By:")
		require.Contains(t, out, md.CreatedBy)
		require.Contains(t, out, "Encoding:")
		require.Contains(t, out, md.Encoding)
		require.Contains(t, out, "Private:")
		require.Contains(t, out, "false")
		require.Contains(t, out, "Piece Length:")
		require.Contains(t, out, "Pieces:")
		require.Contains(t, out, "Total Length:")
		require.Contains(t, out, "Announce List:")
		require.Contains(t, out, md.Announce)
		require.Contains(t, out, "Files:")
	})

	t.Run("returns an error when the info cannot be unmarshalled", func(t *testing.T) {
		md := metainfo.MetaInfo{InfoBytes: []byte("not valid bencode")}

		var buf bytes.Buffer
		require.Error(t, metainfox.NewPrinter(&md).Print(&buf))
	})
}
