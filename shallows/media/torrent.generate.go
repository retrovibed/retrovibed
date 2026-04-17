package media

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// GenerateTorrent creates a torrent for a library item that doesn't have one.
// It creates a temp directory with a symlink to the file, generates torrent info using NewFromPath,
// then moves the content to the final location.
func GenerateTorrent(ctx context.Context, q sqlx.Queryer, mvfs, tvfs fsx.Virtual, lmd *library.Metadata) (tmd tracking.Metadata, err error) {
	writeTorrentFile := func(path string, md metainfo.MetaInfo) error {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		return md.Write(f)
	}
	var (
		info      *metainfo.Info
		infobytes []byte
	)

	src, err := blockcache.NewDirectoryCache(mvfs.Path(lmd.ID))
	if err != nil {
		return tmd, errorsx.Wrap(err, "unable to create content reader")
	}

	info, err = metainfo.NewFromReader(
		io.NewSectionReader(src, 0, int64(lmd.Bytes)),
		metainfo.OptionDisplayName(lmd.Description),
	)
	if err != nil {
		return tmd, errorsx.Wrap(err, "unable to generate torrent info")
	}

	if infobytes, err = metainfo.Encode(*info); err != nil {
		return tmd, errorsx.Wrap(err, "unable to encode torrent info")
	}

	md := metainfo.MetaInfo{
		InfoBytes: infobytes,
	}

	infohash := md.HashInfoBytes()
	tmd = tracking.NewMetadata(
		langx.Autoptr(int160.FromByteArray(infohash)),
		tracking.MetadataOptionFromInfo(info),
		tracking.MetadataOptionBytes(lmd.Bytes),
		tracking.MetadataOptionDownloaded(lmd.Bytes),
		tracking.MetadataOptionEntropySeed(infohash[:], []byte(lmd.EncryptionSeed)), // TODO: this should be passed in. we want to pull the community seed.
		tracking.MetadataOptionKnownMediaID(lmd.KnownMediaID),
		tracking.MetadataOptionAutoSeeding,
		tracking.MetadataOptionCompleted,
	)

	if err = tracking.MetadataInsertWithDefaults(ctx, q, tmd).Scan(&tmd); err != nil {
		return tmd, errorsx.Wrap(err, "unable to insert torrent metadata")
	}

	if err = writeTorrentFile(tvfs.Path(fmt.Sprintf("%s.torrent", infohash.String())), md); err != nil {
		return tmd, errorsx.Wrap(err, "unable to write torrent file")
	}

	torrentdir := tvfs.Path(infohash.String())
	dst, err := blockcache.NewDirectoryCache(torrentdir)
	if err != nil {
		return tmd, errorsx.Wrap(err, "unable to create torrent destination")
	}

	if _, err = io.Copy(io.NewOffsetWriter(dst, 0), io.NewSectionReader(src, 0, int64(lmd.Bytes))); err != nil {
		return tmd, errorsx.Wrap(err, "failed to copy media to the torrent")
	}

	// delete the original data now that we've copied it.
	// TODO: validate checksums.
	if err = os.RemoveAll(mvfs.Path(lmd.ID)); err != nil {
		return tmd, errorsx.Wrap(err, "failed to copy media to the torrent")
	}

	if err = os.Symlink(torrentdir, mvfs.Path(lmd.ID)); err != nil {
		return tmd, errorsx.Wrapf(err, "unable to create symlink to library file: %s -> %s", torrentdir, mvfs.Path(lmd.ID))
	}

	if err = library.MetadataSetTorrentID(ctx, q, lmd.ID, tmd.ID).Scan(lmd); err != nil {
		return tmd, errorsx.Wrap(err, "unable to update library metadata with torrent id")
	}

	return tmd, nil
}
