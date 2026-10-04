package cmdmedia

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/retroapi/iterx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/envx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

// golden copy of the generated media metadata stored in deeppool.
type knowngolden struct {
	Env    goldenenv    `cmd:"" help:"extract environment information for known media import from the deeppool golden copy"`
	Upload goldenupload `cmd:"" help:"processes stdin known media jsonl and uploads it to the deeppool golden copy"`
	Sync   goldensync   `cmd:"" help:"sync known media from the deeppool golden copy and writes to stdout in known media jsonl format for importing"`
}

type goldenenv struct {
	ExcludeSource []string `flag:"" name:"exclude-source" help:"exclude the specified source(s) from the calculation" optional:""`
	IncludeSource []string `flag:"" name:"include-source" help:"limit the calculation to the specified source(s)" optional:""`
}

func (t goldenenv) Run(gctx *cmdopts.Global, sshid *cmdopts.SSHID) error {
	httpc, err := cmdopts.DeeppoolClientDefault{SSHID: sshid}.HTTPClient(gctx.Context)
	if err != nil {
		return errorsx.Wrap(err, "unable to create deeppool client")
	}

	return t.run(gctx.Context, os.Stdout, httpc)
}

func (t goldenenv) run(ctx context.Context, w io.Writer, httpc *http.Client) error {
	resp, err := deeppool.NewMediaMetadata(httpc).Env(ctx, t.IncludeSource, t.ExcludeSource)
	if err != nil {
		return errorsx.Wrap(err, "unable to retrieve environment from deeppool")
	}

	return envx.Build().Var(
		"RETROVIBED_ARCHIVE_START_DATE",
		resp.Start,
	).CopyTo(w)
}

type goldenupload struct {
	Batch int `flag:"" name:"batch" help:"number of records to upload per request" default:"1024"`
}

func (t goldenupload) Run(gctx *cmdopts.Global, sshid *cmdopts.SSHID, stdin cmdopts.Stdin) error {
	httpc, err := cmdopts.DeeppoolClientDefault{SSHID: sshid}.HTTPClient(gctx.Context)
	if err != nil {
		return errorsx.Wrap(err, "unable to create deeppool client")
	}

	return t.run(gctx.Context, langx.FirstNonNil[io.Reader](stdin, os.Stdin), httpc)
}

func (t goldenupload) run(ctx context.Context, r io.Reader, httpc *http.Client) error {
	client := deeppool.NewMediaMetadata(httpc)
	d := jsonl.Iter[library.Known](jsonl.NewDecoder(r))

	for chunk := range iterx.Chunk(d.Each(ctx), t.Batch) {
		items := make([]*deeppool.MediaMetadata, 0, len(chunk))
		for _, k := range chunk {
			items = append(items, &deeppool.MediaMetadata{
				Uid:              k.UID,
				Md5:              k.Md5,
				Source:           k.Source,
				Id:               k.ID,
				TombstonedAt:     timex.RFC3339NanoEncode(k.TombstonedAt).Format(time.RFC3339Nano),
				Released:         timex.RFC3339NanoEncode(k.Released).Format(time.RFC3339Nano),
				Title:            k.Title,
				Popularity:       k.Popularity,
				Adult:            k.Adult,
				BackdropPath:     k.BackdropPath,
				PosterPath:       k.PosterPath,
				OriginalLanguage: k.OriginalLanguage,
				OriginalTitle:    k.OriginalTitle,
				Overview:         k.Overview,
				Mimetype:         k.Mimetype,
				Collation:        k.Collation,
				Subtitle:         k.Subtitle,
				ParentUid:        k.ParentUID,
			})
		}

		resp, err := client.Upload(ctx, items...)
		if err != nil {
			return errorsx.Wrap(err, "unable to upload media metadata")
		}

		log.Println("uploaded", resp.Received, "received", resp.Upserted, "upserted")
	}

	return d.Err()
}

type goldensync struct {
	Sync string `flag:"" name:"sync" help:"uuidv7 sync id to resume from" default:"00000000-0000-0000-0000-000000000000"`
}

func (t goldensync) Run(gctx *cmdopts.Global, sshid *cmdopts.SSHID) error {
	httpc, err := cmdopts.DeeppoolClientDefault{SSHID: sshid}.HTTPClient(gctx.Context)
	if err != nil {
		return errorsx.Wrap(err, "unable to create deeppool client")
	}

	return t.run(gctx.Context, json.NewEncoder(os.Stdout), httpc)
}

func (t goldensync) run(ctx context.Context, enc *json.Encoder, httpc *http.Client) error {
	client := deeppool.NewMediaMetadata(httpc)
	cursor := t.Sync

	for {
		resp, err := client.Sync(ctx, cursor)
		if err != nil {
			return errorsx.Wrap(err, "failed to sync media metadata from deeppool")
		}

		if len(resp.Items) == 0 {
			return nil
		}

		for _, m := range resp.Items {
			released, err := time.Parse(time.RFC3339Nano, m.Released)
			if err != nil {
				return errorsx.Wrapf(err, "invalid released timestamp: %s", m.Uid)
			}

			tombstoned, err := time.Parse(time.RFC3339Nano, m.TombstonedAt)
			if err != nil {
				return errorsx.Wrapf(err, "invalid tombstoned timestamp: %s", m.Uid)
			}

			md5 := uuid.FromStringOrNil(m.Md5)

			if err := enc.Encode(library.Known{
				UID:              m.Uid,
				Md5:              md5.String(),
				Md5Lower:         binary.LittleEndian.Uint64(uuidx.LowN(md5, 64)),
				Source:           m.Source,
				ID:               m.Id,
				TombstonedAt:     tombstoned,
				Released:         released,
				Title:            m.Title,
				Popularity:       m.Popularity,
				Adult:            m.Adult,
				BackdropPath:     m.BackdropPath,
				PosterPath:       m.PosterPath,
				OriginalLanguage: m.OriginalLanguage,
				OriginalTitle:    m.OriginalTitle,
				Overview:         m.Overview,
				Mimetype:         m.Mimetype,
				Collation:        m.Collation,
				Subtitle:         m.Subtitle,
				ParentUID:        m.ParentUid,
			}); err != nil {
				return errorsx.Wrap(err, "unable to encode media")
			}
		}

		cursor = slicesx.LastOrZero(resp.Items...).SyncId
		log.Println("syncing media metadata", cursor)
	}
}
