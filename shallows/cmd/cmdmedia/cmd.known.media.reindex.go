package cmdmedia

import (
	"context"
	"database/sql"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type knownreindex struct {
	Database      string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	CacheDatabase string `flag:"" name:"cache-database" help:"cache database containing known media" default:"${vars_user_cache_directory}/cache.db"`
	Directory     string `flag:"" name:"directory" help:"torrent storage directory containing the downloaded archives, defaults to the standard torrent directory"`
}

func (t knownreindex) Run(gctx *cmdopts.Global) (err error) {
	var db *sql.DB

	if db, err = cmdopts.DatabaseCustom(gctx.Context, t.Database, t.CacheDatabase); err != nil {
		return err
	}
	defer db.Close()

	tvfs := fsx.DirVirtual(stringsx.FirstNonBlank(t.Directory, env.TorrentDir()))

	return t.run(gctx.Context, db, tvfs, blockcache.NewTorrentFromVirtualFS(tvfs))
}

func (t knownreindex) run(ctx context.Context, db *sql.DB, tvfs fsx.Virtual, tstore storage.ClientImpl) (err error) {
	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryMediaArchive(),
			tracking.MetadataQueryCompleted(true),
		},
	)

	s := sqlx.Scan(tracking.MetadataSearch(ctx, db, q))
	for md := range s.Iter() {
		if err = tracking.MetadataImportResetByID(ctx, db, md.ID).Scan(&md); err != nil {
			return errorsx.Wrapf(err, "unable to reset archive import: %s", md.ID)
		}

		log.Println("reset archive import", md.ID, md.Description)
	}

	if err = s.Err(); err != nil {
		return errorsx.Wrap(err, "unable to scan media archives")
	}

	if _, err = db.ExecContext(ctx, "TRUNCATE cache.library_known_media"); err != nil {
		return errorsx.Wrap(err, "unable to truncate known media cache")
	}

	return daemons.MediaMetadataImport(ctx, db, tvfs, tstore)
}
