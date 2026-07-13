package daemons

import (
	"archive/tar"
	"context"
	"io"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/tarx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

func MediaMetadataImport(ctx context.Context, db sqlx.Queryer, tvfs fsx.Virtual, tstore storage.ClientImpl) (err error) {
	type metadataImportBatch struct {
		metadata  tracking.Metadata
		records   []library.Known
		completed bool
	}

	log.Println("import latest media metadata initiated")
	defer log.Println("import latest media metadata completed")

	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryMediaArchive(),
			tracking.MetadataQueryCompleted(true),
			tracking.MetadataQueryNotImported(),
		},
	).OrderBy("created_at ASC")

	mdcache := torrent.NewMetadataCache(tvfs.Path())
	iter := sqlx.Scan(tracking.MetadataSearch(ctx, db, q))

	insert := asynccompute.New(func(ctx context.Context, batch metadataImportBatch) error {
		ts := time.Now()
		s := library.NewKnownBatchInsertWithDefaults(ctx, db, batch.records...)

		if err := sqlx.Discard(sqlx.Scan(s)); err != nil {
			return errorsx.Wrap(err, "failed to insert batch")
		}

		if _, err := db.ExecContext(ctx, "CHECKPOINT"); err != nil {
			return errorsx.Wrap(err, "failed to checkpoint batch")
		}

		log.Println("imported", time.Since(ts), len(batch.records), "records")

		if batch.completed {
			if err := tracking.MetadataImportedByID(ctx, db, batch.metadata.ID).Scan(&batch.metadata); err != nil {
				return errorsx.Wrap(err, "unable to mark archive as imported")
			}
		}

		return nil
	}, asynccompute.Workers[metadataImportBatch](1))

	pool := asynccompute.New(func(ctx context.Context, _md tracking.Metadata) error {
		id := int160.FromBytes(_md.Infohash)

		md, err := mdcache.Read(id)
		if err != nil {
			return errorsx.Wrap(err, "unable to read torrent metadata")
		}
		info, err := md.Metainfo().UnmarshalInfo()
		if err != nil {
			return errorsx.Wrap(err, "unable to decode torrent info")
		}

		disk, err := tstore.OpenTorrent(&info, md.ID)
		if err != nil {
			return errorsx.Wrap(err, "unable to open torrent reader")
		}

		iter, err := tarx.UnpackSeq(io.NewSectionReader(disk, 0, int64(_md.Bytes)))
		if err != nil {
			return errorsx.Wrap(err, "unable to open read archive")
		}

		importtarfile := func(header *tar.Header, content *tar.Reader) error {
			log.Println("media metadata import initiated", id, _md.Description, header.Name)
			defer log.Println("media metadata import completed", id, _md.Description, header.Name)

			d := jsonl.Iter[library.Known](jsonl.NewDecoder(content))
			for chunk := range iterx.Chunk(d.Each(ctx), 8192) {
				chunk = slicesx.Map(func(v library.Known) library.Known {
					v.AutoDescription = stringsx.Join("\n", v.Title, v.OriginalTitle, v.Overview)
					return v
				}, chunk...)

				if err := insert.Run(ctx, metadataImportBatch{metadata: _md, records: chunk}); err != nil {
					return err
				}
			}
			return nil
		}

		for header, content := range iter {
			if err := importtarfile(header, content); err != nil {
				log.Println("failed to import tarfile", header.Name)
				return err
			}
		}

		if err := insert.Run(ctx, metadataImportBatch{metadata: _md, completed: true}); err != nil {
			return errorsx.Wrap(err, "unable to mark archive as imported")
		}

		return nil
	})

	// pool must drain before insert is shut down, since pool's workers feed
	// insert; deferring unconditionally ensures both run even if this
	// function returns early (e.g. iter.Err()), so no enqueued metadata is
	// abandoned mid-import.
	defer func() {
		err = errorsx.Compact(err, asynccompute.Shutdown(ctx, pool, insert))
	}()

	for _md := range iter.Iter() {
		if err := pool.Run(ctx, _md); err != nil {
			return errorsx.Wrap(err, "unable to enqueue for import")
		}
	}

	if err := iter.Err(); err != nil {
		return errorsx.Wrap(iter.Err(), "failed to ingest latest media metadata")
	}

	return nil
}
