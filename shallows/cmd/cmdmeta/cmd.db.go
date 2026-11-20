package cmdmeta

import (
	"context"
	"database/sql"

	"github.com/egdaemon/duckproxy/duckkong"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
)

// DB opens a sql shell, with psql style meta commands (i.e. \d), on the main
// database: meta.db with cache.db attached as the cache catalog. the database is
// opened through duckshare, so the shell works while the daemon is running.
type DB struct {
	Command string `name:"command" short:"c" help:"run the statements, or meta command, and exit instead of starting the shell"`
}

func (t DB) Run(gctx *cmdopts.Global) error {
	return duckkong.Cmd{Command: t.Command, Prompt: "retrovibe%R%# "}.Run(gctx.Context, func(ctx context.Context, _ string) (*sql.DB, error) {
		return cmdopts.Database(ctx)
	})
}
