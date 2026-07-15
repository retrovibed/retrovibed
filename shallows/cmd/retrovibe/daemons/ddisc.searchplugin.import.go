package daemons

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

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

// wasmMagic is the 4-byte binary header every valid wasm module begins with
// ("\0asm"), used to verify a discovery-search-module payload before it's
// dropped into search.d, rather than trusting the torrent's declared
// mimetype/description.
var wasmMagic = []byte{0x00, 0x61, 0x73, 0x6D}

// SearchPluginImport finds completed, unimported RetrovibedDiscoverySearch
// torrents, verifies their payload is a real wasm module via its magic
// bytes, and copies valid ones into plugindir (searchplugin's watched
// search.d directory) so the existing fsnotify watch there loads them
// automatically - no explicit Load() call is needed here. Payloads that
// aren't valid wasm are tombstoned instead of copied, so a bad torrent
// doesn't linger forever nor get silently marked imported.
func SearchPluginImport(ctx context.Context, db sqlx.Queryer, plugindir string, tvfs fsx.Virtual, tstore storage.ClientImpl) error {
	log.Println("import latest search plugins initiated")
	defer log.Println("import latest search plugins completed")

	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryDiscoverySearch(),
			tracking.MetadataQueryCompleted(true),
			tracking.MetadataQueryNotImported(),
			tracking.MetadataQueryNotTombstoned(),
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

		magic := make([]byte, len(wasmMagic))
		if _, err := io.ReadFull(io.NewSectionReader(disk, 0, int64(_md.Bytes)), magic); errorsx.Ignore(err, io.EOF, io.ErrUnexpectedEOF) != nil {
			return errorsx.Wrap(err, "unable to sniff wasm magic bytes")
		}

		if !bytes.Equal(magic, wasmMagic) {
			log.Println("rejecting search plugin, not a valid wasm module", _md.ID, _md.Description)
			if err := tracking.MetadataTombstoneByID(ctx, db, _md.ID).Scan(&_md); err != nil {
				return errorsx.Wrap(err, "unable to tombstone invalid search plugin")
			}
			return nil
		}

		dstname := _md.Description
		if !strings.HasSuffix(dstname, ".wasm") {
			dstname += ".wasm"
		}

		dstpath := filepath.Join(plugindir, dstname)
		dst, err := os.Create(dstpath)
		if err != nil {
			return errorsx.Wrapf(err, "unable to open search plugin destination: %s", dstpath)
		}
		defer dst.Close()

		log.Println("importing search plugin", dstpath)

		if _, err := io.Copy(dst, io.NewSectionReader(disk, 0, int64(_md.Bytes))); err != nil {
			return errorsx.Wrapf(err, "unable to copy search plugin to destination: %s - %s - %s", _md.ID, _md.Description, dstpath)
		}

		if err = tracking.MetadataImportedByID(ctx, db, _md.ID).Scan(&_md); err != nil {
			return errorsx.Wrap(err, "unable to mark search plugin as imported")
		}

		return nil
	})

	for _md := range iter.Iter() {
		if err := pool.Run(ctx, _md); err != nil {
			return errorsx.Wrap(err, "unable to enqueue for import")
		}
	}

	if err := iter.Err(); err != nil {
		return errorsx.Wrap(iter.Err(), "failed to ingest latest search plugins")
	}

	if err := asynccompute.Shutdown(ctx, pool); err != nil {
		return err
	}

	return nil
}
