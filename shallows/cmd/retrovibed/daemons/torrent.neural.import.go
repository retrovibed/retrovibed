package daemons

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

func NeuralImport(ctx context.Context, db sqlx.Queryer, neuraldir string, tvfs fsx.Virtual, tstore storage.ClientImpl) error {
	log.Println("import latest neural initiated")
	defer log.Println("import latest neural completed")

	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryNeural(),
			tracking.MetadataQueryCompleted(true),
			tracking.MetadataQueryNotImported(),
		},
	).OrderBy("created_at ASC")

	mdcache := torrent.NewMetadataCache(tvfs.Path())
	iter := sqlx.Scan(tracking.MetadataSearch(ctx, db, q))

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

		dstpath := filepath.Join(neuraldir, _md.Description)
		dst, err := os.Create(dstpath)
		if err != nil {
			return errorsx.Wrapf(err, "unable to open neural destination: %s", dstpath)
		}
		defer dst.Close()

		log.Println("importing neural weights", dstpath)

		if _, err := io.Copy(dst, io.NewSectionReader(disk, 0, int64(_md.Bytes))); err != nil {
			return errorsx.Wrapf(err, "unable to copy neural to destination: %s - %s - %s", _md.ID, _md.Description, dstpath)
		}

		if err = tracking.MetadataImportedByID(ctx, db, _md.ID).Scan(&_md); err != nil {
			return errorsx.Wrap(err, "unable to mark neural as imported")
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

	return nil
}
