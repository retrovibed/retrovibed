package tracking

import (
	"context"
	"io"
	"log"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/deeppool/internal/x/duckdbx"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/slicesx"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/internal/x/squirrelx"
	"github.com/james-lawrence/deeppool/internal/x/timex"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"golang.org/x/time/rate"
)

func MetadataOptionNoop(*Metadata) {}

func MetadataOptionInitiate(md *Metadata) {
	md.InitiatedAt = time.Now()
}

func MetadataOptionFromInfo(i *metainfo.Info) func(*Metadata) {
	return func(m *Metadata) {
		m.Description = strings.ToValidUTF8(i.Name, "\uFFFD")
		m.Bytes = uint64(i.TotalLength())
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
	)

	pctx, done := context.WithCancel(ctx)
	defer done()

	// update the progress.
	go DownloadProgress(pctx, q, md, t)

	// just copying as we receive data to block until done.
	if downloaded, err = torrent.DownloadInto(ctx, io.Discard, t); err != nil {
		return errorsx.Wrap(err, "download failed")
	}

	log.Println("download completed", md.ID, md.Description, downloaded)
	if err := MetadataProgressByID(ctx, q, md.ID, 0, uint64(downloaded)).Scan(md); err != nil {
		return errorsx.Wrap(err, "progress update failed")
	}

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
			log.Printf("%s: peers(%d:%d:%d) pieces(%d:%d:%d:%d)\n", dl.Metainfo().HashInfoBytes().HexString(), stats.ActivePeers, stats.PendingPeers, stats.TotalPeers, stats.Missing, stats.Outstanding, stats.Unverified, stats.Completed)
		case <-sub.Values:
			statst.Reset(statsfreq)
			if !l.Allow() {
				continue
			}

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
