package media

import (
	"context"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// ValidateTorrent verifies the stored data matches the torrent piece hashes
// and, if valid, marks the torrent as completed and seeding in the database.
func ValidateTorrent(ctx context.Context, q sqlx.Queryer, tvfs fsx.Virtual, tmd *tracking.Metadata) error {
	infohashHex := int160.FromBytes(tmd.Infohash).String()

	mi, err := metainfo.LoadFromFile(tvfs.Path(infohashHex + tracking.TorrentSuffix))
	if err != nil {
		return errorsx.Wrap(err, "unable to load torrent file")
	}

	info, err := mi.UnmarshalInfo()
	if err != nil {
		return errorsx.Wrap(err, "unable to unmarshal torrent info")
	}

	src, err := blockcache.NewDirectoryCache(tvfs.Path(infohashHex))
	if err != nil {
		return errorsx.Wrap(err, "unable to open torrent data")
	}

	if err = mi.Verify(src); err != nil {
		return err
	}

	// downloaded stays 0: verifying already-present on-disk data involves no
	// peer transfer. available becomes the full size, since we just confirmed
	// every piece is present and hashes correctly.
	if err = tracking.MetadataCompleteByID(ctx, q, tmd.ID, 0, uint64(info.TotalLength()), 0, 0, uint64(info.TotalLength())).Scan(tmd); err != nil {
		return errorsx.Wrap(err, "unable to mark torrent as completed")
	}

	return nil
}
