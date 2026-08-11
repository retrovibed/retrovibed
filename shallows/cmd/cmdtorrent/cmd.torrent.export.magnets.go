package cmdtorrent

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// torrentRecord is the JSONL transfer encoding for export/import peer.
type torrentRecord struct {
	Magnet         string `json:"magnet"`
	EncryptionSeed string `json:"encryption_seed,omitempty"`
}

type exportMagnets struct {
	Path         string    `arg:"" name:"path" help:"file to write magnet urls out to, defaults to stdout" default:"" required:"false"`
	Completed    bool      `flag:"" name:"completed" help:"only export completed torrents" default:"false"`
	Hidden       bool      `flag:"" name:"hidden" help:"include hidden torrents" default:"false"`
	CreatedAfter time.Time `flag:"" name:"created-after" help:"only export torrents created after this timestamp" required:"false"`
}

func (t exportMagnets) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	dst := os.Stdout
	if stringsx.Present(t.Path) {
		if dst, err = os.OpenFile(t.Path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600); err != nil {
			return err
		}
	}

	db, err := cmdopts.DatabaseMeta(gctx.Context)
	if err != nil {
		return err
	}
	defer db.Close()

	filters := squirrel.And{tracking.MetadataQueryNotTombstoned()}
	if !t.Hidden {
		filters = append(filters, tracking.MetadataQueryNotHidden())
	}
	if t.Completed {
		filters = append(filters, tracking.MetadataQueryCompleted(true))
	}
	if !t.CreatedAfter.IsZero() {
		filters = append(filters, tracking.MetadataQueryCreatedAfter(t.CreatedAfter))
	}

	q := tracking.MetadataSearchBuilder().Where(filters)
	enc := json.NewEncoder(dst)

	scanner := sqlx.Scan(tracking.MetadataSearch(gctx.Context, db, q))
	for lmd := range scanner.Iter() {
		mg := metainfo.NewMagnetFromInfohash(
			lmd.Infohash,
			metainfo.MagnetOptionTrackers(lmd.Tracker),
			metainfo.MagnetOptionDisplayName(lmd.Description),
		)
		rec := torrentRecord{Magnet: mg.String(), EncryptionSeed: lmd.EncryptionSeed}
		if err = enc.Encode(rec); err != nil {
			return err
		}
	}

	return scanner.Err()
}
