package tracking

import (
	"context"
	"crypto/md5"
	"log"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/deeppool/internal/x/duckdbx"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/md5x"
	"github.com/james-lawrence/deeppool/internal/x/slicesx"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/internal/x/squirrelx"
	"github.com/james-lawrence/deeppool/internal/x/timex"
	"github.com/james-lawrence/deeppool/library"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"golang.org/x/time/rate"

	"github.com/gabriel-vasile/mimetype"
)

func MetadataOptionNoop(*Metadata) {}

func MetadataOptionInitiate(md *Metadata) {
	md.InitiatedAt = time.Now()
}

func MetadataOptionFromInfo(i *metainfo.Info) func(*Metadata) {
	return func(m *Metadata) {
		m.Description = strings.ToValidUTF8(i.Name, "\uFFFD")
		m.Bytes = uint64(i.TotalLength())
		m.Private = i.Private
	}
}

func MetadataOptionDescription(d string) func(*Metadata) {
	return func(m *Metadata) {
		m.Description = d
	}
}

// Currently will select just the first tracker due to poor list support in duckdb.
func MetadataOptionTrackers(d ...string) func(*Metadata) {
	return func(m *Metadata) {
		m.Tracker = slicesx.FirstOrZero(d...)
	}
}

func MetadataOptionJSONSafeEncode(p *Metadata) {
	p.CreatedAt = timex.RFC3339NanoEncode(p.CreatedAt)
	p.UpdatedAt = timex.RFC3339NanoEncode(p.UpdatedAt)
}

func NewMetadata(md *metainfo.Hash, options ...func(*Metadata)) (m Metadata) {
	r := langx.Clone(Metadata{
		ID:       HashUID(md),
		Infohash: md.Bytes(),
	}, options...)
	return r
}

func MetadataQueryNotInitiated() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.initiated_at = 'infinity'")
}

func MetadataQueryInitiated() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.initiated_at < NOW()")
}

func MetadataQueryIncomplete() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.downloaded < torrents_metadata.bytes")
}

func MetadataQueryNotPaused() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.paused_at = 'infinity'")
}

func MetadataSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) MetadataScanner {
	return NewMetadataScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func MetadataQuerySearch(q string, columns ...string) squirrel.Sqlizer {
	return duckdbx.FTSSearch("fts_main_torrents_metadata", q, columns...)
}

func MetadataSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(MetadataScannerStaticColumns)...).From("torrents_metadata")
}

func Download(ctx context.Context, q sqlx.Queryer, md *Metadata, t torrent.Torrent) (err error) {
	var (
		downloaded int64
		mhash      = md5.New()
	)

	pctx, done := context.WithCancel(ctx)
	defer done()

	// update the progress.
	go DownloadProgress(pctx, q, md, t)

	// just copying as we receive data to block until done.
	if downloaded, err = torrent.DownloadInto(ctx, mhash, t); err != nil {
		return errorsx.Wrap(err, "download failed")
	}

	log.Println("download completed", md.ID, md.Description, downloaded)
	if err := MetadataProgressByID(ctx, q, md.ID, 0, uint64(downloaded)).Scan(md); err != nil {
		return errorsx.Wrap(err, "progress update failed")
	}

	content := t.NewReader()
	defer content.Close()
	cmimetype, err := mimetype.DetectReader(content)
	if err != nil {
		return errorsx.Wrap(err, "unable to determine mimetype")
	}

	lmd := library.NewMetadata(
		md5x.FormatString(mhash),
		library.MetadataOptionDescription(md.Description),
		library.MetadataOptionBytes(md.Bytes),
		library.MetadataOptionTorrentID(md.ID),
		library.MetadataOptionMimetype(cmimetype.String()),
	)

	if err := library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd); err != nil {
		return errorsx.Wrap(err, "unable to record library metadata")
	}

	log.Println("new library content", spew.Sdump(lmd))
	return nil
}

func DownloadProgress(ctx context.Context, q sqlx.Queryer, md *Metadata, dl torrent.Torrent) {
	const (
		statsfreq = 10 * time.Second
	)
	log.Println("monitoring download progress initiated", md.ID, md.Description, md.Tracker)
	defer log.Println("monitoring download progress completed", md.ID, md.Description, md.Tracker)
	sub := dl.SubscribePieceStateChanges()
	defer sub.Close()

	statst := time.NewTimer(statsfreq)
	l := rate.NewLimiter(rate.Every(time.Second), 1)
	for {
		select {
		case <-statst.C:
			stats := dl.Stats()
			log.Printf("%s: seeding(%t), peers(%d:%d:%d) pieces(%d:%d:%d:%d)\n", dl.Metainfo().HashInfoBytes().HexString(), stats.Seeding, stats.ActivePeers, stats.PendingPeers, stats.TotalPeers, stats.Missing, stats.Outstanding, stats.Unverified, stats.Completed)
		case <-sub.Values:
			if !l.Allow() {
				continue
			}

			statst.Reset(statsfreq)

			current := uint64(dl.BytesCompleted())
			if md.Downloaded == current {
				continue
			}

			stats := dl.Stats()
			log.Printf("%s: peers(%d:%d:%d) pieces(%d:%d:%d:%d)\n", dl.Metainfo().HashInfoBytes().HexString(), stats.ActivePeers, stats.PendingPeers, stats.TotalPeers, stats.Missing, stats.Outstanding, stats.Unverified, stats.Completed)

			if err := MetadataProgressByID(ctx, q, md.ID, uint16(stats.ActivePeers), current).Scan(md); err != nil {
				log.Println("failed to update progress", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
