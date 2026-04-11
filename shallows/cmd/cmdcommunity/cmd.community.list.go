package cmdcommunity

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/community"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/jsonl"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/metaapi"
)

type cmdCommunityList struct {
	Name string `arg:"" name:"name" help:"name of the community" required:"true"`
}

func (t cmdCommunityList) Run(gctx *cmdopts.Global, dpc cmdopts.DeeppoolClient) (err error) {
	var (
		c        *http.Client
		commresp *meta.CommunityFindResponse
		db       *sql.DB
	)

	if c, err = dpc.HTTPClient(gctx.Context); err != nil {
		return err
	}

	if commresp, err = metaapi.CommunityInfo(gctx.Context, c, t.Name); err != nil {
		return errorsx.Wrap(err, "failed to locate community")
	}

	if db, err = cmdmeta.Database(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	enc := jsonl.NewEncoder(os.Stdout)
	if err = enc.Encode(commresp.Community); err != nil {
		return errorsx.Wrap(err, "failed to encode community")
	}

	q := sqlx.Scan(community.PublishedContentFindByCommunityID(gctx.Context, db, commresp.Community.Id))
	for pc := range q.Iter() {
		tmp := langx.Clone(meta.PublishedContent{}, community.PublishedContentOptionFromDB(langx.Clone(pc, timex.JSONSafeEncodeOption)))
		tmp.CommunityId = commresp.Community.Id
		if err = enc.Encode(&tmp); err != nil {
			return errorsx.Wrap(err, "failed to encode published content")
		}
	}

	return q.Err()
}
