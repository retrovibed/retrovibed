package daemons

import (
	"archive/tar"
	"context"
	"io"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/internal/asynccompute"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/jsonl"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/stringsx"
	"github.com/retrovibed/retrovibed/internal/tarx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/retrovibed/retrovibed/tracking"
)

func MediaMetadataImport(ctx context.Context, db sqlx.Queryer, tvfs fsx.Virtual, tstore storage.ClientImpl) error {
	log.Println("import latest media metadata initiated")
	defer log.Println("import latest media metadata completed")

	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryMetadataArchive(),
			tracking.MetadataQueryCompleted(true),
			tracking.MetadataQueryNotImported(),
		},
	).OrderBy("created_at ASC")

	mdcache := torrent.NewMetadataCache(tvfs.Path())
	iter := sqlx.Scan(tracking.MetadataSearch(ctx, db, q))

	insert := asynccompute.New(func(ctx context.Context, v library.Known) error {
		return library.KnownInsertWithDefaults(ctx, db, v).Scan(&v)
	}, asynccompute.Workers[library.Known](1))

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
			var (
				derr error
				i    uint64
				v    library.Known
			)

			log.Println("media metadata import initiated", id, _md.Description, header.Name)
			defer log.Println("media metadata import completed", id, _md.Description, header.Name)
			d := jsonl.NewDecoder(content)

			for derr = d.Decode(&v); derr == nil; i, derr = i+1, d.Decode(&v) {
				v.AutoDescription = stringsx.Join("\n", v.Title, v.OriginalTitle, v.Overview)
				if err := insert.Run(ctx, v); err != nil {
					return err
				}
			}

			if err := errorsx.Ignore(derr, io.EOF); err != nil {
				return err
			}

			return nil
		}

		for header, content := range iter {
			if err := importtarfile(header, content); err != nil {
				return err
			}
		}

		if err = tracking.MetadataImportedByID(ctx, db, _md.ID).Scan(&_md); err != nil {
			return errorsx.Wrap(err, "unable to mark archive as imported")
		}

		return nil
	})

	for _md := range iter.Iter() {
		if err := pool.Run(ctx, _md); err != nil {
			return errorsx.Wrap(err, "unable to enqueue for import")
		}
	}

	if err := iter.Err(); err != nil {
		return errorsx.Wrap(iter.Err(), "failed to ingest latest media metadata")
	}

	if err := asynccompute.Shutdown(ctx, pool); err != nil {
		return err
	}

	return asynccompute.Shutdown(ctx, insert)
}
