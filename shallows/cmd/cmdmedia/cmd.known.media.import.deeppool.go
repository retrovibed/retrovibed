package cmdmedia

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type deeppoolimport struct {
	StartAt time.Time `flag:"" name:"start" help:"date to start retrieving data from" format:"2006-01-02" default:"1970-01-01"`
	Source  string    `flag:"" name:"source" help:"short id for the data source" hidden:"true" default:"deeppool"`
	Limit   uint64    `flag:"" name:"limit" help:"maximum number of records to emit, 0 for unlimited" default:"131072"`
}

func (t deeppoolimport) Run(gctx *cmdopts.Global, sshid *cmdopts.SSHID) error {
	httpc, err := cmdopts.DeeppoolClientDefault{SSHID: sshid, Global: gctx}.HTTPClient(gctx.Context)
	if err != nil {
		return errorsx.Wrap(err, "unable to create deeppool client")
	}

	return t.run(gctx.Context, json.NewEncoder(os.Stdout), httpc)
}

func (t deeppoolimport) run(ctx context.Context, enc *json.Encoder, httpc *http.Client) error {
	client := communityapi.NewDeeppoolPublished(httpc)
	// UUIDv7 timestamps are an unsigned 48 bit count of milliseconds since the
	// Unix epoch; offsets before the epoch wrap around instead of clamping.
	// This is safe for now because retrovibed was created well after the Unix epoch.
	cursor := uuidx.NewV7AtTime(timex.Max(t.StartAt, time.Unix(0, 0))).String()

	seq := iterx.FromFn(func(ctx context.Context, yield func(library.Known) bool) error {
		for {
			resp, err := client.Sync(ctx, cursor)
			if err != nil {
				return errorsx.Wrap(err, "failed to sync published content from deeppool")
			}

			for _, pc := range resp.Items {
				if !yield(t.knownFromPublished(pc)) {
					return nil
				}
			}

			if len(resp.Items) == 0 {
				return nil
			}

			cursor = slicesx.LastOrZero(resp.Items...).Id
			log.Println("syncing published content", cursor)
		}
	})
	seq = iterx.Take(seq, langx.FirstNonZero(t.Limit, math.MaxUint64))

	for v := range seq.Each(ctx) {
		if err := enc.Encode(v); err != nil {
			return errorsx.Wrap(err, "unable to encode media")
		}
	}

	return seq.Err()
}

func (t deeppoolimport) knownFromPublished(pc *communityapi.PublishedContent) library.Known {
	_md5 := md5x.JSON(pc)
	uidmd5 := uuidx.FromBytesOrNil(_md5.Sum(nil))

	uid := stringsx.First(pc.KnownMediaId, ddiscapi.ImportedMediaUUID(t.Source, uuidx.FromStringOrNil(pc.Id)).String())

	return library.Known{
		Source:    t.Source,
		UID:       uid,
		Md5:       uidmd5.String(),
		Md5Lower:  binary.LittleEndian.Uint64(uuidx.LowN(uidmd5, 64)),
		ID:        pc.Id,
		Title:     pc.Title,
		Overview:  pc.Description,
		Mimetype:  stringsx.First(pc.Mimetype, mimex.Video),
		ParentUID: uuidx.Nil().String(),
	}
}
