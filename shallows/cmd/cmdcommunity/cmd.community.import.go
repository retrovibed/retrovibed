package cmdcommunity

import (
	"database/sql"
	"io"
	"os"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/community"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/jsonl"
	"github.com/retrovibed/retrovibed/internal/stringsx"
	"github.com/retrovibed/retrovibed/meta"
)

type cmdCommunityImport struct {
	DryRun bool   `flag:"" name:"dry-run" help:"validate and show what would be imported without making changes"`
	Input  string `flag:"" name:"input" help:"input file path" required:"true"`
	Output string `flag:"" name:"output" help:"output file path" required:"true"`
}

func (t cmdCommunityImport) Run(gctx *cmdopts.Global) (err error) {
	var (
		db            *sql.DB
		comm          meta.Community
		pc            meta.PublishedContent
		input, output *os.File
	)

	if input, err = os.Open(t.Input); err != nil {
		return errorsx.Wrap(err, "failed to open input file")
	}
	defer input.Close()

	if output, err = os.Create(t.Output); err != nil {
		return errorsx.Wrap(err, "failed to create output file")
	}
	defer output.Close()

	dec := jsonl.NewDecoder(input)
	enc := jsonl.NewEncoder(output)

	if err = dec.Decode(&comm); err != nil {
		return errorsx.Wrap(err, "failed to decode community")
	}

	if db, err = cmdmeta.Database(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	for err = dec.Decode(&pc); err == nil; err = dec.Decode(&pc) {
		pc.CommunityId = comm.Id

		if t.DryRun {
			err = enc.Encode(&pc)
		} else {
			dbpc := community.PublishedContent{
				CommunityID:   pc.CommunityId,
				KnownMediaID:  pc.KnownMediaId,
				MagnetURI:     pc.MagnetUri,
				LibraryID:     stringsx.FirstNonBlank(pc.LibraryId, uuid.Nil.String()),
				OAuthGoogleID: uuid.Nil.String(),
			}
			err = community.PublishedContentInsertWithDefaults(gctx.Context, db, dbpc).Scan(&dbpc)
		}

		if err != nil {
			return err
		}
		pc = meta.PublishedContent{}
	}

	if err != io.EOF {
		return err
	}

	return nil
}
