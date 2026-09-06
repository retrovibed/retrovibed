package daemons

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// PublishPluginTorrentImport finds completed, unimported RetrovibedPublishModule
// torrents, verifies their payload is a real wasm module via its magic bytes,
// and copies valid ones into plugindir (publishplugin's watched publish.d
// directory) so the registry's fsnotify watch loads them automatically - no
// explicit Load() call is needed here. Payloads that aren't valid wasm are
// tombstoned instead of copied, so a bad torrent doesn't linger forever nor get
// silently marked imported.
//
// This is the community distribution path, the counterpart to
// cmdcommunity's publisher install (compile from source) and the /c/publishers
// upload: a community whose mimetype is RetrovibedPublishModule is a plugin
// channel, and subscribing to it installs what it publishes.
//
// Loading a module is not the same as being able to use it - see
// PublishPluginImport, which is what gives it the catalog row publishing fans
// out over. Callers run that afterwards.
func PublishPluginTorrentImport(ctx context.Context, db sqlx.Queryer, plugindir string, tvfs fsx.Virtual, tstore storage.ClientImpl) error {
	log.Println("import latest publish plugins initiated")
	defer log.Println("import latest publish plugins completed")

	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryPublishModule(),
			tracking.MetadataQueryCompleted(true),
			tracking.MetadataQueryNotImported(),
			tracking.MetadataQueryNotTombstoned(),
		},
	).OrderBy("created_at ASC")

	mdcache := torrent.NewMetadataCache(tvfs.Path())
	iter := sqlx.Scan(tracking.MetadataSearch(ctx, db, q))

	pool := asynccompute.New(func(ctx context.Context, _md tracking.Metadata) error {
		id := int160.FromBytes(_md.Infohash)

		// the description is whatever the publishing community called it, so it
		// decides a filename under publish.d - sanitize before it ever reaches
		// filepath.Join. SanitizeName also matches how PublishPluginImport
		// derives an id from the filename, so reinstalling under the same name
		// reuses the id and a community's existing selection survives.
		dstname := publishplugin.SanitizeName(_md.Description)
		if dstname == "" {
			log.Println("rejecting publish plugin, unusable name", _md.ID, _md.Description)
			return errorsx.Wrap(tracking.MetadataTombstoneByID(ctx, db, _md.ID).Scan(&_md), "unable to tombstone unnamable publish plugin")
		}

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

		if err := publishplugin.VerifyWasmMagic(io.NewSectionReader(disk, 0, int64(_md.Bytes))); errors.Is(err, publishplugin.ErrNotWasi) {
			log.Println("rejecting publish plugin, not a valid wasm module", _md.ID, _md.Description)
			return errorsx.Wrap(tracking.MetadataTombstoneByID(ctx, db, _md.ID).Scan(&_md), "unable to tombstone invalid publish plugin")
		} else if err != nil {
			return errorsx.Wrap(err, "unable to sniff wasm magic bytes")
		}

		dstpath := filepath.Join(plugindir, dstname+".wasm")
		dst, err := os.Create(dstpath)
		if err != nil {
			return errorsx.Wrapf(err, "unable to open publish plugin destination: %s", dstpath)
		}
		defer dst.Close()

		log.Println("importing publish plugin", dstpath)

		if _, err := io.Copy(dst, io.NewSectionReader(disk, 0, int64(_md.Bytes))); err != nil {
			return errorsx.Wrapf(err, "unable to copy publish plugin to destination: %s - %s - %s", _md.ID, _md.Description, dstpath)
		}

		if err = tracking.MetadataImportedByID(ctx, db, _md.ID).Scan(&_md); err != nil {
			return errorsx.Wrap(err, "unable to mark publish plugin as imported")
		}

		return nil
	})

	for _md := range iter.Iter() {
		if err := pool.Run(ctx, _md); err != nil {
			return errorsx.Wrap(err, "unable to enqueue for import")
		}
	}

	if err := iter.Err(); err != nil {
		return errorsx.Wrap(iter.Err(), "failed to ingest latest publish plugins")
	}

	if err := asynccompute.Shutdown(ctx, pool); err != nil {
		return err
	}

	return nil
}
