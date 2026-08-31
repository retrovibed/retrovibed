package library

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/shirou/gopsutil/v4/disk"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// handles clearing out tombstoned media.
func NewTombstonedCleanup(ctx context.Context, dir fsx.Virtual, q sqlx.Queryer) error {
	query := MetadataSearchBuilder().Where(squirrel.And{
		MetadataQueryTombstoned(),
	})

	log.Println("tombstoned cleanup initiated")
	defer log.Println("tombstoned cleanup completed")

	v := sqlx.Scan(MetadataSearch(ctx, q, query))
	for md := range v.Iter() {
		log.Println("cleaning", md.ID, md.Description, dir.Path(md.ID))
		errorsx.Log(errorsx.Wrapf(os.Remove(dir.Path(md.ID)), "failed to cleanup %s", dir.Path(md.ID)))
		errorsx.Log(errorsx.Wrap(MetadataDeleteByID(ctx, q, md.ID).Scan(&md), "unable to delete metadata"))
	}

	return v.Err()

}

// Moves archivable data from disk to cloud storage.
func NewAutoArchive(ctx context.Context, c *http.Client, dir fsx.Virtual, q sqlx.Queryer, async *asyncx.Wakeup, archivedisk bool) error {
	// nothing yet knows how to archive a folder: it has no bytes to upload, and the CAS
	// row it would produce is keyed on the md5 of a body it does not have.
	query := MetadataSearchBuilder().Where(squirrel.And{
		MetadataQueryArchivable(),
		MetadataQueryNotTombstoned(),
		MetadataQueryDirectory(false),
	})

	log.Println("auto archive initiated")
	defer log.Println("auto archive completed")

	archive := func(ctx context.Context) error {
		var processed uint64
		log.Println("archival initiated")
		defer log.Println("archival completed")

		a := deeppool.NewArchiver(c)

		v := sqlx.Scan(MetadataSearch(ctx, q, query))
		for md := range v.Iter() {
			processed++
			log.Println("------------------------- archivable initiated", md.ID, md.ArchiveID)

			if archivedisk {
				if err := Archive(ctx, q, &md, dir, a); err != nil {
					// if we attempt to archive a content and it fails because the content is missing that
					// means the data was deleted from disk before we archived it and there is nothing we can do to resolve it.
					// mark the data as tombstoned and move on.
					if errors.Is(err, os.ErrNotExist) {
						log.Println("media no longer exists tombstoning", md.ID, dir.Path(md.ID), err)
						errorsx.Log(errorsx.Wrap(MetadataTombstoneByID(ctx, q, md.ID).Scan(&md), "failed to tombstone media"))
					}

					errorsx.Log(errorsx.Wrapf(err, "archival upload failed: %s", md.ID))
					continue
				}
			} else {
				log.Println("dry-run - not archiving", md.ID)
			}

			log.Println("------------------------- archivable completed", md.ID)
		}

		log.Println("archival processed", processed, "records")
		return v.Err()
	}

	if err := asyncx.Run(ctx, async, archive); errorsx.Is(err, context.Canceled, context.DeadlineExceeded) {
		return nil
	} else if errorsx.Is(err, asyncx.ErrWakeupClosed) {
		// run a final time to ensure everything is archived.
		return errorsx.Ignore(archive(ctx), context.Canceled, context.DeadlineExceeded)
	} else {
		return errorsx.Wrap(err, "archival failed")
	}
}

// NewSlowDiskReclaim implements chunk-based disk reclamation that removes end chunks
// from files while preserving the beginning for immediate access. This provides a
// degraded but functional state where missing chunks can be recovered from the swarm.
func NewSlowDiskReclaim(ctx context.Context, dir fsx.Virtual, q sqlx.Queryer, async *asyncx.Wakeup, threshold float64, reclaimdisk bool) error {
	query := MetadataSearchBuilder().Where(squirrel.And{
		MetadataQueryArchived(),
		MetadataQueryHidden(false),
		MetadataQueryNotTombstoned(),
	}).OrderBy("library_metadata.updated_at ASC")

	reclaim := func(rctx context.Context) error {
		var processed uint64
		var totalReclaimed int64
		log.Println("slow disk reclaim initiated - removing one chunk per file per pass")
		defer func() {
			log.Printf("slow disk reclaim completed - processed %d files, reclaimed %d bytes", processed, totalReclaimed)
		}()

		v := sqlx.Scan(MetadataSearch(rctx, q, query))
		for md := range v.Iter() {
			if !fsx.Exists(dir.Path(md.ID)) {
				continue
			}

			if usage, err := disk.UsageWithContext(rctx, dir.Path()); err != nil {
				log.Println(errorsx.Wrap(err, "unable to retrieve disk"))
				continue
			} else if usage.UsedPercent <= threshold {
				log.Println("disk usage below threshold, stopping slow reclaim:", usage.UsedPercent)
				return v.Err()
			} else {
				log.Println("usage slow reclaiming", dir.Path(), usage.UsedPercent, usage.Fstype, md.ID, md.ArchiveID)
			}

			processed++
			log.Println("------------------------- slow reclaim initiated", md.ID, md.ArchiveID)
			if reclaimdisk {
				if bytesReclaimed, err := ReclaimEndChunks(rctx, md, dir); err != nil {
					errorsx.Log(errorsx.Wrapf(err, "slow disk reclaimation failed: %s", md.ID))
					continue
				} else {
					totalReclaimed += bytesReclaimed
				}
			} else {
				log.Println("dry-run - not reclaiming chunks for", md.ID)
			}

			log.Println("------------------------- slow reclaim completed", md.ID, md.Bytes, md.DiskUsage)
		}

		return v.Err()
	}

	if err := asyncx.Run(ctx, async, reclaim); errorsx.Is(err, context.Canceled, context.DeadlineExceeded) {
		return nil
	} else if errorsx.Is(err, asyncx.ErrWakeupClosed) {
		return errorsx.Ignore(reclaim(context.Background()), context.Canceled, context.DeadlineExceeded)
	} else {
		return errorsx.Wrap(err, "slow reclaim failed")
	}
}
