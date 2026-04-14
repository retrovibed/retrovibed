package media_test

import (
	"testing"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/torrenttestx"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestDownloadOptionFromTorrent(t *testing.T) {
	tclient := torrenttestx.QuickClient(t)

	t.Run("single file torrent reports correct bytes", func(t *testing.T) {
		info := &metainfo.Info{
			PieceLength: 256 * 1024,
			Length:      2 * 1024 * 1024,
			Pieces:      make([]byte, 20*8),
		}
		tmd, err := torrent.NewFromInfo(info)
		require.NoError(t, err)

		dl, _, err := tclient.Start(tmd)
		require.NoError(t, err)

		result := langx.Clone(
			media.Download{},
			media.DownloadOptionFromTorrentMetadata(tracking.NewMetadata(langx.Autoptr(tmd.ID))),
			media.DownloadOptionFromTorrent(dl),
		)

		require.EqualValues(t, 2*1024*1024, result.Bytes)
	})

	t.Run("multi file torrent reports correct bytes", func(t *testing.T) {
		info := &metainfo.Info{
			PieceLength: 256 * 1024,
			Pieces:      make([]byte, 20*8),
			Files: []metainfo.FileInfo{
				{Length: 1024 * 1024, Path: []string{"a.txt"}},
				{Length: 1024 * 1024, Path: []string{"b.txt"}},
			},
		}
		tmd, err := torrent.NewFromInfo(info)
		require.NoError(t, err)

		dl, _, err := tclient.Start(tmd)
		require.NoError(t, err)

		result := langx.Clone(
			media.Download{},
			media.DownloadOptionFromTorrentMetadata(tracking.NewMetadata(langx.Autoptr(tmd.ID))),
			media.DownloadOptionFromTorrent(dl),
		)

		require.EqualValues(t, 2*1024*1024, result.Bytes, "multi-file torrent bytes should be sum of all file lengths")
	})
}
