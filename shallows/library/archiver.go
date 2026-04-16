package library

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/retrovibed/retrovibed/shallows/blockcache"
	"github.com/retrovibed/retrovibed/shallows/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

type Archiver interface {
	Upload(ctx context.Context, mimetype string, r io.Reader) (m *deeppool.Media, _ error)
}

func Archive(ctx context.Context, q sqlx.Queryer, md *Metadata, vfs fsx.Virtual, dst Archiver) (err error) {
	log.Println("archive initiated", md.ID)
	defer func() {
		log.Println("archive completed", md.ID, md.ArchiveID)
	}()
	seed := MetadataChaCha8(*md)
	cache, err := blockcache.NewDirectoryCache(vfs.Path(md.ID))
	if err != nil {
		return errorsx.Wrapf(err, "unable to create reader: %s", md.ID)
	}

	src, err := cryptox.NewReaderChaCha20(seed, io.NewSectionReader(cache, int64(md.DiskOffset), int64(md.Bytes)))
	if err != nil {
		return err
	}

	uploaded, err := dst.Upload(ctx, md.Mimetype, io.LimitReader(src, int64(md.Bytes)))
	if err != nil {
		return err
	}

	return MetadataArchivedByID(ctx, q, md.ID, uploaded.Id, uploaded.Usage).Scan(md)
}

func Reclaim(ctx context.Context, md Metadata, vfs fsx.Virtual) error {
	log.Println("reclaiming disk space initiated", md.ID)
	defer log.Println("reclaiming disk space completed", md.ID)

	dst, err := os.Readlink(vfs.Path(md.ID))
	if err != nil {
		return errorsx.Wrap(err, "failed to resolve link")
	}

	log.Println("reclaiming", vfs.Path(md.ID), "->", dst)

	return errorsx.Wrap(os.RemoveAll(dst), "failed to reclaim disk space")
}

// ReclaimEndChunks removes the last existing chunk from local storage.
// Missing chunks can be recovered via the torrent swarm or caching layer when needed.
// Returns the amount of data removed in bytes.
func ReclaimEndChunks(ctx context.Context, md Metadata, vfs fsx.Virtual) (int64, error) {
	log.Println("slow reclaiming last chunk initiated", md.ID)
	defer log.Println("slow reclaiming last chunk completed", md.ID)

	dst, err := os.Readlink(vfs.Path(md.ID))
	if err != nil {
		return 0, errorsx.Wrap(err, "failed to resolve link")
	}

	blockLength := int64(blockcache.DefaultBlockLength)
	totalChunks := int64((md.Bytes + uint64(blockLength) - 1) / uint64(blockLength))

	log.Println("file stats:", "bytes:", md.Bytes, "total_chunks:", totalChunks)

	for chunkNum := totalChunks - 1; chunkNum >= 0; chunkNum-- {
		chunkPath := filepath.Join(dst, strconv.FormatInt(chunkNum, 10))

		stat, err := os.Stat(chunkPath)
		if os.IsNotExist(err) {
			continue // chunk doesn't exist, skip
		} else if err != nil {
			return 0, errorsx.Wrap(err, "failed to stat chunk")
		}

		chunkSize := stat.Size()
		log.Println("removing last chunk", chunkPath, "size:", chunkSize)
		if err := os.Remove(chunkPath); err != nil {
			log.Println("failed to remove chunk", chunkPath, "error:", err)
			return 0, errorsx.Wrap(err, "failed to remove chunk")
		}
		return chunkSize, nil
	}

	return 0, nil
}
