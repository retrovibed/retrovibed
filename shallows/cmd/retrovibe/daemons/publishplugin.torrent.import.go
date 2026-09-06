package daemons

import (
	"context"
	"crypto/md5"
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
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// PublishPluginTorrentImport finds completed, unimported RetrovibedPublishModule
// torrents, verifies their payload is a real wasm module via its magic bytes,
// and installs valid ones into plugindir (publishplugin's watched publish.d
// directory) so the registry's fsnotify watch loads them automatically - no
// explicit Load() call is needed here. Payloads that aren't valid wasm are
// tombstoned instead of installed, so a bad torrent doesn't linger forever nor
// get silently marked imported.
//
// This is the community distribution path, the counterpart to cmdcommunity's
// publisher install (compile from source) and the /c/publishers upload: a
// community whose mimetype is RetrovibedPublishModule is a plugin channel, and
// subscribing to it installs what it publishes.
//
// Modules land under the content addressed {digest}.wasm name the upload
// endpoint also writes, and the catalog row is recorded here rather than left
// to PublishPluginImport, so a plugin is selectable the moment it arrives. The
// publishing community's description becomes the row's label and never reaches
// the filesystem - two communities distributing the same module converge on one
// file and one row regardless of what either of them called it.
func PublishPluginTorrentImport(ctx context.Context, db sqlx.Queryer, plugindir string, tvfs fsx.Virtual, tstore storage.ClientImpl) error {
	log.Println("import latest publish plugins initiated")
	defer log.Println("import latest publish plugins completed")

	if err := os.MkdirAll(plugindir, 0o700); err != nil {
		return errorsx.Wrapf(err, "unable to create publish plugin directory: %s", plugindir)
	}

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

		tmp, err := os.CreateTemp(plugindir, "retrovibed.import.*")
		if err != nil {
			return errorsx.Wrapf(err, "unable to create temporary file: %s", plugindir)
		}
		// a successful rename leaves nothing behind to remove; anything else
		// did, and the watcher must never see a partial module.
		defer func() {
			errorsx.Log(errorsx.Wrap(fsx.IgnoreIsNotExist(os.Remove(tmp.Name())), "unable to remove tmp"))
		}()
		defer tmp.Close()

		digest := md5.New()
		if _, err := io.Copy(io.MultiWriter(tmp, digest), io.NewSectionReader(disk, 0, int64(_md.Bytes))); err != nil {
			return errorsx.Wrapf(err, "unable to copy publish plugin: %s - %s", _md.ID, _md.Description)
		}

		publisher := md5x.FormatUUID(digest)
		dstpath := filepath.Join(plugindir, publisher+".wasm")
		if err := os.Rename(tmp.Name(), dstpath); err != nil {
			return errorsx.Wrapf(err, "unable to install publish plugin: %s", dstpath)
		}

		log.Println("importing publish plugin", dstpath)

		var inserted community.PluginPublisher
		if err := community.PluginPublisherInsertWithDefaults(ctx, db, community.PluginPublisher{
			ID:          publisher,
			Path:        dstpath,
			Description: _md.Description,
			Mimetype:    _md.Mimetype,
		}).Scan(&inserted); err != nil {
			return errorsx.Wrapf(err, "unable to record publish plugin: %s", dstpath)
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
