package cmdmedia

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type deeppoolimport struct {
	StartAt time.Time `flag:"" name:"start" help:"date to start retrieving data from" format:"2006-01-02" default:"1970-01-01"`
	Source  string    `flag:"" name:"source" help:"short id for the data source" hidden:"true" default:"deeppool"`
}

func (t deeppoolimport) Run(gctx *cmdopts.Global) error {
	httpc, err := cmdopts.DeeppoolClientDefault{}.HTTPClient(gctx.Context)
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
	cursor := uuid.Must(uuid.NewV7AtTime(timex.Max(t.StartAt, time.Unix(0, 0)))).String()

	for {
		resp, err := client.Sync(ctx, cursor)
		if err != nil {
			return errorsx.Wrap(err, "failed to sync published content from deeppool")
		}

		for _, pc := range resp.Items {
			if err := enc.Encode(t.knownFromPublished(pc)); err != nil {
				return errorsx.Wrap(err, "unable to encode media")
			}
		}

		if len(resp.Items) == 0 {
			return nil
		}

		cursor = slicesx.LastOrZero(resp.Items...).Id
		log.Println("syncing published content", cursor)
	}
}

func (t deeppoolimport) knownFromPublished(pc *communityapi.PublishedContent) library.Known {
	_md5 := md5x.JSON(pc)
	uidmd5 := uuid.FromBytesOrNil(_md5.Sum(nil))

	uid := stringsx.FirstNonBlank(pc.KnownMediaId, library.KnownImportedUUID(t.Source, uuid.FromStringOrNil(pc.Id)).String())

	return library.Known{
		Source:   t.Source,
		UID:      uid,
		Md5:      uidmd5.String(),
		Md5Lower: binary.LittleEndian.Uint64(uuidx.LowN(uidmd5, 64)),
		ID:       pc.Id,
		Title:    pc.Title,
		Overview: pc.Description,
		Mimetype: stringsx.FirstNonBlank(pc.Mimetype, mimex.Video),
	}
}
