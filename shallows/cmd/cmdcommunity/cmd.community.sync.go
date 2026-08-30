package cmdcommunity

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/communityapi"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type cmdCommunitySync struct {
	Community    string `arg:"" name:"community" help:"community ID to sync"`
	Autodownload bool   `flag:"" name:"autodownload" help:"automatically download synced torrents" default:"false"`
}

func (t cmdCommunitySync) Run(gctx *cmdopts.Global, sshid *cmdopts.SSHID) (err error) {
	var (
		db    *sql.DB
		httpc *http.Client
	)

	log.Println("community sync initiated", t.Community)
	defer log.Println("community sync completed", t.Community)

	id, err := sshid.Signer()
	if err != nil {
		return errorsx.Wrap(err, "unable to generate signer id")
	}

	if _, err = authn.AutoRegistration(gctx.Context, id); err != nil {
		return errorsx.Wrap(err, "unable to register with archival service")
	}

	if httpc, err = authn.AutoJWTClient(gctx.Context, id); err != nil {
		return errorsx.Wrap(err, "unable to create api client")
	}

	if db, err = cmdopts.DatabaseMeta(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	client := communityapi.NewDeeppoolPublished(httpc)
	synced, err := communityapi.SyncContentFromDeeppool(gctx.Context, db, client, t.Community, t.Autodownload, 0)
	if err != nil {
		return err
	}

	log.Println("synced", synced, "items")
	return nil
}
