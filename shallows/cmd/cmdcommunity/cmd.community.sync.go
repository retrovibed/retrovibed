package cmdcommunity

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type cmdCommunitySync struct {
	Community    string `arg:"" name:"community" help:"community ID to sync"`
	Autodownload bool   `flag:"" name:"autodownload" help:"automatically download synced torrents" default:"false"`
}

func (t cmdCommunitySync) Run(gctx *cmdopts.Global) (err error) {
	var (
		db    *sql.DB
		httpc *http.Client
	)

	log.Println("community sync initiated", t.Community)
	defer log.Println("community sync completed", t.Community)

	if _, err = metaapi.Register(gctx.Context); err != nil {
		return errorsx.Wrap(err, "unable to register with archival service")
	}

	if httpc, err = metaapi.AutoJWTClient(gctx.Context); err != nil {
		return errorsx.Wrap(err, "unable to create api client")
	}

	if db, err = cmdmeta.Database(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	client := deeppool.NewPublished(httpc)
	synced, err := community.SyncContentFromDeeppool(gctx.Context, db, client, t.Community, t.Autodownload)
	if err != nil {
		return err
	}

	log.Println("synced", synced, "items")
	return nil
}
