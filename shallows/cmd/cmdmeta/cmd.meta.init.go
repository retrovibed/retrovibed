package cmdmeta

import (
	"log"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
)

type Meta struct {
	Init MetaInit `cmd:"" help:"create and migrate the metadata and cache databases, then exit"`
}

type MetaInit struct {
	Database      string `flag:"" name:"database" help:"metadata database to initialize" default:"${vars_user_configuration_directory}/meta.db"`
	CacheDatabase string `flag:"" name:"cache-database" help:"cache database to initialize" default:"${vars_user_cache_directory}/cache.db"`
}

// Run opens both databases, which creates their parent directories and applies
// any pending migrations, and does nothing else.
func (t MetaInit) Run(gctx *cmdopts.Global) (err error) {
	db, err := cmdopts.DatabaseCustom(gctx.Context, t.Database, t.CacheDatabase)
	if err != nil {
		return errorsx.Wrap(err, "unable to initialize database")
	}

	if err = db.Close(); err != nil {
		return errorsx.Wrap(err, "unable to close database")
	}

	log.Println("initialized", t.Database, t.CacheDatabase)
	return nil
}
