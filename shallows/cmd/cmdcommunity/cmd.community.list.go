package cmdcommunity

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

type cmdCommunityList struct {
	Name string `arg:"" name:"name" help:"name of the community" required:"true"`
}

func (t cmdCommunityList) Run(gctx *cmdopts.Global, dpc cmdopts.DeeppoolClient) (err error) {
	var (
		c        *http.Client
		commresp *communityapi.CommunityFindResponse
		db       *sql.DB
	)

	if c, err = dpc.HTTPClient(gctx.Context); err != nil {
		return err
	}

	if commresp, err = communityapi.CommunityInfo(gctx.Context, c, t.Name); err != nil {
		return errorsx.Wrap(err, "failed to locate community")
	}

	if db, err = cmdopts.DatabaseMeta(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	enc := jsonl.NewEncoder(os.Stdout)
	if err = enc.Encode(commresp.Community); err != nil {
		return errorsx.Wrap(err, "failed to encode community")
	}

	q := sqlx.Scan(community.PublishedContentSearch(gctx.Context, db, community.PublishedContentSearchBuilder().Where(
		squirrel.And{
			community.PublishedContentQueryCommunityID(commresp.Community.Id),
			community.PublishedContentQueryNotTombstoned(),
		},
	).OrderBy("published_at DESC")))

	for pc := range q.Iter() {
		tmp := langx.Clone(communityapi.PublishedContent{}, communityapi.PublishedContentOptionFromDB(langx.Clone(pc, timex.JSONSafeEncodeOption)))
		tmp.CommunityId = commresp.Community.Id
		if err = enc.Encode(&tmp); err != nil {
			return errorsx.Wrap(err, "failed to encode published content")
		}
	}

	return q.Err()
}
