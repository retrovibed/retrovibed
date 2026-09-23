package cmdmedia

import (
	"database/sql"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/daemons"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type identify struct {
	Database      string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	CacheDatabase string `flag:"" name:"cache-database" help:"cache database containing known media" default:"${vars_user_cache_directory}/cache.db"`
}

func (t identify) Run(gctx *cmdopts.Global) (err error) {
	var db *sql.DB

	if db, err = cmdopts.DatabaseCustom(gctx.Context, t.Database, t.CacheDatabase); err != nil {
		return err
	}
	defer db.Close()

	return daemons.IdentifyLibraryMedia(gctx.Context, db, library.NewQueryerCleanerAuto())
}
