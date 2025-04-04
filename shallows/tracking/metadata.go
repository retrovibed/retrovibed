package tracking

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/internal/duckdbx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/internal/slicesx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/squirrelx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/library"
)

func MetadataOptionNoop(*Metadata) {}

func MetadataOptionInitiate(md *Metadata) {
	md.InitiatedAt = time.Now()
}

func MetadataOptionFromInfo(i *metainfo.Info) func(*Metadata) {
	return func(m *Metadata) {
		m.Description = strings.ToValidUTF8(i.Name, "\uFFFD")
		m.Bytes = uint64(i.TotalLength())
		m.Private = langx.Autoderef(i.Private)
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

func MetadataQuerySeeding() squirrel.Sqlizer {
	return squirrel.Expr("torrents_metadata.seeding")
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

func Verify(ctx context.Context, t torrent.Torrent) error {
	log.Println("verify initiated", t.Metadata().DisplayName)
	defer log.Println("verify completed", t.Metadata().DisplayName, spew.Sdump(t.Stats()))
	return torrent.Verify(ctx, t)
}

func Download(ctx context.Context, q sqlx.Queryer, vfs fsx.Virtual, md *Metadata, t torrent.Torrent) (err error) {
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

	mediavfs := fsx.DirVirtual(vfs.Path("media"))

	log.Println("content transfer to library initiated")
	defer log.Println("content transfer to library completed")
	// need to get the path to the torrent media.
	for tx, cause := range library.ImportFilesystem(ctx, library.ImportSymlinkFile(mediavfs), vfs.Path("torrent", t.Metadata().ID.HexString())) {
		if cause != nil {
			log.Println(cause)
			err = errorsx.Compact(err, cause)
			continue
		}

		lmd := library.NewMetadata(
			md5x.FormatUUID(tx.MD5),
			library.MetadataOptionDescription(filepath.Base(tx.Path)),
			library.MetadataOptionBytes(tx.Bytes),
			library.MetadataOptionTorrentID(md.ID),
			library.MetadataOptionMimetype(tx.Mimetype.String()),
		)

		if err := library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd); err != nil {
			return errorsx.Wrap(err, "unable to record library metadata")
		}

		log.Println("new library content", spew.Sdump(lmd))
	}

	if err != nil {
		return errorsx.Wrap(err, "failed to transfer files into library")
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
	// Revisit once resume is working.
	// sub := dl.SubscribePieceStateChanges()
	// defer sub.Close()

	statst := time.NewTimer(statsfreq)
	// l := rate.NewLimiter(rate.Every(time.Second), 1)
	for {
		select {
		case <-statst.C:
			stats := dl.Stats()
			log.Printf("%s - %s: seeding(%t), peers(%d:%d:%d) pieces(%d:%d:%d:%d)\n", md.ID, hex.EncodeToString(md.Infohash), stats.Seeding, stats.ActivePeers, stats.PendingPeers, stats.TotalPeers, stats.Missing, stats.Outstanding, stats.Unverified, stats.Completed)
		// case <-sub.Values:
		// 	if !l.Allow() {
		// 		continue
		// 	}

		// 	statst.Reset(statsfreq)

		// 	current := uint64(dl.BytesCompleted())
		// 	if md.Downloaded == current {
		// 		continue
		// 	}

		// 	stats := dl.Stats()
		// 	log.Printf("%s: peers(%d:%d:%d) pieces(%d:%d:%d:%d)\n", dl.Metainfo().HashInfoBytes().HexString(), stats.ActivePeers, stats.PendingPeers, stats.TotalPeers, stats.Missing, stats.Outstanding, stats.Unverified, stats.Completed)

		// 	if err := MetadataProgressByID(ctx, q, md.ID, uint16(stats.ActivePeers), current).Scan(md); err != nil {
		// 		log.Println("failed to update progress", err)
		// 	}
		case <-ctx.Done():
			return
		}
	}
}
